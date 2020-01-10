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

// TestDeleteClusterConfiguration checks the command 'delete configuration'
func TestDeleteClusterConfiguration(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMock{}

	captured, err := capture.StandardOutput(func() {
		commands.DeleteClusterConfiguration(restAPIMock, "1")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Configuration 1 has been deleted") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestDeleteClusterConfigurationError checks the command 'delete configuration' when error is reported by REST API
func TestDeleteClusterConfigurationError(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMockErrors{}

	captured, err := capture.StandardOutput(func() {
		commands.DeleteClusterConfiguration(restAPIMock, "0")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Error communicating with the service") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestEnableClusterConfiguration checks the command 'enable configuration'
func TestEnableClusterConfiguration(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMock{}

	captured, err := capture.StandardOutput(func() {
		commands.EnableClusterConfiguration(restAPIMock, "1")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Configuration 1 has been enabled") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestEnableClusterConfigurationError checks the command 'enable configuration' when error is reported by REST API
func TestEnableClusterConfigurationError(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMockErrors{}

	captured, err := capture.StandardOutput(func() {
		commands.EnableClusterConfiguration(restAPIMock, "0")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Error communicating with the service") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestDisableClusterConfiguration checks the command 'disable configuration'
func TestDisableClusterConfiguration(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMock{}

	captured, err := capture.StandardOutput(func() {
		commands.DisableClusterConfiguration(restAPIMock, "1")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Configuration 1 has been disabled") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestDisableClusterConfigurationError checks the command 'disable configuration' when error is reported by REST API
func TestDisableClusterConfigurationError(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMockErrors{}

	captured, err := capture.StandardOutput(func() {
		commands.DisableClusterConfiguration(restAPIMock, "0")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Error communicating with the service") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestDescribeConfiguration checks the command 'describe configuration'
func TestDescribeConfiguration(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMock{}

	captured, err := capture.StandardOutput(func() {
		commands.DescribeConfiguration(restAPIMock, "0")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Configuration for cluster 0") {
		t.Fatal("Unexpected output:\n", captured)
	}
	if !strings.Contains(captured, "configuration#0") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestDescribeConfigurationError checks the command 'describe configuration' when error is reported by REST API
func TestDescribeConfigurationError(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMockErrors{}

	captured, err := capture.StandardOutput(func() {
		commands.DescribeConfiguration(restAPIMock, "0")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Error reading cluster configuration") {
		t.Fatal("Unexpected output:\n", captured)
	}
}
