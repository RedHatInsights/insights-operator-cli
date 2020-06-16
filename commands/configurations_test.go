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
	"os"
	"strings"
	"testing"
)

func tryToFindConfiguration(t *testing.T, captured string, configuration string) {
	if !strings.Contains(captured, configuration) {
		t.Fatal("Can not find configuration:", configuration)
	}
}

// changeDirectory tries to change current directory with additional test
// whether the operation has been correct or not
func changeDirectory(t *testing.T, path string) {
	err := os.Chdir(path)
	if err != nil {
		t.Fatal(err)
	}
}

// TestListOfConfigurations checks whether the non-empty list of configurations
// read via REST API is displayed correctly
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

	// Mocked REST API returns three configurations, so we expect at least
	// one caption + 3 other lines in the output
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

// TestListOfConfigurationsEmptyList checks whether the empty list of
// configurations read via REST API is displayed correctly
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

	// Mocked REST API returns empty list, so just one caption + one
	// message is expected
	if numlines > 2 {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestListOfConfigurationsErrorHandling checks whether error returned by REST
// API is handled correctly
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

// TestDeleteClusterConfigurationError checks the command 'delete
// configuration' when error is reported by REST API
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

// TestEnableClusterConfigurationError checks the command 'enable
// configuration' when error is reported by REST API
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

// TestDisableClusterConfigurationError checks the command 'disable
// configuration' when error is reported by REST API
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

// TestDescribeConfigurationError checks the command 'describe configuration'
// when error is reported by REST API
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

// TestAddClusterConfigurationImpl checks the command 'add configuration'
func TestAddClusterConfigurationImpl(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMock{}

	captured, err := capture.StandardOutput(func() {
		changeDirectory(t, "../")
		commands.AddClusterConfigurationImpl(restAPIMock, "tester", "cluster0", "reason", "description", "configuration1.json")
		changeDirectory(t, "./commands")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Configuration has been created") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestAddClusterConfigurationImplError checks the command 'add configuration'
// when REST API fails with error
func TestAddClusterConfigurationImplError(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMockErrors{}

	captured, err := capture.StandardOutput(func() {
		changeDirectory(t, "../")
		commands.AddClusterConfigurationImpl(restAPIMock, "tester", "cluster0", "reason", "description", "configuration1.json")
		changeDirectory(t, "./commands")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Error communicating with the service") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestAddClusterConfigurationImplBadConfiguration checks the command 'add
// configuration' for non-existing configuration file
func TestAddClusterConfigurationImplBadConfiguration(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMock{}

	captured, err := capture.StandardOutput(func() {
		changeDirectory(t, "../")
		commands.AddClusterConfigurationImpl(restAPIMock, "tester", "cluster0", "reason", "description", "strange_configuration1.json")
		changeDirectory(t, "./commands")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.Contains(captured, "Cannot read configuration file") {
		t.Fatal("Unexpected output:\n", captured)
	}
	if !strings.Contains(captured, "no such file or directory") {
		t.Fatal("Unexpected output:\n", captured)
	}
	if strings.Contains(captured, "Configuration has been created") {
		t.Fatal("Configuration should not be created when configuration file does not exist")
	}
}
