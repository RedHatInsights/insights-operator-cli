/*
Copyright © 2019, 2020, 2021 Red Hat, Inc.

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

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-cli/packages/commands/license_test.html

import (
	"github.com/RedHatInsights/insights-operator-cli/commands"
	"github.com/tisnik/go-capture"
	"strings"
	"testing"
)

// TestCommandLicense checks if the command 'license' displays actual license
func TestCommandLicense(t *testing.T) {
	// turn off any colorization on standard output
	configureColorizer()

	// use go-capture package to capture all writes to standard output
	captured, err := capture.StandardOutput(func() {
		commands.PrintLicense()
	})

	// check if capture operation was finished correctly
	if err != nil {
		t.Fatal("Unable to capture standard output", err)
	}

	// check if capture operation captures at least some output
	if captured == "" {
		t.Fatal("Standard output is empty")
	}

	// check if header was written into standard output
	if !strings.HasPrefix(captured, "License\n") {
		t.Fatal("Unexpected output:\n", captured)
	}

	// check if message with license was written onto standard output
	numlines := strings.Count(captured, "\n")
	if numlines <= 1 {
		t.Fatal("License is trimmed:\n", captured)
	}
}
