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

func checkCapturedOutput(t *testing.T, captured string, err error) {
	if err != nil {
		t.Fatal("Unable to capture standard output", err)
	}
	if captured == "" {
		t.Fatal("Standard output is empty")
	}
}

func tryToFindCluster(t *testing.T, captured string, clusterName string) {
	if !strings.Contains(captured, clusterName) {
		t.Fatal("Can not find cluster:", clusterName)
	}
}

// TestListOfClusters checks whether the non-empty list of clusters read via REST API is displayed correctly
func TestListOfClusters(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMock{}

	captured, err := capture.StandardOutput(func() {
		commands.ListOfClusters(restAPIMock)
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "List of clusters") {
		t.Fatal("Unexpected output:\n", captured)
	}

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

// TestListOfClusters checks whether the empty list of clusters read via REST API is displayed correctly
func TestListOfClustersNoClusters(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMockEmpty{}

	captured, err := capture.StandardOutput(func() {
		commands.ListOfClusters(restAPIMock)
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "List of clusters") {
		t.Fatal("Unexpected output:\n", captured)
	}

	numlines := strings.Count(captured, "\n")
	if numlines > 2 {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestListOfClusters checks whether error returned by REST API is handled correctly
func TestListOfClustersErrorHandling(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMockErrors{}

	captured, err := capture.StandardOutput(func() {
		commands.ListOfClusters(restAPIMock)
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Error reading list of clusters") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestAddCluster checks the command 'add cluster'
func TestAddCluster(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMock{}

	captured, err := capture.StandardOutput(func() {
		commands.AddCluster(restAPIMock, "clusterName")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Cluster has been added") {
		t.Fatal("Unexpected output:\n", captured)
	}
}

// TestAddClusterError checks the command 'add cluster' when error is reported by REST API
func TestAddClusterError(t *testing.T) {
	configureColorizer()
	restAPIMock := RestAPIMockErrors{}

	captured, err := capture.StandardOutput(func() {
		commands.AddCluster(restAPIMock, "clusterName")
	})

	checkCapturedOutput(t, captured, err)
	if !strings.HasPrefix(captured, "Error communicating with the service") {
		t.Fatal("Unexpected output:\n", captured)
	}
}
