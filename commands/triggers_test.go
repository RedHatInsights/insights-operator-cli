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
	"regexp"
	"strings"
	"testing"
)

func tryToFindTrigger(t *testing.T, captured string, trigger string) {
	if !strings.Contains(captured, trigger) {
		t.Fatal("Can not find trigger:", trigger)
	}
}

// TestListOfTriggers checks whether the non-empty list of triggers read via REST API is displayed correctly
func TestListOfTriggers(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMock{}

	captured, err := capture.StandardOutput(func() {
		commands.ListOfTriggers(restAPIMock)
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "List of triggers for all clusters") {
		t.Fatal("Unexpected output:\n", captured)
	}

	// we expect six lines - title, column headers and four triggers
	numlines := strings.Count(captured, "\n")
	if numlines < 6 {
		t.Fatal("Not all triggers are listed in the output:\n", captured)
	}

	expectedTriggers := []string{
		"must-gather",
		"must-gather",
		"must-gather",
		"different-trigger",
	}
	for _, expectedTrigger := range expectedTriggers {
		tryToFindTrigger(t, captured, expectedTrigger)
	}
}

// TestListOfTriggers checks whether the empty list of triggers read via REST API is displayed correctly
func TestListOfTriggersEmptyList(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMockEmpty{}

	captured, err := capture.StandardOutput(func() {
		commands.ListOfTriggers(restAPIMock)
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "List of triggers for all clusters") {
		t.Fatal("Unexpected output:\n", captured)
	}

	// we expect two lines - title and column headers
	numlines := strings.Count(captured, "\n")
	if numlines > 2 {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestListOfTriggersErrorHandling checks whether error returned by REST API is handled correctly
func TestListOfTriggersErrorHandling(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMockErrors{}

	captured, err := capture.StandardOutput(func() {
		commands.ListOfTriggers(restAPIMock)
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Error reading list of triggers") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

func TestDescribeActivatedTrigger(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMock{}

	captured, err := capture.StandardOutput(func() {
		commands.DescribeTrigger(restAPIMock, "0")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Trigger info") {
		t.Fatal("Unexpected output:\n", captured)
	}
	if !strings.Contains(captured, "ffffffff-ffff-ffff-ffff-ffffffffffff") {
		t.Fatal("Can not find cluster ID:\n", captured)
	}
	if !strings.Contains(captured, "tester") {
		t.Fatal("Can not find name of user two triggered the trigger:\n", captured)
	}
	match, err := regexp.MatchString(`Active:.*yes`, captured)
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Fatal("Trigger is not activated as expected:\n", captured)
	}
}

func TestDescribeInactivatedTrigger(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMock{}

	captured, err := capture.StandardOutput(func() {
		commands.DescribeTrigger(restAPIMock, "1")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Trigger info") {
		t.Fatal("Unexpected output:\n", captured)
	}
	if !strings.Contains(captured, "ffffffff-ffff-ffff-ffff-ffffffffffff") {
		t.Fatal("Can not find cluster ID:\n", captured)
	}
	if !strings.Contains(captured, "tester") {
		t.Fatal("Can not find name of user two triggered the trigger:\n", captured)
	}
	match, err := regexp.MatchString(`Active:.*no`, captured)
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Fatal("Trigger is not deactivated as expected:\n", captured)
	}
}

func TestDescribeNonMustGatherTrigger(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMock{}

	captured, err := capture.StandardOutput(func() {
		commands.DescribeTrigger(restAPIMock, "2")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Trigger info") {
		t.Fatal("Unexpected output:\n", captured)
	}
	if !strings.Contains(captured, "00000000-0000-0000-0000-000000000000") {
		t.Fatal("Can not find cluster ID:\n", captured)
	}
	if !strings.Contains(captured, "tester") {
		t.Fatal("Can not find name of user two triggered the trigger:\n", captured)
	}
	match, err := regexp.MatchString(`Active:.*no`, captured)
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Fatal("Trigger is not deactivated as expected:\n", captured)
	}
}

func TestDescribeTriggerErrorHandling(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMockErrors{}

	captured, err := capture.StandardOutput(func() {
		commands.DescribeTrigger(restAPIMock, "1")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Error reading selected trigger") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

func TestAddTrigger(t *testing.T) {
}

func TestAddTriggerErrorHandling(t *testing.T) {
}

func TestDeleteTrigger(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMock{}

	captured, err := capture.StandardOutput(func() {
		commands.DeleteTrigger(restAPIMock, "0")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Trigger 0 has been deleted") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

func TestDeleteTriggerErrorHandling(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMockErrors{}

	captured, err := capture.StandardOutput(func() {
		commands.DeleteTrigger(restAPIMock, "0")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Error communicating with the service") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

func TestActivateTrigger(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMock{}

	captured, err := capture.StandardOutput(func() {
		commands.ActivateTrigger(restAPIMock, "0")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Trigger 0 has been activated") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

func TestActivateTriggerErrorHandling(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMockErrors{}

	captured, err := capture.StandardOutput(func() {
		commands.ActivateTrigger(restAPIMock, "0")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Error communicating with the service") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

func TestDeactivateTrigger(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMock{}

	captured, err := capture.StandardOutput(func() {
		commands.DeactivateTrigger(restAPIMock, "0")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Trigger 0 has been deactivated") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

func TestDeactivateTriggerErrorHandling(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMockErrors{}

	captured, err := capture.StandardOutput(func() {
		commands.DeactivateTrigger(restAPIMock, "0")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Error communicating with the service") {
		t.Fatal("Unexpected output:\n", captured)
	}
}
