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
	"strings"
)

// BuildVersion contains the major.minor version of the CLI client
var BuildVersion string = "*not set*"

// BuildTime contains timestamp when the CLI client has been built
var BuildTime string = "*not set*"

var username string
var password string
var api restapi.Api

func tryToLogin(username string, password string) {
	fmt.Println(aurora.Blue("\nDone"))
}

func printVersion() {
	fmt.Println(aurora.Blue("Insights operator CLI client "), "version", aurora.Yellow(BuildVersion), "compiled", aurora.Yellow(BuildTime))
}

func login() {
	username = prompt.Input("login: ", commands.LoginCompleter)
	fmt.Print("password: ")
	p, err := terminal.ReadPassword(0)
	if err != nil {
		fmt.Println(aurora.Red("Password is not set"))
	} else {
		password = string(p)
		tryToLogin(username, password)
	}
}

type commandWithParam struct {
	prefix  string
	handler func(restapi.Api, string)
}

var commandsWithParam = []commandWithParam{
	{"describe profile ", commands.DescribeProfile},
	{"describe configuration ", commands.DescribeConfiguration},
	{"describe trigger ", commands.DescribeTrigger},
	{"enable configuration ", commands.EnableClusterConfiguration},
	{"disable configuration ", commands.DisableClusterConfiguration},
	{"list configurations ", commands.ListOfConfigurations},
	{"delete cluster ", commands.DeleteCluster},
	{"delete configuration ", commands.DeleteClusterConfiguration},
	{"delete trigger ", commands.DeleteTrigger},
	{"activate must-gather ", commands.ActivateTrigger},
	{"activate trigger ", commands.ActivateTrigger},
	{"deactivate must-gather ", commands.DeactivateTrigger},
	{"deactivate trigger ", commands.DeactivateTrigger},
}

func executor(t string) {
	blocks := strings.Split(t, " ")

	// commands with variable parts
	for _, command := range commandsWithParam {
		if strings.HasPrefix(t, command.prefix) {
			command.handler(api, blocks[2])
			return
		}
	}
	// no match? try commands without variable parts
	executeFixedCommand(t)
}

type simpleCommand struct {
	prefix  string
	handler func()
}

var simpleCommands = []simpleCommand{
	{"login", login},
	{"bye", commands.Quit},
	{"exit", commands.Quit},
	{"quit", commands.Quit},
	{"?", commands.PrintHelp},
	{"help", commands.PrintHelp},
	{"version", printVersion},
}

type commandWithApiParam struct {
	prefix  string
	handler func(restapi.Api)
}

var commandsWithApiParam = []commandWithApiParam{
	{"list must-gather", commands.ListOfTriggers},
	{"list triggers", commands.ListOfTriggers},
	{"list clusters", commands.ListOfClusters},
	{"list profiles", commands.ListOfProfiles},
	{"new cluster", commands.AddCluster},
}

func executeFixedCommand(t string) {
	// simple commands without parameters
	for _, command := range simpleCommands {
		if strings.HasPrefix(t, command.prefix) {
			command.handler()
			return
		}
	}
	// fixed commands with API as param
	for _, command := range commandsWithApiParam {
		if strings.HasPrefix(t, command.prefix) {
			command.handler(api)
			return
		}
	}
	switch t {
	case "list configurations":
		commands.ListOfConfigurations(api, "")
	case "add cluster":
		fallthrough
	case "add profile":
		fallthrough
	case "new profile":
		commands.AddConfigurationProfile(api, username)
	case "add configuration":
		fallthrough
	case "new configuration":
		commands.AddClusterConfiguration(api, username)
	case "request must-gather":
		fallthrough
	case "add trigger":
		fallthrough
	case "new trigger":
		commands.AddTrigger(api, username)
	case "describe profile":
		profile := prompt.Input("profile: ", commands.LoginCompleter)
		commands.DescribeProfile(api, profile)
	case "describe configuration":
		configuration := prompt.Input("configuration: ", commands.LoginCompleter)
		commands.DescribeConfiguration(api, configuration)
	case "describe must-gather":
		fallthrough
	case "describe trigger":
		trigger := prompt.Input("trigger: ", commands.LoginCompleter)
		commands.DescribeTrigger(api, trigger)
	case "enable configuration":
		configuration := prompt.Input("configuration: ", commands.LoginCompleter)
		commands.EnableClusterConfiguration(api, configuration)
	case "disable configuration":
		configuration := prompt.Input("configuration: ", commands.LoginCompleter)
		commands.DisableClusterConfiguration(api, configuration)
	case "delete cluster":
		cluster := prompt.Input("cluster: ", commands.LoginCompleter)
		commands.DeleteCluster(api, cluster)
	case "delete configuration":
		configuration := prompt.Input("configuration: ", commands.LoginCompleter)
		commands.DeleteClusterConfiguration(api, configuration)
	case "delete profile":
		profile := prompt.Input("profile: ", commands.LoginCompleter)
		commands.DeleteConfigurationProfile(api, profile)
	case "delete trigger":
		trigger := prompt.Input("trigger: ", commands.LoginCompleter)
		commands.DeleteTrigger(api, trigger)
	case "activate must-gather":
		fallthrough
	case "activate trigger":
		trigger := prompt.Input("trigger: ", commands.LoginCompleter)
		commands.ActivateTrigger(api, trigger)
	case "deactivate must-gather":
		fallthrough
	case "deactivate trigger":
		trigger := prompt.Input("trigger: ", commands.LoginCompleter)
		commands.DeactivateTrigger(api, trigger)
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
	// enable operations
	secondWord["enable"] = []prompt.Suggest{
		{Text: "configuration", Description: "enable cluster configuration"},
	}
	// disable operations
	secondWord["disable"] = []prompt.Suggest{
		{Text: "configuration", Description: "disable cluster configuration"},
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
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	controllerURL := viper.GetString("CONTROLLER_URL")
	p := prompt.New(executor, completer)
	api = restapi.NewRestApi(controllerURL)

	p.Run()
}
