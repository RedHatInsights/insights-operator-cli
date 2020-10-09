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

// Implementation of functional tests for the CLI client (application). These
// tests start CLI client, send commands to it and check the output (if it is
// the same as expected). In order to work properly, the CLI client needs to be
// configured to not to use TAB-completion and color output needs to be
// disabled as well. Additionally, the controller service needs to be started
// in background, because CLI client calls this service for almost all
// commands.
package main

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-cli/packages/tests/functional_test.html

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ThomasRooney/gexpect"
)

const (
	commandTimeout = 2 * time.Second
	prompt         = "> "
)

// startCLI function starts CLI application w/o color output and w/o command-line completer.
func startCLI(t *testing.T) *gexpect.ExpectSubprocess {
	// we need to know the current working directory so the CLI tool will be started from the right place
	// (tests might be stored in different sub-directory)
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	for !strings.HasSuffix(dir, "/insights-operator-cli") { // make sure it's executed from the correct path
		// try to traverse through super-directories
		err := os.Chdir("../")
		if err != nil {
			panic(err)
		}
		newDir, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		if strings.HasSuffix(newDir, "/insights-operator-cli") {
			break
		}
	}

	// start the CLI tool
	child, err := gexpect.Spawn("./insights-operator-cli --colors=false --completer=false")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("CLI client has been started")

	// now, after the start, we should see prompt
	expectPrompt(t, child)
	return child
}

// quitCLI function quits CLI tool that was started as children
func quitCLI(t *testing.T, child *gexpect.ExpectSubprocess) {
	// check if child process has been started before
	if child == nil {
		t.Fatal("Child process has not been started")
	}

	// the 'quit' command should quit the application
	sendCommand(t, child, "quit")

	// make it breath a bit
	err := child.Wait()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("CLI client has been stopped")
}

// expectOutput function expects the specified output from the tested CLI client.
func expectOutput(t *testing.T, child *gexpect.ExpectSubprocess, output string) {
	// check if the expected output has been displayed
	err := child.ExpectTimeout(output, commandTimeout)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Expected output '" + output + "' has been found")
}

// expectPrompt function expects the prompt from the tested CLI client.
func expectPrompt(t *testing.T, child *gexpect.ExpectSubprocess) {
	expectOutput(t, child, prompt)
}

// sendCommand sends command to the tested CLI client.
func sendCommand(t *testing.T, child *gexpect.ExpectSubprocess, command string) {
	// send a command to the tested CLI client
	err := child.SendLine(command)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Command '" + command + "' has been sent to CLI client")
}
