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

package main

import (
	"testing"
)

// TestCheckQuitCommand check whether the client can be started and stopped using the 'quit', 'exit' or 'bye' command.
func TestCheckQuitCommand(t *testing.T) {
	// all commands that can be used to quit the CLI client
	commands := []string{
		"quit", "exit", "bye",
	}

	// try all quit-like commands
	for _, command := range commands {

		// start the CLI client
		child := startCLI(t)

		// send quit-like command to it
		sendCommand(t, child, command)

		// and check if the client has been stopped
		err := child.Wait()
		if err != nil {
			t.Fatal(err)
		}
		t.Log("CLI client has been stopped")
	}
}

// TestCheckVersionCommand check the output of 'version' command.
func TestCheckVersionCommand(t *testing.T) {
	// start the CLI client
	child := startCLI(t)

	// client needs to be shut down at the end of this test
	defer quitCLI(t, child)

	// send the 'version' command to the CLI client and check its basic output
	sendCommand(t, child, "version")
	expectOutput(t, child, "Insights operator CLI client")
	// TODO: more thorought check of 'version' value
	expectOutput(t, child, "version")
	// TODO: more thorought check of 'compiled' value
	expectOutput(t, child, "compiled")

	// at the end, the standard prompt has to be displayed
	expectPrompt(t, child)
}
