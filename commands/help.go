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

package commands

// Generated documentation is available at:
// https://pkg.go.dev/github.com/RedHatInsights/insights-operator-cli/commands
//
// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-cli/packages/commands/help.html

import (
	"fmt"
)

const (
	commandAlias = "alias for previous command"
)

// PrintHelp function can be used to display help on (color) terminal.
// Colorization is optional and depends on configuration.
func PrintHelp() {
	fmt.Println(colorizer.Magenta("HELP:"))
	fmt.Println()

	// cluster manipulation commands
	fmt.Println(colorizer.Blue("Cluster operations:        "))
	fmt.Println(colorizer.Yellow("list clusters            "), "list all clusters known to the service")
	fmt.Println(colorizer.Yellow("delete cluster ##        "), "delete selected cluster")
	fmt.Println(colorizer.Yellow("add cluster              "), "create new cluster")
	fmt.Println(colorizer.Yellow("new cluster              "), commandAlias)
	fmt.Println()

	// configuration profiles manipulation commands
	fmt.Println(colorizer.Blue("Configuration profiles:    "))
	fmt.Println(colorizer.Yellow("list profiles            "), "list all profiles known to the service")
	fmt.Println(colorizer.Yellow("describe profile ##      "), "describe profile selected by its ID")
	fmt.Println(colorizer.Yellow("add profile              "), "create new configuration profile")
	fmt.Println(colorizer.Yellow("delete profile ##        "), "delete profile selected by its ID")
	fmt.Println()

	// cluster configuraration commands
	fmt.Println(colorizer.Blue("Cluster configurations:    "))
	fmt.Println(colorizer.Yellow("list configurations      "), "list all configurations known to the service")
	fmt.Println(colorizer.Yellow("describe configuration ##"), "describe cluster configuration selected by its ID")
	fmt.Println(colorizer.Yellow("add configuration        "), "add new configuration")
	fmt.Println(colorizer.Yellow("new configuration        "), commandAlias)
	fmt.Println(colorizer.Yellow("enable configuration ##  "), "enable cluster configuration selected by its ID")
	fmt.Println(colorizer.Yellow("disable configuration ## "), "disable cluster configuration selected by its ID")
	fmt.Println(colorizer.Yellow("delete configuration ##  "), "delete configuration selected by its ID")
	fmt.Println()

	// must-gather triggering related commands
	fmt.Println(colorizer.Blue("Must-gather trigger:       "))
	fmt.Println(colorizer.Yellow("list triggers            "), "list all triggers")
	fmt.Println(colorizer.Yellow("describe trigger ##      "), "describe trigger selected by its ID")
	fmt.Println(colorizer.Yellow("add trigger              "), "add new trigger")
	fmt.Println(colorizer.Yellow("new trigger              "), commandAlias)
	fmt.Println(colorizer.Yellow("activate trigger ##      "), "activate trigger selected by its ID")
	fmt.Println(colorizer.Yellow("deactivate trigger ##    "), "deactivate trigger selected by its ID")
	fmt.Println(colorizer.Yellow("delete trigger ##        "), "delete trigger")
	fmt.Println()

	// other commands
	fmt.Println(colorizer.Blue("Other commands:"))
	fmt.Println(colorizer.Yellow("version                  "), "print version information")
	fmt.Println(colorizer.Yellow("quit                     "), "quit the application")
	fmt.Println(colorizer.Yellow("exit                     "), commandAlias)
	fmt.Println(colorizer.Yellow("bye                      "), commandAlias)
	fmt.Println(colorizer.Yellow("help                     "), "this help")
	fmt.Println()
}
