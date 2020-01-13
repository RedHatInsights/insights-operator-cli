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
