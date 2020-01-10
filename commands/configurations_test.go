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

package commands_test

import (
	"github.com/redhatinsighs/insights-operator-cli/commands"
	"github.com/tisnik/go-capture"
	"strings"
	"testing"
)

func tryToFindConfiguration(t *testing.T, captured string, configuration string) {
	if !strings.Contains(captured, configuration) {
		t.Fatal("Can not find configuration:", configuration)
	}
}

// TestListOfConfigurations checks whether the non-empty list of configurations read via REST API is displayed correctly
func TestListOfConfigurations(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMock{}

	captured, err := capture.StandardOutput(func() {
		commands.ListOfConfigurations(restAPIMock, "")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "List of configurations for all clusters") {
		t.Fatal("Unexpected output:\n", captured)
	}

	numlines := strings.Count(captured, "\n")
	if numlines <= 4 {
		t.Fatal("Configurations are not listed in the output:\n", captured)
	}
	expectedConfigurations := []string{
		"configuration1",
		"configuration2",
		"configuration3",
	}
	for _, expectedConfiguration := range expectedConfigurations {
		tryToFindConfiguration(t, captured, expectedConfiguration)
	}
}

// TestListOfConfigurationsEmptyList checks whether the empty list of configurations read via REST API is displayed correctly
func TestListOfConfigurationsEmptyList(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMockEmpty{}

	captured, err := capture.StandardOutput(func() {
		commands.ListOfConfigurations(restAPIMock, "")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "List of configurations for all clusters") {
		t.Fatal("Unexpected output:\n", captured)
	}

	numlines := strings.Count(captured, "\n")
	if numlines > 2 {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestListOfConfigurationsErrorHandling checks whether error returned by REST API is handled correctly
func TestListOfConfigurationsErrorHandling(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMockErrors{}

	captured, err := capture.StandardOutput(func() {
		commands.ListOfConfigurations(restAPIMock, "")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Error reading list of configurations") {
		t.Fatal("Unexpected output:\n", captured)
	}
}
