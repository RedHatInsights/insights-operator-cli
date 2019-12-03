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

package main

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ThomasRooney/gexpect"
)

// startCLI starts CLI application w/o color output and w/o command-line completer.
func startCLI(t *testing.T) *gexpect.ExpectSubprocess {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	for !strings.HasSuffix(dir, "ccx/insights-operator-cli") { // make sure it's executed from the correct path
		err := os.Chdir("../")
		if err != nil {
			panic(err)
		}
		newDir, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		if strings.HasSuffix(newDir, "ccx/insights-operator-cli") {
			break
		}
	}

	child, err := gexpect.Spawn("./insights-operator-cli --colors=false --completer=false")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("CLI client has been started")
	return child
}

// pause just pauses test for a given amount of time to wait for tested application.
func pause() {
	time.Sleep(time.Second)
}

// expectOutput expects the specified output from the tested CLI client.
func expectOutput(t *testing.T, child *gexpect.ExpectSubprocess, output string) {
	err := child.ExpectTimeout(output, 2*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Expected output '" + output + "' has been found")
}

// expectPrompt expects the prompt from the tested CLI client.
func expectPrompt(t *testing.T, child *gexpect.ExpectSubprocess) {
	expectOutput(t, child, "> ")
}

// sendCommand sends command to the tested CLI client.
func sendCommand(t *testing.T, child *gexpect.ExpectSubprocess, command string) {
	err := child.SendLine(command)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Command '" + command + "' has been sent to CLI client")
}

// TestCheckQuitCommand check whether the client can be started and stopped using the 'quit' command.
func TestCheckQuitCommand(t *testing.T) {
	child := startCLI(t)
	expectPrompt(t, child)
	sendCommand(t, child, "quit")
	child.Wait()
}

// TestCheckVersionCommand check the 'version' command.
func TestCheckVersionCommand(t *testing.T) {
	child := startCLI(t)
	expectPrompt(t, child)
	sendCommand(t, child, "version")
	expectOutput(t, child, "Insights operator CLI client")
	expectOutput(t, child, "version")
	expectOutput(t, child, "compiled")
	expectPrompt(t, child)
	sendCommand(t, child, "quit")
	child.Wait()
}
