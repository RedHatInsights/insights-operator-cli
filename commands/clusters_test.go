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
// https://redhatinsights.github.io/insights-operator-cli/packages/commands/clusters_test.html

import (
	"github.com/RedHatInsights/insights-operator-cli/commands"
	"github.com/tisnik/go-capture"
	"strings"
	"testing"
)

// checkCapturedOutput function checks if the capturing of standard output was
// correct.
func checkCapturedOutput(t *testing.T, captured string, err error) {
	if err != nil {
		t.Fatal("Unable to capture standard output", err)
	}
	if captured == "" {
		t.Fatal("Standard output is empty")
	}
}

// tryToFindCluster helper function checks if captured standard output contains
// provided cluster name or not.
func tryToFindCluster(t *testing.T, captured, clusterName string) {
	if !strings.Contains(captured, clusterName) {
		t.Fatal("Can not find cluster:", clusterName)
	}
}

// TestListOfClusters function checks whether the non-empty list of clusters
// read via REST API is displayed correctly
func TestListOfClusters(t *testing.T) {
	// turn off any colorization on standard output
	configureColorizer()

	// use mocked REST API instead of the real one
	restAPIMock := RestAPIMock{}

	// use go-capture package to capture all writes to standard output
	captured, err := capture.StandardOutput(func() {
		commands.ListOfClusters(restAPIMock)
	})

	// check if capture operation was finished correctly
	checkCapturedOutput(t, captured, err)

	// check if header was written into standard output
	if !strings.HasPrefix(captured, "List of clusters") {
		t.Fatal("Unexpected output:\n", captured)
	}

	// check if list of clusters was written onto standard output
	numlines := strings.Count(captured, "\n")
	if numlines <= 4 {
		t.Fatal("Clusters are not listed in the output:\n", captured)
	}
	expectedClusters := []string{
		"c8590f31-e97e-4b85-b506-c45ce1911a12",
		"ffffffff-ffff-ffff-ffff-ffffffffffff",
		"00000000-0000-0000-0000-000000000000",
	}
	for _, expectedCluster := range expectedClusters {
		tryToFindCluster(t, captured, expectedCluster)
	}
}

// TestListOfClustersNoClusters function checks whether the empty list of
// clusters read via REST API is displayed correctly
func TestListOfClustersNoClusters(t *testing.T) {
	// turn off any colorization on standard output
	configureColorizer()

	// use mocked REST API instead of the real one
	restAPIMock := RestAPIMockEmpty{}

	// use go-capture package to capture all writes to standard output
	captured, err := capture.StandardOutput(func() {
		commands.ListOfClusters(restAPIMock)
	})

	// check if capture operation was finished correctly
	checkCapturedOutput(t, captured, err)

	// check if header was written into standard output
	if !strings.HasPrefix(captured, "List of clusters") {
		t.Fatal("Unexpected output:\n", captured)
	}

	// check if list of clusters was written onto standard output
	numlines := strings.Count(captured, "\n")
	if numlines > 2 {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestListOfClusters function checks whether error returned by REST API is
// handled correctly
func TestListOfClustersErrorHandling(t *testing.T) {
	// turn off any colorization on standard output
	configureColorizer()

	// use mocked REST API instead of the real one
	restAPIMock := RestAPIMockErrors{}

	// use go-capture package to capture all writes to standard output
	captured, err := capture.StandardOutput(func() {
		commands.ListOfClusters(restAPIMock)
	})

	// check if capture operation was finished correctly
	checkCapturedOutput(t, captured, err)

	// check if error message was written into standard output
	if !strings.HasPrefix(captured, commands.ErrorReadingListOfClusters) {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestAddCluster function checks the command 'add cluster'
func TestAddCluster(t *testing.T) {
	// turn off any colorization on standard output
	configureColorizer()

	// use mocked REST API instead of the real one
	restAPIMock := RestAPIMock{}

	// use go-capture package to capture all writes to standard output
	captured, err := capture.StandardOutput(func() {
		commands.AddCluster(restAPIMock, "c8590f31-e97e-4b85-b506-c45ce1911a12")
	})

	// check if capture operation was finished correctly
	checkCapturedOutput(t, captured, err)

	// check if cluster has been added via mocked REST API
	if !strings.HasPrefix(captured, "Cluster c8590f31-e97e-4b85-b506-c45ce1911a12 has been added") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestAddClusterError checks the command 'add cluster' when error is reported
// by REST API
func TestAddClusterError(t *testing.T) {
	// turn off any colorization on standard output
	configureColorizer()

	// use mocked REST API instead of the real one
	restAPIMock := RestAPIMockErrors{}

	// use go-capture package to capture all writes to standard output
	captured, err := capture.StandardOutput(func() {
		commands.AddCluster(restAPIMock, "clusterName")
	})

	// check if capture operation was finished correctly
	checkCapturedOutput(t, captured, err)

	// check if error message was written into standard output as expected
	if !strings.HasPrefix(captured, "Error communicating with the service") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestDeleteCluster checks the command 'delete cluster'
func TestDeleteCluster(t *testing.T) {
	// turn off any colorization on standard output
	configureColorizer()
	restAPIMock := RestAPIMock{}

	// use go-capture package to capture all writes to standard output
	captured, err := capture.StandardOutput(func() {
		commands.DeleteCluster(restAPIMock, "c8590f31-e97e-4b85-b506-c45ce1911a12", false)
	})

	// check if capture operation was finished correctly
	checkCapturedOutput(t, captured, err)

	// check if cluster has been deleted via mocked REST API
	if !strings.HasPrefix(captured, "Cluster c8590f31-e97e-4b85-b506-c45ce1911a12 has been deleted") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestDeleteClusterError checks the command 'delete cluster' when error is
// reported by REST API
func TestDeleteClusterError(t *testing.T) {
	// turn off any colorization on standard output
	configureColorizer()
	restAPIMock := RestAPIMockErrors{}

	// use go-capture package to capture all writes to standard output
	captured, err := capture.StandardOutput(func() {
		commands.DeleteCluster(restAPIMock, "c8590f31-e97e-4b85-b506-c45ce1911a12", false)
	})

	// check if capture operation was finished correctly
	checkCapturedOutput(t, captured, err)

	// check if error message was written into standard output as expected
	if !strings.HasPrefix(captured, "Error communicating with the service") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestDeleteClusterNoConfirm checks the command 'delete cluster' w/o
// confirmation of the command
func TestDeleteClusterNoConfirm(t *testing.T) {
	// turn off any colorization on standard output
	configureColorizer()
	restAPIMock := RestAPIMock{}

	// use go-capture package to capture all writes to standard output
	captured, err := capture.StandardOutput(func() {
		commands.DeleteClusterNoConfirm(restAPIMock, "c8590f31-e97e-4b85-b506-c45ce1911a12")
	})

	// check if capture operation was finished correctly
	checkCapturedOutput(t, captured, err)

	// check if cluster has been deleted via mocked REST API
	if !strings.HasPrefix(captured, "Cluster c8590f31-e97e-4b85-b506-c45ce1911a12 has been deleted") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestDeleteClusterNoConfirmError checks the command 'delete cluster' when
// error is reported by REST API
func TestDeleteClusterNoConfirmError(t *testing.T) {
	// turn off any colorization on standard output
	configureColorizer()
	restAPIMock := RestAPIMockErrors{}

	// use go-capture package to capture all writes to standard output
	captured, err := capture.StandardOutput(func() {
		commands.DeleteClusterNoConfirm(restAPIMock, "c8590f31-e97e-4b85-b506-c45ce1911a12")
	})

	// check if capture operation was finished correctly
	checkCapturedOutput(t, captured, err)

	// check if error message was written into standard output as expected
	if !strings.HasPrefix(captured, "Error communicating with the service") {
		t.Fatal("Unexpected output:\n", captured)
	}
}
