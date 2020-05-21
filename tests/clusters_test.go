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

package main

import (
	"fmt"
	"math/rand"
	"testing"
)

// generateClusterName generates cluster name in expected GUID-like format
func generateClusterName(t *testing.T) string {
	// random data generation
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		t.Fatal(err)
	}
	// format random data as proper cluster name
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// TestListClustersCommand function checks the 'list clusters' command
func TestListClustersCommand(t *testing.T) {
	// start the CLI client
	child := startCLI(t)

	// client needs to be shut down at the end of this test
	defer quitCLI(t, child)

	sendCommand(t, child, "list clusters")
	expectOutput(t, child, "List of clusters")

	// at the end, the standard prompt has to be displayed
	expectPrompt(t, child)
}

// TestAddClusterCommand function checks the 'add cluster' command followed by 'list clusters' one
func TestAddClusterCommand(t *testing.T) {
	// start the CLI client
	child := startCLI(t)

	// client needs to be shut down at the end of this test
	defer quitCLI(t, child)

	clusterName := generateClusterName(t)

	// try to add (register) new cluster
	command := fmt.Sprintf("add cluster %s", clusterName)
	sendCommand(t, child, command)
	expectOutput(t, child, "Cluster "+clusterName+" has been added")

	// now list clusters and check if the new cluster has been really added
	sendCommand(t, child, "list clusters")
	expectOutput(t, child, clusterName)

	// at the end, the standard prompt has to be displayed
	expectPrompt(t, child)
}
