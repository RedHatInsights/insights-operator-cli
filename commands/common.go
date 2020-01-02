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

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/logrusorgru/aurora"
	"os"
	"path/filepath"
)

var files []prompt.Suggest

var colorizer aurora.Aurora

// SetColorizer set the terminal colorizer
func SetColorizer(c aurora.Aurora) {
	colorizer = c
}

// LoginCompleter implements a no-op completer needed to input random data
func LoginCompleter(in prompt.Document) []prompt.Suggest {
	return nil
}

// ConfigFileCompleter implements completer based on list of files
func ConfigFileCompleter(in prompt.Document) []prompt.Suggest {
	return prompt.FilterHasPrefix(files, in.Text, true)
}

// ProceedQuestion ask user about y/n answer.
func ProceedQuestion(question string) bool {
	fmt.Println(colorizer.Red(question))
	proceed := prompt.Input("proceed? [y/n] ", LoginCompleter)
	if proceed != "y" {
		fmt.Println(colorizer.Blue("cancelled"))
		return false
	}
	return true
}

// FillInConfigurationList prepares a list of configuration files
func FillInConfigurationList(directory string) error {
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

// Quit will exit from the CLI client
func Quit() {
	fmt.Println(colorizer.Magenta("Quitting"))
	os.Exit(0)
}
