/*
Copyright © 2019, 2020, 2021 Red Hat, Inc.

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

// Implementation of command-line client for the insights operator
// instrumentation service.
package main

// Generated documentation is available at:
// https://pkg.go.dev/github.com/RedHatInsights/insights-operator-cli
//
// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-cli/packages/ioc.html

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/RedHatInsights/insights-operator-cli/commands"
	"github.com/RedHatInsights/insights-operator-cli/restapi"
	"github.com/c-bata/go-prompt"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
)

// prompts
const (
	profilePrompt = "profile: "
)

// BuildVersion contains the major.minor version of the CLI client
var BuildVersion string = "*not set*"

// BuildTime contains timestamp when the CLI client has been built
var BuildTime string = "*not set*"

// Configuration represent insights operator CLI configuration
type Configuration struct {
	// enable or disable asking for confirmation for selected actions
	// (like all delete commands)
	askForConfirmation *bool

	// enable colors usage on CLI interface
	colors *bool

	// enable or disable Tab-completion
	useCompleter *bool
}

// configuration represents current CLI configuration
var configuration Configuration

// username used to access REST API
var username string

// password used to access REST API
var password string

// api represents instance of REST API
var api restapi.API

// colorizer represents implementation of interface used to provide (display)
// color output on terminal
var colorizer aurora.Aurora

// tryToLogin tries to login to service via REST API
func tryToLogin(username string, password string) {
	fmt.Println(colorizer.Blue("\nDone"))
}

// printVersion displays version of Insights operator CLI client
func printVersion() {
	fmt.Println(colorizer.Blue("Insights operator CLI client "),
		"version", colorizer.Yellow(BuildVersion), "compiled",
		colorizer.Yellow(BuildTime))
}

// login prompts for username and password and then tries to login to the
// service via REST API
func login() {
	username = prompt.Input("login: ", commands.LoginCompleter)
	fmt.Print("password: ")
	p, err := terminal.ReadPassword(0)
	if err != nil {
		fmt.Println(colorizer.Red("Password is not set"))
	} else {
		password = string(p)
		tryToLogin(username, password)
	}
}

// commandWithParam represents any commands than needs to be followed by
// parameter
type commandWithParam struct {
	prefix  string
	handler func(restapi.API, string)
}

// commandsWithParam is a list of all supported commands with parameters
var commandsWithParam = []commandWithParam{
	{"describe profile ", commands.DescribeProfile},
	{"describe configuration ", commands.DescribeConfiguration},
	{"describe trigger ", commands.DescribeTrigger},
	{"enable configuration ", commands.EnableClusterConfiguration},
	{"disable configuration ", commands.DisableClusterConfiguration},
	{"list configurations ", commands.ListOfConfigurations},
	{"delete cluster ", commands.DeleteClusterNoConfirm},
	{"delete profile ", commands.DeleteConfigurationProfileNoConfirm},
	{"delete configuration ", commands.DeleteClusterConfiguration},
	{"delete trigger ", commands.DeleteTrigger},
	{"activate must-gather ", commands.ActivateTrigger},
	{"activate trigger ", commands.ActivateTrigger},
	{"deactivate must-gather ", commands.DeactivateTrigger},
	{"deactivate trigger ", commands.DeactivateTrigger},
	{"add cluster ", commands.AddCluster},
	{"new cluster ", commands.AddCluster},
}

// executor tries to call the command specified on command line
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

// simpleCommand represents any command without parameter
type simpleCommand struct {
	prefix  string
	handler func()
}

// simpleCommands is a list of all supported commands without parameters
var simpleCommands = []simpleCommand{
	{"login", login},
	{"bye", commands.Quit},
	{"exit", commands.Quit},
	{"quit", commands.Quit},
	{"?", commands.PrintHelp},
	{"help", commands.PrintHelp},
	{"version", printVersion},
	{"license", commands.PrintLicense},
	{"authors", commands.PrintAuthors},
}

// commandsWithAPIParam represents any command with REST API parameter
type commandWithAPIParam struct {
	prefix  string
	handler func(restapi.API)
}

// commandsWithAPIParam is a list of all supported commands with REST API
// parameter
var commandsWithAPIParam = []commandWithAPIParam{
	{"list must-gather", commands.ListOfTriggers},
	{"list triggers", commands.ListOfTriggers},
	{"list clusters", commands.ListOfClusters},
	{"list profiles", commands.ListOfProfiles},
}

// configurationPrompt function displays input prompt for configuration file
// (with completer)
func configurationPrompt() string {
	return prompt.Input("configuration: ", commands.LoginCompleter)
}

// triggerPrompt function displays input prompt for trigger configuration file
// (with completer)
func triggerPrompt() string {
	return prompt.Input("trigger: ", commands.LoginCompleter)
}

// executeFixedCommand tries to execute command stored in argument
func executeFixedCommand(t string) {
	// simple commands without parameters
	for _, command := range simpleCommands {
		if strings.HasPrefix(t, command.prefix) {
			command.handler()
			return
		}
	}
	// fixed commands with API as param
	for _, command := range commandsWithAPIParam {
		if strings.HasPrefix(t, command.prefix) {
			command.handler(api)
			return
		}
	}
	// other commands, usually ones constructed from two words
	switch t {
	case "list configurations":
		commands.ListOfConfigurations(api, "")
	case "add cluster":
		clusterName := prompt.Input("clusterName: ",
			commands.LoginCompleter)
		commands.AddCluster(api, clusterName)
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
		profile := prompt.Input(profilePrompt, commands.LoginCompleter)
		commands.DescribeProfile(api, profile)
	case "describe configuration":
		configuration := configurationPrompt()
		commands.DescribeConfiguration(api, configuration)
	case "describe must-gather":
		fallthrough
	case "describe trigger":
		trigger := triggerPrompt()
		commands.DescribeTrigger(api, trigger)
	case "enable configuration":
		configuration := configurationPrompt()
		commands.EnableClusterConfiguration(api, configuration)
	case "disable configuration":
		configuration := configurationPrompt()
		commands.DisableClusterConfiguration(api, configuration)
	case "delete cluster":
		cluster := prompt.Input("cluster to delete: ",
			commands.LoginCompleter)
		commands.DeleteCluster(api, cluster,
			*configuration.askForConfirmation)
	case "delete configuration":
		configuration := configurationPrompt()
		commands.DeleteClusterConfiguration(api, configuration)
	case "delete profile":
		profile := prompt.Input(profilePrompt, commands.LoginCompleter)
		commands.DeleteConfigurationProfile(api, profile,
			*configuration.askForConfirmation)
	case "delete trigger":
		trigger := triggerPrompt()
		commands.DeleteTrigger(api, trigger)
	case "activate must-gather":
		fallthrough
	case "activate trigger":
		trigger := triggerPrompt()
		commands.ActivateTrigger(api, trigger)
	case "deactivate must-gather":
		fallthrough
	case "deactivate trigger":
		trigger := triggerPrompt()
		commands.DeactivateTrigger(api, trigger)
	default:
		fmt.Println("Command not found")
	}
}

// completer function is called by Aurora to autocomplete command typed by user
// on command line.
func completer(in prompt.Document) []prompt.Suggest {
	// suggestions for the first word for all commands
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
		{Text: "copyright", Description: "displays copyright notice"},
		{Text: "license", Description: "displays license used by this project"},
		{Text: "authors", Description: "displays list of authors"},
	}

	// map with autocompletion for commands consisting from two words
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

	// handle commands with two words or one word + parameter
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

// readConfiguration function reads configuration from configuration file and
// via CLI flags.
func readConfiguration(filename string) (Configuration, error) {
	var config Configuration

	// read configuration first
	viper.SetConfigName(filename)
	viper.AddConfigPath(".")

	// try to read configuration and check for possible errors
	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	// parse command line arguments and flags
	config.colors = flag.Bool("colors", true, "enable or disable colors")
	config.useCompleter = flag.Bool("completer", true,
		"enable or disable command line completer")
	config.askForConfirmation = flag.Bool("confirmation", true,
		"enable or disable asking for confirmation for selected actions (like delete)")
	flag.Parse()

	return config, nil
}

// main function represents entry point to CLI client called right after the
// process is started.
func main() {
	// read configuration
	configuration, err := readConfiguration("config")
	if err != nil {
		panic(err)
	}

	// initialize colorizers
	colorizer = aurora.NewAurora(*configuration.colors)
	commands.SetColorizer(colorizer)

	// initialize REST API connection to service
	controllerURL := viper.GetString("CONTROLLER_URL")
	api = restapi.NewRestAPI(controllerURL)

	// start the command line
	if *configuration.useCompleter {
		// command line prompt with autocompleter
		p := prompt.New(executor, completer)
		p.Run()
	} else {
		// command line prompt without autocompleter
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("> ")
		for scanner.Scan() {
			line := scanner.Text()
			executor(line)
			fmt.Print("> ")
		}
	}
}
