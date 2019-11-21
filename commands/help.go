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

package commands

import (
	"fmt"
	"github.com/logrusorgru/aurora"
)

// PrintHelp can be used to display help on (color) terminal.
func PrintHelp() {
	fmt.Println(aurora.Magenta("HELP:"))
	fmt.Println()
	fmt.Println(aurora.Blue("Cluster operations:        "))
	fmt.Println(aurora.Yellow("list clusters            "), "list all clusters known to the service")
	fmt.Println(aurora.Yellow("delete cluster ##        "), "delete selected cluster")
	fmt.Println(aurora.Yellow("add cluster              "), "create new cluster")
	fmt.Println(aurora.Yellow("new cluster              "), "alias for previous command")
	fmt.Println()
	fmt.Println(aurora.Blue("Configuration profiles:    "))
	fmt.Println(aurora.Yellow("list profiles            "), "list all profiles known to the service")
	fmt.Println(aurora.Yellow("describe profile ##      "), "describe profile selected by its ID")
	fmt.Println(aurora.Yellow("delete profile ##        "), "delete profile selected by its ID")
	fmt.Println()
	fmt.Println(aurora.Blue("Cluster configurations:    "))
	fmt.Println(aurora.Yellow("list configurations      "), "list all configurations known to the service")
	fmt.Println(aurora.Yellow("describe configuration ##"), "describe cluster configuration selected by its ID")
	fmt.Println(aurora.Yellow("add configuration        "), "add new configuration")
	fmt.Println(aurora.Yellow("new configuration        "), "alias for previous command")
	fmt.Println(aurora.Yellow("enable configuration ##  "), "enable cluster configuration selected by its ID")
	fmt.Println(aurora.Yellow("disable configuration ## "), "disable cluster configuration selected by its ID")
	fmt.Println(aurora.Yellow("delete configuration ##  "), "delete configuration selected by its ID")
	fmt.Println()
	fmt.Println(aurora.Blue("Must-gather trigger:       "))
	fmt.Println(aurora.Yellow("list triggers            "), "list all triggers")
	fmt.Println(aurora.Yellow("describe trigger ##      "), "describe trigger selected by its ID")
	fmt.Println(aurora.Yellow("add trigger              "), "add new trigger")
	fmt.Println(aurora.Yellow("new trigger              "), "alias for previous command")
	fmt.Println(aurora.Yellow("activate trigger ##      "), "activate trigger selected by its ID")
	fmt.Println(aurora.Yellow("deactivate trigger ##    "), "deactivate trigger selected by its ID")
	fmt.Println(aurora.Yellow("delete trigger           "), "delete trigger")
	fmt.Println()
	fmt.Println(aurora.Blue("Other commands:"))
	fmt.Println(aurora.Yellow("version                  "), "print version information")
	fmt.Println(aurora.Yellow("quit                     "), "quit the application")
	fmt.Println(aurora.Yellow("exit                     "), "dtto")
	fmt.Println(aurora.Yellow("bye                      "), "dtto")
	fmt.Println(aurora.Yellow("help                     "), "this help")
	fmt.Println()
}
