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

func tryToFindProfile(t *testing.T, captured string, profileDescription string) {
	if !strings.Contains(captured, profileDescription) {
		t.Fatal("Can not find profile with description:", profileDescription)
	}
}

// TestListOfProfiles checks whether the non-empty list of configuration profiles read via REST API is displayed correctly
func TestListOfProfiles(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMock{}

	captured, err := capture.StandardOutput(func() {
		commands.ListOfProfiles(restAPIMock)
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "List of configuration profiles") {
		t.Fatal("Unexpected output:\n", captured)
	}

	// we expect four lines - title, column headers and two configuration profiles
	numlines := strings.Count(captured, "\n")
	if numlines < 4 {
		t.Fatal("Configuration profiles are not listed in the output:\n", captured)
	}

	expectedProfiles := []string{
		"default configuration profile",
		"another configuration profile",
	}
	for _, expectedProfile := range expectedProfiles {
		tryToFindProfile(t, captured, expectedProfile)
	}
}

// TestListOfProfilesNoProfiles checks whether the empty list of configuration profiles read via REST API is displayed correctly
func TestListOfProfilesNoProfiles(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMockEmpty{}

	captured, err := capture.StandardOutput(func() {
		commands.ListOfProfiles(restAPIMock)
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "List of configuration profiles") {
		t.Fatal("Unexpected output:\n", captured)
	}

	// we expect two lines - title and column headers
	numlines := strings.Count(captured, "\n")
	if numlines > 2 {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestListOfProfilesNoProfiles checks whether error returned by REST API is handled correctly
func TestListOfProfilesErrorHandling(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMockErrors{}

	captured, err := capture.StandardOutput(func() {
		commands.ListOfProfiles(restAPIMock)
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Error reading list of configuration profiles") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestDescribeProfile checks how the configuration profile is displayed on CLI
func TestDescribeProfile(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMock{}

	captured, err := capture.StandardOutput(func() {
		commands.DescribeProfile(restAPIMock, "0")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Configuration profile") {
		t.Fatal("Unexpected output:\n", captured)
	}
	if !strings.Contains(captured, "*configuration*") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestDescribeProfile checks error handling of REST API
func TestDescribeProfileErrorHandling(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMockErrors{}

	captured, err := capture.StandardOutput(func() {
		commands.DescribeProfile(restAPIMock, "0")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Error reading configuration profile") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestAddConfigurationProfileImpl checks the command 'add profile' when no error is reported by REST API
func TestAddConfigurationProfileImpl(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMock{}

	captured, err := capture.StandardOutput(func() {
		os.Chdir("../")
		commands.AddConfigurationProfileImpl(restAPIMock, "tester", "description", "configuration1.json")
		os.Chdir("./commands")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Configuration profile has been created") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestAddConfigurationProfileImplWrongConfiguration checks the command 'add profile' when configuration file does not exist
func TestAddConfigurationProfileImplWrongConfiguration(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMock{}

	captured, err := capture.StandardOutput(func() {
		os.Chdir("../")
		commands.AddConfigurationProfileImpl(restAPIMock, "tester", "description", "non-existing-configuration.json")
		os.Chdir("./commands")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Cannot read configuration file") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestAddConfigurationProfileImplErrorHandling checks the command 'add profile' when error is reported by REST API
func TestAddConfigurationProfileImplErrorHandling(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMockErrors{}

	captured, err := capture.StandardOutput(func() {
		os.Chdir("../")
		commands.AddConfigurationProfileImpl(restAPIMock, "tester", "description", "configuration1.json")
		os.Chdir("./commands")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Error communicating with the service") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestDeleteConfigurationProfile checks the command 'delete profile' when no error is reported by REST API
func TestDeleteConfigurationProfile(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMock{}

	captured, err := capture.StandardOutput(func() {
		commands.DeleteConfigurationProfile(restAPIMock, "0", false)
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Configuration profile 0 has been deleted") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestDeleteConfigurationProfile checks the command 'delete profile' when error is reported by REST API
func TestDeleteConfigurationProfileErrorHandling(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMockErrors{}

	captured, err := capture.StandardOutput(func() {
		commands.DeleteConfigurationProfile(restAPIMock, "0", false)
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Error communicating with the service") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestDeleteConfigurationProfileNoConfirm checks the command 'delete profile' when no error is reported by REST API
func TestDeleteConfigurationProfileNoConfirm(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMock{}

	captured, err := capture.StandardOutput(func() {
		commands.DeleteConfigurationProfileNoConfirm(restAPIMock, "0")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Configuration profile 0 has been deleted") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestDeleteConfigurationProfileNoConfirm checks the command 'delete profile' when error is reported by REST API
func TestDeleteConfigurationProfileNoConfirmErrorHandling(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMockErrors{}

	captured, err := capture.StandardOutput(func() {
		commands.DeleteConfigurationProfileNoConfirm(restAPIMock, "0")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Error communicating with the service") {
		t.Fatal("Unexpected output:\n", captured)
	}
}
