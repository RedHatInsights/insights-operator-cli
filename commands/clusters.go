/*
Copyright © 2019 Red Hat, Inc.

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

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/redhatinsighs/insights-operator-cli/restapi"
)

// ListOfClusters displays list of clusters gathered via REST API call to the controller service
func ListOfClusters(api restapi.API) {
	clusters, err := api.ReadListOfClusters()
	if err != nil {
		fmt.Println(colorizer.Red("Error reading list of clusters"))
		fmt.Println(err)
		return
	}

	fmt.Println(colorizer.Magenta("List of clusters"))
	fmt.Printf("%4s %4s %-s\n", "#", "ID", "Name")
	for i, cluster := range clusters {
		fmt.Printf("%4d %4d %-s\n", i, cluster.ID, cluster.Name)
	}
}

// DeleteCluster deletes all info about selected cluster from database
func DeleteCluster(api restapi.API, clusterID string) {
	if !ProceedQuestion("All cluster configurations will be deleted") {
		return
	}

	err := api.DeleteCluster(clusterID)
	if err != nil {
		fmt.Println(colorizer.Red("Error communicating with the service"))
		fmt.Println(err)
		return
	}

	// everything's ok
	fmt.Println(colorizer.Blue("Cluster "+clusterID+" has been"), colorizer.Red("deleted"))
}

// AddCluster inserts new cluster info into the database
func AddCluster(api restapi.API) {
	id := prompt.Input("ID: ", LoginCompleter)
	if id == "" {
		fmt.Println(colorizer.Red("Cancelled"))
		return
	}

	name := prompt.Input("name: ", LoginCompleter)
	if name == "" {
		fmt.Println(colorizer.Red("Cancelled"))
		return
	}

	err := api.AddCluster(id, name)
	if err != nil {
		fmt.Println(colorizer.Red("Error communicating with the service"))
		fmt.Println(err)
		return
	}

	// everything's ok, controller has been registered
	fmt.Println(colorizer.Blue("Controller has been registered"))
}
