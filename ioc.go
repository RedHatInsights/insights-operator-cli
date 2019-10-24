/*
Copyright © 2019 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/c-bata/go-prompt"
	. "github.com/logrusorgru/aurora"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const API_PREFIX = "/api/v1/"

var BuildVersion string = "*not set*"
var BuildTime string = "*not set*"

var controllerUrl string
var username string
var password string
var files []prompt.Suggest

func tryToLogin(username string, password string) {
	fmt.Println(Blue("\nDone"))
}

type Cluster struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ConfigurationProfile struct {
	Id            int    `json:"id"`
	Configuration string `json:"configuration"`
	ChangedAt     string `json:"changed_at"`
	ChangedBy     string `json:"changed_by"`
	Description   string `json:"description"`
}

type ClusterConfiguration struct {
	Id            int    `json:"id"`
	Cluster       string `json:"cluster"`
	Configuration string `json:"configuration"`
	ChangedAt     string `json:"changed_at"`
	ChangedBy     string `json:"changed_by"`
	Active        string `json:"active"`
	Reason        string `json:"reason"`
}

func performReadRequest(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Communication error with the server %v", err)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Expected HTTP status 200 OK, got %d", response.StatusCode)
	}
	body, readErr := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if readErr != nil {
		return nil, fmt.Errorf("Unable to read response body")
	}

	return body, nil
}

func performWriteRequest(url string, method string, payload io.Reader) error {
	var client http.Client

	request, err := http.NewRequest(method, url, payload)
	if err != nil {
		return fmt.Errorf("Error creating request %v", err)
	}

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("Communication error with the server %v", err)
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected HTTP status 200 OK, got %d", response.StatusCode)
	}
	return nil
}

func readListOfClusters(controllerUrl string, apiPrefix string) ([]Cluster, error) {
	clusters := []Cluster{}

	url := controllerUrl + apiPrefix + "client/cluster"
	body, err := performReadRequest(url)

	err = json.Unmarshal(body, &clusters)
	if err != nil {
		return nil, err
	}
	return clusters, nil
}

func readListOfConfigurationProfiles(controllerUrl string, apiPrefix string) ([]ConfigurationProfile, error) {
	profiles := []ConfigurationProfile{}

	url := controllerUrl + apiPrefix + "client/profile"
	body, err := performReadRequest(url)

	err = json.Unmarshal(body, &profiles)
	if err != nil {
		return nil, err
	}
	return profiles, nil
}

func readListOfConfigurations(controllerUrl string, apiPrefix string) ([]ClusterConfiguration, error) {
	configurations := []ClusterConfiguration{}

	url := controllerUrl + apiPrefix + "client/configuration"
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &configurations)
	if err != nil {
		return nil, err
	}
	return configurations, nil
}

func readConfigurationProfile(controllerUrl string, apiPrefix string, profileId string) (*ConfigurationProfile, error) {
	var profile ConfigurationProfile
	url := controllerUrl + apiPrefix + "client/profile/" + profileId
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func readClusterConfigurationById(controllerUrl string, apiPrefix string, configurationId string) (*string, error) {
	url := controllerUrl + apiPrefix + "client/configuration/" + configurationId
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	str := string(body)
	return &str, nil
}

func listOfClusters() {
	clusters, err := readListOfClusters(controllerUrl, API_PREFIX)
	if err != nil {
		fmt.Println(Red("Error reading list of clusters"))
		fmt.Println(err)
		return
	}

	fmt.Println(Magenta("List of clusters"))
	fmt.Printf("%4s %4s %-s\n", "#", "ID", "Name")
	for i, cluster := range clusters {
		fmt.Printf("%4d %4d %-s\n", i, cluster.Id, cluster.Name)
	}
}

func listOfProfiles() {
	profiles, err := readListOfConfigurationProfiles(controllerUrl, API_PREFIX)
	if err != nil {
		fmt.Println(Red("Error reading list of configuration profiles"))
		fmt.Println(err)
		return
	}

	fmt.Println(Magenta("List of configuration profiles"))
	fmt.Printf("%4s %4s %-20s %-20s %s\n", "#", "ID", "Changed at", "Changed by", "Description")
	for i, profile := range profiles {
		changedAt := profile.ChangedAt[0:19]
		fmt.Printf("%4d %4d %-20s %-20s %-s\n", i, profile.Id, changedAt, profile.ChangedBy, profile.Description)
	}
}

func listOfConfigurations(filter string) {
	// TODO: filter in query?
	configurations, err := readListOfConfigurations(controllerUrl, API_PREFIX)
	if err != nil {
		fmt.Println(Red("Error reading list of configurations"))
		fmt.Println(err)
		return
	}

	fmt.Println(Magenta("List of configuration for all clusters"))
	fmt.Printf("%4s %4s %-20s %-20s %-10s %-12s %s\n", "#", "ID", "Cluster", "Changed at", "Changed by", "Active", "Reason")
	for i, configuration := range configurations {
		// poor man's filtering
		if strings.Contains(configuration.Cluster, filter) {
			var active Value
			if configuration.Active == "1" {
				active = Green("yes")
			} else {
				active = Red("no")
			}
			changedAt := configuration.ChangedAt[0:19]
			fmt.Printf("%4d %4d %-20s %-20s %-10s %-12s %s\n", i, configuration.Id, configuration.Cluster, changedAt, configuration.ChangedBy, active, configuration.Reason)
		}
	}
}

func describeProfile(profileId string) {
	profile, err := readConfigurationProfile(controllerUrl, API_PREFIX, profileId)
	if err != nil {
		fmt.Println(Red("Error reading configuration profile"))
		fmt.Println(err)
		return
	}

	fmt.Println(Magenta("Configuration profile"))
	fmt.Println(profile.Configuration)
}

func describeConfiguration(clusterId string) {
	configuration, err := readClusterConfigurationById(controllerUrl, API_PREFIX, clusterId)
	if err != nil {
		fmt.Println(Red("Error reading cluster configuration"))
		fmt.Println(err)
		return
	}

	fmt.Println(Magenta("Configuration for cluster " + clusterId))
	fmt.Println(*configuration)
}

func enableClusterConfiguration(configurationId string) {
	url := controllerUrl + API_PREFIX + "client/configuration/" + configurationId + "/enable"
	err := performWriteRequest(url, "PUT", nil)
	if err != nil {
		fmt.Println(Red("Error communicating with the service"))
		fmt.Println(err)
		return
	} else {
		fmt.Print(Blue("Configuration " + configurationId + " has been "))
		fmt.Println(Green("enabled"))
	}
}

func disableClusterConfiguration(configurationId string) {
	// TODO: refactoring needed - almost the same code as in previous function
	url := controllerUrl + API_PREFIX + "client/configuration/" + configurationId + "/disable"
	err := performWriteRequest(url, "PUT", nil)
	if err != nil {
		fmt.Println(Red("Error communicating with the service"))
		fmt.Println(err)
		return
	} else {
		fmt.Print(Blue("Configuration " + configurationId + " has been "))
		fmt.Println(Red("disabled"))
	}
}

func addClusterConfiguration() {
	if username == "" {
		fmt.Println(Red("Not logged in"))
		return
	}

	cluster := prompt.Input("cluster: ", loginCompleter)
	if cluster == "" {
		fmt.Println(Red("Cancelled"))
		return
	}

	reason := prompt.Input("reason: ", loginCompleter)
	if reason == "" {
		fmt.Println(Red("Cancelled"))
		return
	}

	description := prompt.Input("description: ", loginCompleter)
	if description == "" {
		fmt.Println(Red("Cancelled"))
		return
	}

	// TODO: make the directory fully configurable
	err := fillInConfigurationList("configurations")
	if err != nil {
		fmt.Println(Red("Cannot read any configuration file"))
		fmt.Println(err)
	}

	configurationFileName := prompt.Input("configuration file (TAB to complete): ", configFileCompleter)
	if configurationFileName == "" {
		fmt.Println(Red("Cancelled"))
		return
	}

	// TODO: make the directory fully configurable
	configuration, err := ioutil.ReadFile("configurations/" + configurationFileName)
	if err != nil {
		fmt.Println(Red("Cannot read configuration file"))
		fmt.Println(err)
	}

	query := "username=" + url.QueryEscape(username) + "&reason=" + url.QueryEscape(reason) + "&description=" + url.QueryEscape(description)
	url := controllerUrl + API_PREFIX + "client/cluster/" + url.PathEscape(cluster) + "/configuration?" + query

	err = performWriteRequest(url, "POST", bytes.NewReader(configuration))
	if err != nil {
		fmt.Println(Red("Error communicating with the service"))
		fmt.Println(err)
		return
	} else {
		fmt.Println(Blue("Configuration has been created"))
	}
}

func fillInConfigurationList(directory string) error {
	files = []prompt.Suggest{}

	root := directory
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			suggest := prompt.Suggest{
				Text: info.Name()}
			files = append(files, suggest)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func printHelp() {
	fmt.Println(Magenta("HELP:"))
	fmt.Println()
	fmt.Println(Blue("Cluster operations:        "))
	fmt.Println(Yellow("list clusters            "), "list all clusters known to the service")
	fmt.Println()
	fmt.Println(Blue("Configuration profiles:    "))
	fmt.Println(Yellow("list profiles            "), "list all profiles known to the service")
	fmt.Println(Yellow("describe profile ##      "), "describe profile selected by its ID")
	fmt.Println()
	fmt.Println(Blue("Cluster configurations:    "))
	fmt.Println(Yellow("list configurations      "), "list all configurations known to the service")
	fmt.Println(Yellow("describe configuration ##"), "describe cluster configuration selected by its ID")
	fmt.Println(Yellow("add configuration        "), "add new configuration")
	fmt.Println(Yellow("new configuration        "), "alias for previous command")
	fmt.Println(Yellow("enable ##                "), "enable cluster configuration selected by its ID")
	fmt.Println(Yellow("disable ##               "), "disable cluster configuration selected by its ID")
	fmt.Println()
	fmt.Println(Blue("Other commands:"))
	fmt.Println(Yellow("quit                     "), "quit the application")
	fmt.Println(Yellow("exit                     "), "dtto")
	fmt.Println(Yellow("bye                      "), "dtto")
	fmt.Println(Yellow("help                     "), "this help")
	fmt.Println()
}

func loginCompleter(in prompt.Document) []prompt.Suggest {
	return nil
}

func configFileCompleter(in prompt.Document) []prompt.Suggest {
	return prompt.FilterHasPrefix(files, in.Text, true)
}

func printVersion() {
	fmt.Println(Blue("Insights operator CLI client "), "version", Yellow(BuildVersion), "compiled", Yellow(BuildTime))
}

func executor(t string) {
	blocks := strings.Split(t, " ")
	// commands with variable parts
	switch {
	case strings.HasPrefix(t, "describe profile "):
		describeProfile(blocks[2])
		return
	case strings.HasPrefix(t, "describe configuration "):
		describeConfiguration(blocks[2])
		return
	case strings.HasPrefix(t, "enable "):
		enableClusterConfiguration(blocks[1])
		return
	case strings.HasPrefix(t, "disable "):
		disableClusterConfiguration(blocks[1])
		return
	case strings.HasPrefix(t, "list configurations "):
		listOfConfigurations(blocks[2])
		return
	}

	// fixed commands
	switch t {
	case "login":
		username = prompt.Input("login: ", loginCompleter)
		fmt.Print("password: ")
		p, err := terminal.ReadPassword(0)
		if err != nil {
			fmt.Println(Red("Password is not set"))
		} else {
			password = string(p)
			tryToLogin(username, password)
		}
	case "list clusters":
		listOfClusters()
	case "list profiles":
		listOfProfiles()
	case "list configurations":
		listOfConfigurations("")
	case "add configuration":
		fallthrough
	case "new configuration":
		addClusterConfiguration()
	case "describe profile":
		profile := prompt.Input("profile: ", loginCompleter)
		describeProfile(profile)
	case "describe configuration":
		configuration := prompt.Input("configuration: ", loginCompleter)
		describeConfiguration(configuration)
	case "enable":
		configuration := prompt.Input("configuration: ", loginCompleter)
		enableClusterConfiguration(configuration)
	case "disable":
		configuration := prompt.Input("configuration: ", loginCompleter)
		disableClusterConfiguration(configuration)
	case "bye":
		fallthrough
	case "exit":
		fallthrough
	case "quit":
		fmt.Println(Magenta("Quitting"))
		os.Exit(0)
	case "?":
		fallthrough
	case "help":
		printHelp()
	case "version":
		printVersion()
	default:
		fmt.Println("Command not found")
	}
}

func completer(in prompt.Document) []prompt.Suggest {
	firstWord := []prompt.Suggest{
		{Text: "login", Description: "provide login info"},
		{Text: "help", Description: "show help with all commands"},
		{Text: "exit", Description: "quit the application"},
		{Text: "quit", Description: "quit the application"},
		{Text: "bye", Description: "quit the application"},
		{Text: "list", Description: "list resources (clusters, profiles, configurations)"},
		{Text: "describe", Description: "describe the selected resource"},
		{Text: "add", Description: "add resource (cluster, profile, configuration)"},
		{Text: "new", Description: "alias for add"},
		{Text: "enable", Description: "enable selected cluster profile"},
		{Text: "disable", Description: "disable selected cluster profile"},
		{Text: "version", Description: "prints the build information for CLI executable"},
	}

	secondWord := make(map[string][]prompt.Suggest)

	// list operations
	secondWord["list"] = []prompt.Suggest{
		{Text: "clusters", Description: "show list of all clusters available"},
		{Text: "profiles", Description: "show list of all configuration profiles"},
		{Text: "configurations", Description: "show list all cluster configurations"},
	}
	// add operations
	secondWord["add"] = []prompt.Suggest{
		{Text: "cluster", Description: "add/register new cluster"},
		{Text: "profile", Description: "add new configuration profile"},
		{Text: "configuration", Description: "add new cluster configuration"},
	}
	secondWord["new"] = []prompt.Suggest{
		{Text: "cluster", Description: "add/register new cluster"},
		{Text: "profile", Description: "add new configuration profile"},
		{Text: "configuration", Description: "add new cluster configuration"},
	}
	// descripbe operations
	secondWord["describe"] = []prompt.Suggest{
		{Text: "profile", Description: "describe selected configuration profile"},
		{Text: "configuration", Description: "describe configuration for selected cluster"},
	}

	empty_s := []prompt.Suggest{}

	blocks := strings.Split(in.TextBeforeCursor(), " ")

	if len(blocks) == 2 {
		sec, ok := secondWord[blocks[0]]
		if ok {
			return prompt.FilterHasPrefix(sec, blocks[1], true)
		} else {
			return empty_s
		}
	}
	if in.GetWordBeforeCursor() == "" {
		return nil
	} else {
		return prompt.FilterHasPrefix(firstWord, blocks[0], true)
	}
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	controllerUrl = viper.GetString("CONTROLLER_URL")
	p := prompt.New(executor, completer)
	p.Run()
}
