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

package commands

// Generated documentation is available at:
// https://pkg.go.dev/github.com/RedHatInsights/insights-operator-cli/commands
//
// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-cli/packages/commands/clusters.html

import (
	"fmt"

	"github.com/RedHatInsights/insights-operator-cli/restapi"
)

// ListOfClusters function displays list of clusters gathered via REST API call
// to the controller service. Just basic information about clusters are
// displayed - mainly its internal ID, an official ID, and a name.
func ListOfClusters(api restapi.API) {
	// try to read list of clusters and display error if something wrong
	// happens
	clusters, err := api.ReadListOfClusters()

	// check for any error
	if err != nil {
		// list of clusters operation failed for some reason
		fmt.Println(colorizer.Red("Error reading list of clusters"))
		fmt.Println(err)
		return
	}

	// TODO: handle empty list of clusters

	// list of clusters operation has been successful, let's display them
	fmt.Println(colorizer.Magenta("List of clusters"))
	fmt.Printf("%4s %4s %-s\n", "#", "ID", "Name")
	for i, cluster := range clusters {
		fmt.Printf("%4d %4d %-s\n", i, cluster.ID, cluster.Name)
	}
}

// DeleteClusterNoConfirm function deletes all info about selected cluster w/o
// asking for confirmation of this operation. Usually it function should not be
// called directly from user interface as some confirmation is required.
func DeleteClusterNoConfirm(api restapi.API, clusterID string) {
	DeleteCluster(api, clusterID, false)
}

// DeleteCluster function deletes all info about selected cluster from
// database. Before this operation is performed, user is ask if it is really
// required (this additional operation can be disabled by command line option).
func DeleteCluster(api restapi.API, clusterID string, askForConfirmation bool) {
	if askForConfirmation {
		// the client has been configured to ask for additional confirmation
		// display the confirmation dialog
		if !ProceedQuestion("All cluster configurations will be deleted") {
			// if answer is no, simply skip the rest of this function
			return
		}
	}

	// try to delete cluster and display error if something wrong happens
	err := api.DeleteCluster(clusterID)

	// check for any error
	if err != nil {
		// error has been detected during REST API call or during DB operation
		fmt.Println(colorizer.Red(errorCommunicationWithServiceErrorMessage))
		fmt.Println(err)
		return
	}

	// everything's ok, cluster has been deleted
	fmt.Println(colorizer.Blue("Cluster "+clusterID+" has been"), colorizer.Red(deleted))
}

// AddCluster function inserts new cluster info into the database via REST API
// call to insights operator controller service.
func AddCluster(api restapi.API, clusterName string) {
	// try to add new cluster and display error if something wrong happens
	err := api.AddCluster(clusterName)

	// check for any error
	if err != nil {
		// error has been detected during REST API call or during DB operation
		fmt.Println(colorizer.Red(errorCommunicationWithServiceErrorMessage))
		fmt.Println(err)
		return
	}

	// everything's ok, cluster has been added
	fmt.Println(colorizer.Blue("Cluster " + clusterName + " has been added"))
}
