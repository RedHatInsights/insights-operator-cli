/*
Copyright © 2019, 2020, 2021, 2022 Red Hat, Inc.

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

// Unit tests checking functions for manipulation with profiles via REST API.

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-cli/packages/commands/profiles_test.html

import (
	"github.com/RedHatInsights/insights-operator-cli/commands"
	"github.com/tisnik/go-capture"
	"strings"
	"testing"
)

// tryToFindProfile is a helper function that tries to find a given profile in
// captured output. If profile info is not found, the test that calls this
// function, fails.
func tryToFindProfile(t *testing.T, captured, profileDescription string) {
	if !strings.Contains(captured, profileDescription) {
		// if profile info is not found, the test should fail
		t.Fatal("Can not find profile with description:", profileDescription)
	}
}

// TestListOfProfiles function checks whether the non-empty list of
// configuration profiles read via REST API is displayed correctly.
func TestListOfProfiles(t *testing.T) {
	// turn off any colorization on standard output
	// so we'll be able to capture pure messages w/o terminal control codes
	configureColorizer()

	// use mocked REST API instead of the real one
	// to perform profile-related command or query
	restAPIMock := RestAPIMock{}

	// use go-capture package to capture all writes to standard output
	captured, err := capture.StandardOutput(func() {
		commands.ListOfProfiles(restAPIMock)
	})

	// check if capture was done correctly
	checkCapturedOutput(t, captured, err)

	// test the captured output
	if !strings.HasPrefix(captured, "List of configuration profiles") {
		t.Fatal("Unexpected output:\n", captured)
	}

	// we expect four lines - title, column headers and two configuration profiles
	numlines := strings.Count(captured, "\n")
	if numlines < 4 {
		t.Fatal("Configuration profiles are not listed in the output:\n", captured)
	}

	// list of expected profiles to be displayed on standard output
	expectedProfiles := []string{
		"default configuration profile",
		"another configuration profile",
	}

	// check the actual output displayed on terminal
	for _, expectedProfile := range expectedProfiles {
		tryToFindProfile(t, captured, expectedProfile)
	}
}

// TestListOfProfilesNoProfiles function checks whether the empty list of
// configuration profiles read via REST API is displayed correctly.
func TestListOfProfilesNoProfiles(t *testing.T) {
	// turn off any colorization on standard output
	// so we'll be able to capture pure messages w/o terminal control codes
	configureColorizer()

	// use mocked REST API instead of the real one
	// to perform profile-related command or query
	restAPIMock := RestAPIMockEmpty{}

	// use go-capture package to capture all writes to standard output
	captured, err := capture.StandardOutput(func() {
		commands.ListOfProfiles(restAPIMock)
	})

	// check if capture was done correctly
	checkCapturedOutput(t, captured, err)

	// test the captured output
	if !strings.HasPrefix(captured, "List of configuration profiles") {
		t.Fatal("Unexpected output:\n", captured)
	}

	// we expect two lines - title and column headers
	numlines := strings.Count(captured, "\n")

	// check the actual output displayed on terminal
	if numlines > 2 {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestListOfProfilesNoProfiles function checks whether error returned by REST
// API is handled correctly.
func TestListOfProfilesErrorHandling(t *testing.T) {
	// turn off any colorization on standard output
	// so we'll be able to capture pure messages w/o terminal control codes
	configureColorizer()

	// use mocked REST API instead of the real one
	// to perform profile-related command or query
	restAPIMock := RestAPIMockErrors{}

	// use go-capture package to capture all writes to standard output
	captured, err := capture.StandardOutput(func() {
		commands.ListOfProfiles(restAPIMock)
	})

	// check if capture was done correctly
	checkCapturedOutput(t, captured, err)

	// test the captured output
	if !strings.HasPrefix(captured, ErrorReadingListOfConfigurationProfiles) {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestDescribeProfile function checks how the configuration profile is
// displayed on CLI.
func TestDescribeProfile(t *testing.T) {
	// turn off any colorization on standard output
	// so we'll be able to capture pure messages w/o terminal control codes
	configureColorizer()

	// use mocked REST API instead of the real one
	// to perform profile-related command or query
	restAPIMock := RestAPIMock{}

	// use go-capture package to capture all writes to standard output
	captured, err := capture.StandardOutput(func() {
		commands.DescribeProfile(restAPIMock, "0")
	})

	// check if capture was done correctly
	checkCapturedOutput(t, captured, err)

	// test the captured output
	if !strings.HasPrefix(captured, "Configuration profile") {
		t.Fatal("Unexpected output:\n", captured)
	}

	// configuration profile needs to be displayed
	if !strings.Contains(captured, "*configuration*") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestDescribeProfile function checks error handling of REST API.
func TestDescribeProfileErrorHandling(t *testing.T) {
	// turn off any colorization on standard output
	// so we'll be able to capture pure messages w/o terminal control codes
	configureColorizer()

	// use mocked REST API instead of the real one
	// to perform profile-related command or query
	restAPIMock := RestAPIMockErrors{}

	// use go-capture package to capture all writes to standard output
	captured, err := capture.StandardOutput(func() {
		commands.DescribeProfile(restAPIMock, "0")
	})

	// check if capture was done correctly
	checkCapturedOutput(t, captured, err)

	// test the captured output
	if !strings.HasPrefix(captured, commands.ErrorReadingConfigurationProfile) {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestAddConfigurationProfileImpl function checks the command 'add profile'
// when no error is reported by REST API.
func TestAddConfigurationProfileImpl(t *testing.T) {
	// turn off any colorization on standard output
	// so we'll be able to capture pure messages w/o terminal control codes
	configureColorizer()

	// use mocked REST API instead of the real one
	// to perform profile-related command or query
	restAPIMock := RestAPIMock{}

	// use go-capture package to capture all writes to standard output
	captured, err := capture.StandardOutput(func() {
		changeDirectory(t, "../")
		commands.AddConfigurationProfileImpl(restAPIMock, "tester", "description", "configuration1.json")
		changeDirectory(t, "./commands")
	})

	// check if capture was done correctly
	checkCapturedOutput(t, captured, err)

	// test the captured output
	if !strings.HasPrefix(captured, "Configuration profile has been created") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestAddConfigurationProfileImplWrongConfiguration function checks the
// command 'add profile' when configuration file does not exist.
func TestAddConfigurationProfileImplWrongConfiguration(t *testing.T) {
	// turn off any colorization on standard output
	// so we'll be able to capture pure messages w/o terminal control codes
	configureColorizer()

	// use mocked REST API instead of the real one
	// to perform profile-related command or query
	restAPIMock := RestAPIMock{}

	// use go-capture package to capture all writes to standard output
	captured, err := capture.StandardOutput(func() {
		changeDirectory(t, "../")
		commands.AddConfigurationProfileImpl(restAPIMock, "tester", "description", "non-existing-configuration.json")
		changeDirectory(t, "./commands")
	})

	// check if capture was done correctly
	checkCapturedOutput(t, captured, err)

	// test the captured output
	if !strings.HasPrefix(captured, "Cannot read configuration file") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestAddConfigurationProfileImplErrorHandling function checks the command
// 'add profile' when error is reported by REST API.
func TestAddConfigurationProfileImplErrorHandling(t *testing.T) {
	// turn off any colorization on standard output
	// so we'll be able to capture pure messages w/o terminal control codes
	configureColorizer()

	// use mocked REST API instead of the real one
	// to perform profile-related command or query
	restAPIMock := RestAPIMockErrors{}

	// use go-capture package to capture all writes to standard output
	captured, err := capture.StandardOutput(func() {
		changeDirectory(t, "../")
		commands.AddConfigurationProfileImpl(restAPIMock, "tester", "description", "configuration1.json")
		changeDirectory(t, "./commands")
	})

	// check if capture was done correctly
	checkCapturedOutput(t, captured, err)

	// test the captured output
	if !strings.HasPrefix(captured, "Error communicating with the service") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestDeleteConfigurationProfile function checks the command 'delete profile'
// when no error is reported by REST API.
func TestDeleteConfigurationProfile(t *testing.T) {
	// turn off any colorization on standard output
	// so we'll be able to capture pure messages w/o terminal control codes
	configureColorizer()

	// use mocked REST API instead of the real one
	// to perform profile-related command or query
	restAPIMock := RestAPIMock{}

	// use go-capture package to capture all writes to standard output
	captured, err := capture.StandardOutput(func() {
		commands.DeleteConfigurationProfile(restAPIMock, "0", false)
	})

	// check if capture was done correctly
	checkCapturedOutput(t, captured, err)

	// test the captured output
	if !strings.HasPrefix(captured, "Configuration profile 0 has been deleted") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestDeleteConfigurationProfile function checks the command 'delete profile'
// when error is reported by REST API.
func TestDeleteConfigurationProfileErrorHandling(t *testing.T) {
	// turn off any colorization on standard output
	// so we'll be able to capture pure messages w/o terminal control codes
	configureColorizer()

	// use mocked REST API instead of the real one
	// to perform profile-related command or query
	restAPIMock := RestAPIMockErrors{}

	// use go-capture package to capture all writes to standard output
	captured, err := capture.StandardOutput(func() {
		commands.DeleteConfigurationProfile(restAPIMock, "0", false)
	})

	// check if capture was done correctly
	checkCapturedOutput(t, captured, err)

	// test the captured output
	if !strings.HasPrefix(captured, "Error communicating with the service") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestDeleteConfigurationProfileNoConfirm function checks the command 'delete
// profile' when no error is reported by REST API.
func TestDeleteConfigurationProfileNoConfirm(t *testing.T) {
	// turn off any colorization on standard output
	// so we'll be able to capture pure messages w/o terminal control codes
	configureColorizer()

	// use mocked REST API instead of the real one
	// to perform profile-related command or query
	restAPIMock := RestAPIMock{}

	// use go-capture package to capture all writes to standard output
	captured, err := capture.StandardOutput(func() {
		commands.DeleteConfigurationProfileNoConfirm(restAPIMock, "0")
	})

	// check if capture was done correctly
	checkCapturedOutput(t, captured, err)

	// test the captured output
	if !strings.HasPrefix(captured, "Configuration profile 0 has been deleted") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestDeleteConfigurationProfileNoConfirm function checks the command 'delete
// profile' when error is reported by REST API.
func TestDeleteConfigurationProfileNoConfirmErrorHandling(t *testing.T) {
	// turn off any colorization on standard output
	// so we'll be able to capture pure messages w/o terminal control codes
	configureColorizer()

	// use mocked REST API instead of the real one
	// to perform profile-related command or query
	restAPIMock := RestAPIMockErrors{}

	// use go-capture package to capture all writes to standard output
	captured, err := capture.StandardOutput(func() {
		commands.DeleteConfigurationProfileNoConfirm(restAPIMock, "0")
	})

	// check if capture was done correctly
	checkCapturedOutput(t, captured, err)

	// test the captured output
	if !strings.HasPrefix(captured, "Error communicating with the service") {
		t.Fatal("Unexpected output:\n", captured)
	}
}
