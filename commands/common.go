/*
Copyright © 2019, 2020 Red Hat, Inc.

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
// https://godoc.org/github.com/RedHatInsights/insights-operator-cli/commands
//
// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-cli/packages/commands/common.html

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/logrusorgru/aurora"
	"os"
	"path/filepath"
)

// files will be filled by list of files that are found in given directory and
// displayed to user to select one of them
var files []prompt.Suggest

// colorizer contains instance of terminal colorizer interface
var colorizer aurora.Aurora

// SetColorizer function set the terminal colorizer.
func SetColorizer(c aurora.Aurora) {
	colorizer = c
}

// LoginCompleter implements a no-op completer needed to input random data, for
// example during testing.
func LoginCompleter(in prompt.Document) []prompt.Suggest {
	return nil
}

// ConfigFileCompleter function implements completer based on list of files.
func ConfigFileCompleter(in prompt.Document) []prompt.Suggest {
	return prompt.FilterHasPrefix(files, in.Text, true)
}

// ProceedQuestion ask user about y/n answer.
func ProceedQuestion(question string) bool {
	fmt.Println(colorizer.Red(question))
	// ask user and wait for input
	proceed := prompt.Input("proceed? [y/n] ", LoginCompleter)
	// only 'y' is accepted as "Yes" answer right now
	if proceed != "y" {
		fmt.Println(colorizer.Blue("cancelled"))
		return false
	}
	return true
}

// FillInConfigurationList function prepares a list of configuration files that
// are found in specified directory.
func FillInConfigurationList(directory string) error {
	files = []prompt.Suggest{}

	root := directory
	// iterate over all files and directories found
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			suggest := prompt.Suggest{
				Text: info.Name()}
			files = append(files, suggest)
		}
		return nil
	})

	// check for any error
	if err != nil {
		return err
	}
	return nil
}

// Quit function will exit from the CLI client.
func Quit() {
	fmt.Println(colorizer.Magenta("Quitting"))
	os.Exit(0)
}
