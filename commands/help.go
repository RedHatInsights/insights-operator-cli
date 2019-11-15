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
	. "github.com/logrusorgru/aurora"
)

func PrintHelp() {
	fmt.Println(Magenta("HELP:"))
	fmt.Println()
	fmt.Println(Blue("Cluster operations:        "))
	fmt.Println(Yellow("list clusters            "), "list all clusters known to the service")
	fmt.Println(Yellow("delete cluster ##        "), "delete selected cluster")
	fmt.Println(Yellow("add cluster              "), "create new cluster")
	fmt.Println(Yellow("new cluster              "), "alias for previous command")
	fmt.Println()
	fmt.Println(Blue("Configuration profiles:    "))
	fmt.Println(Yellow("list profiles            "), "list all profiles known to the service")
	fmt.Println(Yellow("describe profile ##      "), "describe profile selected by its ID")
	fmt.Println(Yellow("delete profile ##        "), "delete profile selected by its ID")
	fmt.Println()
	fmt.Println(Blue("Cluster configurations:    "))
	fmt.Println(Yellow("list configurations      "), "list all configurations known to the service")
	fmt.Println(Yellow("describe configuration ##"), "describe cluster configuration selected by its ID")
	fmt.Println(Yellow("add configuration        "), "add new configuration")
	fmt.Println(Yellow("new configuration        "), "alias for previous command")
	fmt.Println(Yellow("enable ##                "), "enable cluster configuration selected by its ID")
	fmt.Println(Yellow("disable ##               "), "disable cluster configuration selected by its ID")
	fmt.Println(Yellow("delete configuration ##  "), "delete configuration selected by its ID")
	fmt.Println()
	fmt.Println(Blue("Must-gather trigger:       "))
	fmt.Println(Yellow("list triggers            "), "list all triggers")
	fmt.Println(Yellow("describe trigger ##      "), "describe trigger selected by its ID")
	fmt.Println(Yellow("add trigger              "), "add new trigger")
	fmt.Println(Yellow("new trigger              "), "alias for previous command")
	fmt.Println(Yellow("activate trigger ##      "), "activate trigger selected by its ID")
	fmt.Println(Yellow("deactivate trigger ##    "), "deactivate trigger selected by its ID")
	fmt.Println(Yellow("delete trigger           "), "delete trigger")
	fmt.Println()
	fmt.Println(Blue("Other commands:"))
	fmt.Println(Yellow("version                  "), "print version information")
	fmt.Println(Yellow("quit                     "), "quit the application")
	fmt.Println(Yellow("exit                     "), "dtto")
	fmt.Println(Yellow("bye                      "), "dtto")
	fmt.Println(Yellow("help                     "), "this help")
	fmt.Println()
}
