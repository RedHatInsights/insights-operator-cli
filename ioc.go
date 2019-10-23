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
	"encoding/json"
	"fmt"
	"github.com/c-bata/go-prompt"
	. "github.com/logrusorgru/aurora"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const API_PREFIX = "/api/v1/"

var controllerUrl string
var username string
var password string

func tryToLogin(username string, password string) {
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

func readListOfClusters(controllerUrl string, apiPrefix string) ([]Cluster, error) {
	clusters := []Cluster{}

	url := controllerUrl + apiPrefix + "client/cluster"
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Expected HTTP status 200 OK, got %d", response.StatusCode)
	}

	body, readErr := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if readErr != nil {
		return nil, fmt.Errorf("Unable to read response body")
	}

	err = json.Unmarshal(body, &clusters)
	if err != nil {
		return nil, err
	}
	return clusters, nil
}

func readListOfConfigurationProfiles(controllerUrl string, apiPrefix string) ([]ConfigurationProfile, error) {
	profiles := []ConfigurationProfile{}

	url := controllerUrl + apiPrefix + "client/profile"
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

	err = json.Unmarshal(body, &profiles)
	if err != nil {
		return nil, err
	}
	return profiles, nil
}

func readListOfConfigurations(controllerUrl string, apiPrefix string) ([]ClusterConfiguration, error) {
	configurations := []ClusterConfiguration{}

	url := controllerUrl + apiPrefix + "client/configuration"
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
		return configurations, fmt.Errorf("Unable to read response body")
	}

	err = json.Unmarshal(body, &configurations)
	if err != nil {
		return nil, err
	}
	return configurations, nil
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
	fmt.Printf("%4s %4s %-20s %-20s %s\n", "#", "ID", "ChangedAt", "ChangedBy", "Description")
	for i, profile := range profiles {
		fmt.Printf("%4d %4d %-20s %-20s %-s\n", i, profile.Id, profile.ChangedAt, profile.ChangedBy, profile.Description)
	}
}

func listOfConfigurations() {
	configurations, err := readListOfConfigurations(controllerUrl, API_PREFIX)
	if err != nil {
		fmt.Println(Red("Error reading list of configurations"))
		fmt.Println(err)
		return
	}

	fmt.Println(Magenta("List of configurations"))
	fmt.Printf("%4s %4s %-20s %-20s %-10s %-12s %s\n", "#", "ID", "Cluster", "ChangedAt", "ChangedBy", "Active", "Reason")
	for i, configuration := range configurations {
		fmt.Printf("%4d %4d %-20s %-20s %-10s %-12s %s\n", i, configuration.Id, configuration.Cluster, configuration.ChangedAt, configuration.ChangedBy, configuration.Active, configuration.Reason)
	}
}

func printHelp() {
	fmt.Println("HELP:\nexit\nquit")
}

func loginCompleter(in prompt.Document) []prompt.Suggest {
	return nil
}

func executor(t string) {
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
		listOfConfigurations()
	case "bye":
		fallthrough
	case "exit":
		fallthrough
	case "quit":
		fmt.Println(Magenta("Quitting"))
		os.Exit(0)
	case "help":
		printHelp()
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
	}

	secondWord := make(map[string][]prompt.Suggest)

	// list operations
	secondWord["list"] = []prompt.Suggest{
		{Text: "clusters", Description: "show list of all clusters available"},
		{Text: "profiles", Description: "show list of all configuration profiles"},
		{Text: "configurations", Description: "show list all cluster configurations"},
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
