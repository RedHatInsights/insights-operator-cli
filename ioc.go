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
	"github.com/logrusorgru/aurora"
	"github.com/redhatinsighs/insights-operator-cli/commands"
	"github.com/redhatinsighs/insights-operator-cli/restapi"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// BuildVersion contains the major.minor version of the CLI client
var BuildVersion string = "*not set*"

// BuildTime contains timestamp when the CLI client has been built
var BuildTime string = "*not set*"

var controllerURL string
var username string
var password string
var files []prompt.Suggest
var api restapi.Api

func tryToLogin(username string, password string) {
	fmt.Println(aurora.Blue("\nDone"))
}

func listOfConfigurations(filter string) {
	// TODO: filter in query?
	configurations, err := api.ReadListOfConfigurations()
	if err != nil {
		fmt.Println(aurora.Red("Error reading list of configurations"))
		fmt.Println(err)
		return
	}

	fmt.Println(aurora.Magenta("List of configuration for all clusters"))
	fmt.Printf("%4s %4s %4s    %-20s %-20s %-10s %-12s %s\n", "#", "ID", "Profile", "Cluster", "Changed at", "Changed by", "Active", "Reason")
	for i, configuration := range configurations {
		// poor man's filtering
		if strings.Contains(configuration.Cluster, filter) {
			var active aurora.Value
			if configuration.Active == "1" {
				active = aurora.Green("yes")
			} else {
				active = aurora.Red("no")
			}
			changedAt := configuration.ChangedAt[0:19]
			fmt.Printf("%4d %4d %4s       %-20s %-20s %-10s %-12s %s\n", i, configuration.ID, configuration.Configuration, configuration.Cluster, changedAt, configuration.ChangedBy, active, configuration.Reason)
		}
	}
}

func describeProfile(profileID string) {
	profile, err := api.ReadConfigurationProfile(profileID)
	if err != nil {
		fmt.Println(aurora.Red("Error reading configuration profile"))
		fmt.Println(err)
		return
	}

	fmt.Println(aurora.Magenta("Configuration profile"))
	fmt.Println(profile.Configuration)
}

func describeConfiguration(clusterID string) {
	configuration, err := api.ReadClusterConfigurationById(clusterID)
	if err != nil {
		fmt.Println(aurora.Red("Error reading cluster configuration"))
		fmt.Println(err)
		return
	}

	fmt.Println(aurora.Magenta("Configuration for cluster " + clusterID))
	fmt.Println(*configuration)
}

func proceedQuestion(question string) bool {
	fmt.Println(aurora.Red(question))
	proceed := prompt.Input("proceed? [y/n] ", loginCompleter)
	if proceed != "y" {
		fmt.Println(aurora.Blue("cancelled"))
		return false
	}
	return true
}

func deleteCluster(clusterID string) {
	if !proceedQuestion("All cluster configurations will be deleted") {
		return
	}

	err := api.DeleteCluster(clusterID)
	if err != nil {
		fmt.Println(aurora.Red("Error communicating with the service"))
		fmt.Println(err)
		return
	}

	// everything's ok
	fmt.Println(aurora.Blue("Cluster "+clusterID+" has been"), aurora.Red("deleted"))
}

func deleteClusterConfiguration(configurationID string) {
	err := api.DeleteClusterConfiguration(configurationID)
	if err != nil {
		fmt.Println(aurora.Red("Error communicating with the service"))
		fmt.Println(err)
		return
	}

	// everything's ok, configuration has been deleted
	fmt.Println(aurora.Blue("Configuration "+configurationID+" has been "), aurora.Red("deleted"))
}

func deleteConfigurationProfile(profileID string) {
	if !proceedQuestion("All configurations based on this profile will be deleted") {
		return
	}

	err := api.DeleteConfigurationProfile(profileID)
	if err != nil {
		fmt.Println(aurora.Red("Error communicating with the service"))
		fmt.Println(err)
		return
	}

	// everything's ok, configuration profile has been deleted
	fmt.Println(aurora.Blue("Configuration profile "+profileID+" has been "), aurora.Red("deleted"))
}

func addCluster() {
	id := prompt.Input("ID: ", loginCompleter)
	if id == "" {
		fmt.Println(aurora.Red("Cancelled"))
		return
	}

	name := prompt.Input("name: ", loginCompleter)
	if name == "" {
		fmt.Println(aurora.Red("Cancelled"))
		return
	}

	err := api.AddCluster(id, name)
	if err != nil {
		fmt.Println(aurora.Red("Error communicating with the service"))
		fmt.Println(err)
		return
	}

	// everything's ok, controller has been registered
	fmt.Println(aurora.Blue("Controller has been registered"))
}

func addProfile() {
	if username == "" {
		fmt.Println(aurora.Red("Not logged in"))
		return
	}

	description := prompt.Input("description: ", loginCompleter)
	if description == "" {
		fmt.Println(aurora.Red("Cancelled"))
		return
	}

	// TODO: make the directory fully configurable
	err := fillInConfigurationList("configurations")
	if err != nil {
		fmt.Println(aurora.Red("Cannot read any configuration file"))
		fmt.Println(err)
	}

	configurationFileName := prompt.Input("configuration file (TAB to complete): ", configFileCompleter)
	if configurationFileName == "" {
		fmt.Println(aurora.Red("Cancelled"))
		return
	}

	// TODO: make the directory fully configurable
	configuration, err := ioutil.ReadFile("configurations/" + configurationFileName)
	if err != nil {
		fmt.Println(aurora.Red("Cannot read configuration file"))
		fmt.Println(err)
	}

	err = api.AddConfigurationProfile(username, description, configuration)
	if err != nil {
		fmt.Println(aurora.Red("Error communicating with the service"))
		fmt.Println(err)
		return
	}

	// everything's ok, configuration profile has been created
	fmt.Println(aurora.Blue("Configuration profile has been created"))
}

func addClusterConfiguration() {
	if username == "" {
		fmt.Println(aurora.Red("Not logged in"))
		return
	}

	cluster := prompt.Input("cluster: ", loginCompleter)
	if cluster == "" {
		fmt.Println(aurora.Red("Cancelled"))
		return
	}

	reason := prompt.Input("reason: ", loginCompleter)
	if reason == "" {
		fmt.Println(aurora.Red("Cancelled"))
		return
	}

	description := prompt.Input("description: ", loginCompleter)
	if description == "" {
		fmt.Println(aurora.Red("Cancelled"))
		return
	}

	// TODO: make the directory fully configurable
	err := fillInConfigurationList("configurations")
	if err != nil {
		fmt.Println(aurora.Red("Cannot read any configuration file"))
		fmt.Println(err)
	}

	configurationFileName := prompt.Input("configuration file (TAB to complete): ", configFileCompleter)
	if configurationFileName == "" {
		fmt.Println(aurora.Red("Cancelled"))
		return
	}

	// TODO: make the directory fully configurable
	configuration, err := ioutil.ReadFile("configurations/" + configurationFileName)
	if err != nil {
		fmt.Println(aurora.Red("Cannot read configuration file"))
		fmt.Println(err)
	}

	err = api.AddClusterConfiguration(username, cluster, reason, description, configuration)
	if err != nil {
		fmt.Println(aurora.Red("Error communicating with the service"))
		fmt.Println(err)
		return
	}

	// everything's ok, configuration has been created
	fmt.Println(aurora.Blue("Configuration has been created"))
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
		fmt.Println(aurora.Red("Not logged in"))
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
	}

	// everything's ok, trigger has been created
	fmt.Println(aurora.Blue("Trigger has been created"))
}

func describeTrigger(triggerID string) {
	trigger, err := api.ReadTriggerById(triggerID)
	if err != nil {
		fmt.Println(aurora.Red("Error reading selected trigger"))
		fmt.Println(err)
		return
	}

	var active aurora.Value
	if trigger.Active == 1 {
		active = aurora.Green("yes")
	} else {
		active = aurora.Red("no")
	}

	triggeredAt := trigger.TriggeredAt[0:19]
	ackedAt := trigger.AckedAt[0:19]

	var ttype aurora.Value
	if trigger.Type == "must-gather" {
		ttype = aurora.Blue(trigger.Type)
	} else {
		ttype = aurora.Magenta(trigger.Type)
	}

	fmt.Println(aurora.Magenta("Trigger info"))
	fmt.Printf("ID:            %d\n", trigger.ID)
	fmt.Printf("Type:          %s\n", ttype)
	fmt.Printf("Cluster:       %s\n", trigger.Cluster)
	fmt.Printf("Triggered at:  %s\n", triggeredAt)
	fmt.Printf("Triggered by:  %s\n", trigger.TriggeredBy)
	fmt.Printf("Active:        %s\n", active)
	fmt.Printf("Acked at:      %s\n", ackedAt)
}

func deleteTrigger(triggerID string) {
	err := api.DeleteTrigger(triggerID)
	if err != nil {
		fmt.Println(aurora.Red("Error communicating with the service"))
		fmt.Println(err)
		return
	}

	// everything's ok, trigger has been deleted
	fmt.Println(aurora.Blue("Trigger "+triggerID+" has been"), aurora.Red("deleted"))
}

func activateTrigger(triggerID string) {
	err := api.ActivateTrigger(triggerID)
	if err != nil {
		fmt.Println(aurora.Red("Error communicating with the service"))
		fmt.Println(err)
		return
	}

	// everything's ok, trigger has been activated
	fmt.Println(aurora.Blue("Trigger "+triggerID+" has been"), aurora.Green("activated"))
}

func deactivateTrigger(triggerID string) {
	err := api.DeactivateTrigger(triggerID)
	if err != nil {
		fmt.Println(aurora.Red("Error communicating with the service"))
		fmt.Println(err)
		return
	}

	// everything's ok, trigger has been deactivated
	fmt.Println(aurora.Blue("Trigger "+triggerID+" has been"), aurora.Green("deactivated"))
}

func loginCompleter(in prompt.Document) []prompt.Suggest {
	return nil
}

func configFileCompleter(in prompt.Document) []prompt.Suggest {
	return prompt.FilterHasPrefix(files, in.Text, true)
}

func printVersion() {
	fmt.Println(aurora.Blue("Insights operator CLI client "), "version", aurora.Yellow(BuildVersion), "compiled", aurora.Yellow(BuildTime))
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
			fmt.Println(aurora.Red("Password is not set"))
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
		fmt.Println(aurora.Magenta("Quitting"))
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

	emptySuggest := []prompt.Suggest{}

	blocks := strings.Split(in.TextBeforeCursor(), " ")

	if len(blocks) == 2 {
		sec, ok := secondWord[blocks[0]]
		if ok {
			return prompt.FilterHasPrefix(sec, blocks[1], true)
		}
		// second word is not known
		return emptySuggest
	}
	if in.GetWordBeforeCursor() == "" {
		return nil
	}

	// commands consisting of just one word
	return prompt.FilterHasPrefix(firstWord, blocks[0], true)
}

func main() {
	// read configuration first
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	controllerURL = viper.GetString("CONTROLLER_URL")
	p := prompt.New(executor, completer)
	api = restapi.NewRestApi(controllerURL)

	p.Run()
}
