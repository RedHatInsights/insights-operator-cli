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

// Implementation of command-line client for the insights operator instrumentation service.
package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	. "github.com/logrusorgru/aurora"
	"github.com/redhatinsighs/insights-operator-cli/commands"
	"github.com/redhatinsighs/insights-operator-cli/restapi"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var BuildVersion string = "*not set*"
var BuildTime string = "*not set*"

var controllerUrl string
var username string
var password string
var files []prompt.Suggest
var api restapi.Api

func tryToLogin(username string, password string) {
	fmt.Println(Blue("\nDone"))
}

func listOfConfigurations(filter string) {
	// TODO: filter in query?
	configurations, err := api.ReadListOfConfigurations()
	if err != nil {
		fmt.Println(Red("Error reading list of configurations"))
		fmt.Println(err)
		return
	}

	fmt.Println(Magenta("List of configuration for all clusters"))
	fmt.Printf("%4s %4s %4s    %-20s %-20s %-10s %-12s %s\n", "#", "ID", "Profile", "Cluster", "Changed at", "Changed by", "Active", "Reason")
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
			fmt.Printf("%4d %4d %4s       %-20s %-20s %-10s %-12s %s\n", i, configuration.Id, configuration.Configuration, configuration.Cluster, changedAt, configuration.ChangedBy, active, configuration.Reason)
		}
	}
}

func describeProfile(profileId string) {
	profile, err := api.ReadConfigurationProfile(profileId)
	if err != nil {
		fmt.Println(Red("Error reading configuration profile"))
		fmt.Println(err)
		return
	}

	fmt.Println(Magenta("Configuration profile"))
	fmt.Println(profile.Configuration)
}

func describeConfiguration(clusterId string) {
	configuration, err := api.ReadClusterConfigurationById(clusterId)
	if err != nil {
		fmt.Println(Red("Error reading cluster configuration"))
		fmt.Println(err)
		return
	}

	fmt.Println(Magenta("Configuration for cluster " + clusterId))
	fmt.Println(*configuration)
}

func proceedQuestion(question string) bool {
	fmt.Println(Red(question))
	proceed := prompt.Input("proceed? [y/n] ", loginCompleter)
	if proceed != "y" {
		fmt.Println(Blue("cancelled"))
		return false
	}
	return true
}

func deleteCluster(clusterId string) {
	if !proceedQuestion("All cluster configurations will be deleted") {
		return
	}

	err := api.DeleteCluster(clusterId)
	if err != nil {
		fmt.Println(Red("Error communicating with the service"))
		fmt.Println(err)
		return
	} else {
		fmt.Println(Blue("Cluster "+clusterId+" has been"), Red("deleted"))
	}
}

func deleteClusterConfiguration(configurationId string) {
	err := api.DeleteClusterConfiguration(configurationId)
	if err != nil {
		fmt.Println(Red("Error communicating with the service"))
		fmt.Println(err)
		return
	} else {
		fmt.Println(Blue("Configuration "+configurationId+" has been "), Red("deleted"))
	}
}

func deleteConfigurationProfile(profileId string) {
	if !proceedQuestion("All configurations based on this profile will be deleted") {
		return
	}

	err := api.DeleteConfigurationProfile(profileId)
	if err != nil {
		fmt.Println(Red("Error communicating with the service"))
		fmt.Println(err)
		return
	} else {
		fmt.Println(Blue("Configuration profile "+profileId+" has been "), Red("deleted"))
	}
}

func addCluster() {
	id := prompt.Input("ID: ", loginCompleter)
	if id == "" {
		fmt.Println(Red("Cancelled"))
		return
	}

	name := prompt.Input("name: ", loginCompleter)
	if name == "" {
		fmt.Println(Red("Cancelled"))
		return
	}

	err := api.AddCluster(id, name)
	if err != nil {
		fmt.Println(Red("Error communicating with the service"))
		fmt.Println(err)
		return
	} else {
		fmt.Println(Blue("Configuration profile has been created"))
	}
}

func addProfile() {
	if username == "" {
		fmt.Println(Red("Not logged in"))
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

	err = api.AddConfigurationProfile(username, description, configuration)
	if err != nil {
		fmt.Println(Red("Error communicating with the service"))
		fmt.Println(err)
		return
	} else {
		fmt.Println(Blue("Configuration profile has been created"))
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

	err = api.AddClusterConfiguration(username, cluster, reason, description, configuration)
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

func addTrigger() {
	if username == "" {
		fmt.Println(Red("Not logged in"))
		return
	}

	clusterName := prompt.Input("cluster name: ", loginCompleter)
	reason := prompt.Input("reason: ", loginCompleter)
	link := prompt.Input("link: ", loginCompleter)

	err := api.AddTrigger(username, clusterName, reason, link)
	if err != nil {
		fmt.Println("Error communicating with the service")
		fmt.Println(err)
		return
	} else {
		fmt.Println(Blue("Trigger has been created"))
	}
}

func describeTrigger(triggerId string) {
	trigger, err := api.ReadTriggerById(triggerId)
	if err != nil {
		fmt.Println(Red("Error reading selected trigger"))
		fmt.Println(err)
		return
	}

	var active Value
	if trigger.Active == 1 {
		active = Green("yes")
	} else {
		active = Red("no")
	}

	triggeredAt := trigger.TriggeredAt[0:19]
	ackedAt := trigger.AckedAt[0:19]

	var ttype Value
	if trigger.Type == "must-gather" {
		ttype = Blue(trigger.Type)
	} else {
		ttype = Magenta(trigger.Type)
	}

	fmt.Println(Magenta("Trigger info"))
	fmt.Printf("ID:            %d\n", trigger.Id)
	fmt.Printf("Type:          %s\n", ttype)
	fmt.Printf("Cluster:       %s\n", trigger.Cluster)
	fmt.Printf("Triggered at:  %s\n", triggeredAt)
	fmt.Printf("Triggered by:  %s\n", trigger.TriggeredBy)
	fmt.Printf("Active:        %s\n", active)
	fmt.Printf("Acked at:      %s\n", ackedAt)
}

func deleteTrigger(triggerId string) {
	err := api.DeleteTrigger(triggerId)
	if err != nil {
		fmt.Println(Red("Error communicating with the service"))
		fmt.Println(err)
		return
	} else {
		fmt.Println(Blue("Trigger "+triggerId+" has been"), Red("deleted"))
	}
}

func activateTrigger(triggerId string) {
	err := api.ActivateTrigger(triggerId)
	if err != nil {
		fmt.Println(Red("Error communicating with the service"))
		fmt.Println(err)
		return
	} else {
		fmt.Println(Blue("Trigger "+triggerId+" has been"), Green("activated"))
	}
}

func deactivateTrigger(triggerId string) {
	err := api.DeactivateTrigger(triggerId)
	if err != nil {
		fmt.Println(Red("Error communicating with the service"))
		fmt.Println(err)
		return
	} else {
		fmt.Println(Blue("Trigger "+triggerId+" has been"), Green("deactivated"))
	}
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
	case strings.HasPrefix(t, "describe trigger "):
		describeTrigger(blocks[2])
		return
	case strings.HasPrefix(t, "enable "):
		commands.EnableClusterConfiguration(api, blocks[1])
		return
	case strings.HasPrefix(t, "disable "):
		commands.DisableClusterConfiguration(api, blocks[1])
		return
	case strings.HasPrefix(t, "list configurations "):
		listOfConfigurations(blocks[2])
		return
	case strings.HasPrefix(t, "delete cluster "):
		deleteCluster(blocks[2])
		return
	case strings.HasPrefix(t, "delete configuration "):
		deleteClusterConfiguration(blocks[2])
		return
	case strings.HasPrefix(t, "delete profile "):
		deleteConfigurationProfile(blocks[2])
		return
	case strings.HasPrefix(t, "delete trigger "):
		deleteTrigger(blocks[2])
		return
	case strings.HasPrefix(t, "activate must gather "):
		fallthrough
	case strings.HasPrefix(t, "activate trigger "):
		activateTrigger(blocks[2])
		return
	case strings.HasPrefix(t, "deactivate must gather "):
		fallthrough
	case strings.HasPrefix(t, "deactivate trigger "):
		deactivateTrigger(blocks[2])
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
	case "list must-gather":
		fallthrough
	case "list triggers":
		commands.ListOfTriggers(api)
	case "list clusters":
		commands.ListOfClusters(api)
	case "list profiles":
		commands.ListOfProfiles(api)
	case "list configurations":
		listOfConfigurations("")
	case "add cluster":
		fallthrough
	case "new cluster":
		addCluster()
	case "add profile":
		fallthrough
	case "new profile":
		addProfile()
	case "add configuration":
		fallthrough
	case "new configuration":
		addClusterConfiguration()
	case "request must-gather":
		fallthrough
	case "add trigger":
		fallthrough
	case "new trigger":
		addTrigger()
	case "describe profile":
		profile := prompt.Input("profile: ", loginCompleter)
		describeProfile(profile)
	case "describe configuration":
		configuration := prompt.Input("configuration: ", loginCompleter)
		describeConfiguration(configuration)
	case "describe must-gather":
		fallthrough
	case "describe trigger":
		trigger := prompt.Input("trigger: ", loginCompleter)
		describeTrigger(trigger)
	case "enable":
		configuration := prompt.Input("configuration: ", loginCompleter)
		commands.EnableClusterConfiguration(api, configuration)
	case "disable":
		configuration := prompt.Input("configuration: ", loginCompleter)
		commands.DisableClusterConfiguration(api, configuration)
	case "delete cluster":
		cluster := prompt.Input("cluster: ", loginCompleter)
		deleteCluster(cluster)
	case "delete configuration":
		configuration := prompt.Input("configuration: ", loginCompleter)
		deleteClusterConfiguration(configuration)
	case "delete profile":
		profile := prompt.Input("profile: ", loginCompleter)
		deleteConfigurationProfile(profile)
	case "delete trigger":
		trigger := prompt.Input("trigger: ", loginCompleter)
		deleteTrigger(trigger)
	case "activate must-gather":
		fallthrough
	case "activate trigger":
		trigger := prompt.Input("trigger: ", loginCompleter)
		activateTrigger(trigger)
	case "deactivate must-gather":
		fallthrough
	case "deactivate trigger":
		trigger := prompt.Input("trigger: ", loginCompleter)
		deactivateTrigger(trigger)
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
		commands.PrintHelp()
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
		{Text: "list", Description: "list resources (clusters, profiles, configurations, triggers)"},
		{Text: "describe", Description: "describe the selected resource"},
		{Text: "request", Description: "request selected operation to be performed"},
		{Text: "add", Description: "add resource (cluster, profile, configuration, trigger)"},
		{Text: "new", Description: "alias for add"},
		{Text: "delete", Description: "delete resource (configuration, trigger)"},
		{Text: "enable", Description: "enable selected cluster profile"},
		{Text: "disable", Description: "disable selected cluster profile"},
		{Text: "activate", Description: "activate resource (trigger)"},
		{Text: "deactivate", Description: "deactivate resource (trigger)"},
		{Text: "version", Description: "prints the build information for CLI executable"},
	}

	secondWord := make(map[string][]prompt.Suggest)

	// list operations
	secondWord["list"] = []prompt.Suggest{
		{Text: "clusters", Description: "show list of all clusters available"},
		{Text: "profiles", Description: "show list of all configuration profiles"},
		{Text: "configurations", Description: "show list all cluster configurations"},
		{Text: "must-gather", Description: "show list all must-gathers"},
		{Text: "triggers", Description: "show list with all must-gather triggers"},
	}
	// add operations
	secondWord["add"] = []prompt.Suggest{
		{Text: "cluster", Description: "add/register new cluster"},
		{Text: "profile", Description: "add new configuration profile"},
		{Text: "configuration", Description: "add new cluster configuration"},
		{Text: "trigger", Description: "add new must-gather trigger"},
	}
	// new operations (aliases for add)
	secondWord["new"] = []prompt.Suggest{
		{Text: "cluster", Description: "add/register new cluster"},
		{Text: "profile", Description: "add new configuration profile"},
		{Text: "configuration", Description: "add new cluster configuration"},
		{Text: "trigger", Description: "add new must-gather trigger"},
	}
	// request operations
	secondWord["request"] = []prompt.Suggest{
		{Text: "must-gather", Description: "request must-gather"},
	}
	// delete operations
	secondWord["delete"] = []prompt.Suggest{
		{Text: "cluster", Description: "delete cluster and its configuration"},
		{Text: "profile", Description: "delete configuration profile"},
		{Text: "configuration", Description: "delete cluster configuration"},
		{Text: "trigger", Description: "delete trigger"},
	}
	// descripbe operations
	secondWord["describe"] = []prompt.Suggest{
		{Text: "profile", Description: "describe selected configuration profile"},
		{Text: "configuration", Description: "describe configuration for selected cluster"},
		{Text: "trigger", Description: "describe selected must-gather trigger"},
		{Text: "must-gather", Description: "describe selected must-gather trigger"},
	}
	// activate operations
	secondWord["activate"] = []prompt.Suggest{
		{Text: "trigger", Description: "activate selected must-gather trigger"},
		{Text: "must-gather", Description: "activate selected must-gather"},
	}
	// deactivate operations
	secondWord["deactivate"] = []prompt.Suggest{
		{Text: "trigger", Description: "deactivate selected must-gather trigger"},
		{Text: "must-gather", Description: "deactivate selected must-gather"},
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
	// read configuration first
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	controllerUrl = viper.GetString("CONTROLLER_URL")
	p := prompt.New(executor, completer)
	api = restapi.NewRestApi(controllerUrl)

	p.Run()
}
