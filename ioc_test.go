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
package main_test

// Documentation in literate-programming-style is available at:
// https://RedHatInsights.github.io/insights-operator-cli/packages/ioc_test.html

import (
	"github.com/c-bata/go-prompt"
	"github.com/logrusorgru/aurora"
	"testing"

	"github.com/RedHatInsights/insights-operator-cli"
)

// createDocumentWithCommand function constructs an instance of prompt.Document
// containing the command and cursor position.
func createDocumentWithCommand(t *testing.T, command string) prompt.Document {
	// try to allocate a buffer
	buffer := prompt.NewBuffer()
	if buffer == nil {
		t.Fatal("Error in prompt library - can not constructs new buffer")
	}

	// insert command into buffer
	buffer.InsertText(command, false, true)

	// and gather instance of new document
	document := buffer.Document()
	if document == nil {
		t.Fatal("Error in prompt library - can not get document for a buffer")
	}
	return *document
}

// checkSuggestionCount function checks the number of suggestions returned by
// suggester.
func checkSuggestionCount(t *testing.T, suggests []prompt.Suggest, expected int) {
	// test if number of suffestions is expected
	if len(suggests) != expected {
		t.Fatal("Invalid suggestion returned by completer:", suggests)
	}
}

// checkSuggestionCount function checks the suggestion text and description.
func checkSuggestion(t *testing.T, suggest prompt.Suggest, command, description string) {
	// test suggestion text by comparing it with command
	if suggest.Text != command {
		t.Fatal("Invalid suggestion command:", suggest.Text)
	}

	// test suggestion description
	if suggest.Description != description {
		t.Fatal("Invalid suggestion description:", suggest.Description)
	}
}

// TestCompleterEmptyInput function checks which suggestions are returned for
// empty input.
func TestCompleterEmptyInput(t *testing.T) {
	// test the suggestion(s) for empty input
	suggests := main.Completer(createDocumentWithCommand(t, ""))
	// no suggestions are expected
	checkSuggestionCount(t, suggests, 0)
}

// TestCompleterHelpCommand function checks which suggestions are returned for
// 'help' input.
func TestCompleterHelpCommand(t *testing.T) {
	// test the suggestion(s) for help command
	suggests := main.Completer(createDocumentWithCommand(t, "help"))

	// just one suggestion is expected
	checkSuggestionCount(t, suggests, 1)
	checkSuggestion(t, suggests[0], "help", "show help with all commands")
}

// TestReadConfiguration function tries to read configuration from existing
// configuration file.
func TestReadConfiguration(t *testing.T) {
	// test the suggestion(s) for command for reading configuration file
	_, err := main.ReadConfiguration("config")
	if err != nil {
		t.Fatal("Error during reading configuration", err)
	}
}

// TestReadConfigurationNegative function tries to read configuration from
// non-existing configuration file.
func TestReadConfigurationNegative(t *testing.T) {
	// test the suggestion(s) for command for reading configuration file
	_, err := main.ReadConfiguration("this_does_not_exists")
	if err == nil {
		t.Fatal("Error expected during reading configuration from non-existing file")
	}
}

// TestPrintVersion function is dummy ATM - we'll check versions etc. in
// integration tests.
func TestPrintVersion(t *testing.T) {
	// make sure the colorizer is initialized
	*main.Colorizer = aurora.NewAurora(true)

	// just print the version w/o any checks
	main.PrintVersion()
}
