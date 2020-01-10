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

package commands

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/logrusorgru/aurora"
	"github.com/redhatinsighs/insights-operator-cli/restapi"
	"io/ioutil"
	"strings"
)

// ListOfConfigurations displays list of all configurations gathered via REST API call to controller service
func ListOfConfigurations(api restapi.API, filter string) {
	// TODO: filter in query?
	configurations, err := api.ReadListOfConfigurations()
	if err != nil {
		fmt.Println(colorizer.Red("Error reading list of configurations"))
		fmt.Println(err)
		return
	}

	fmt.Println(colorizer.Magenta("List of configurations for all clusters"))
	fmt.Printf("%4s %4s %4s    %-20s %-20s %-10s %-12s %s\n", "#", "ID", "Profile", "Cluster", "Changed at", "Changed by", "Active", "Reason")
	for i, configuration := range configurations {
		// poor man's filtering
		if strings.Contains(configuration.Cluster, filter) {
			var active aurora.Value
			if configuration.Active == "1" {
				active = colorizer.Green("yes")
			} else {
				active = colorizer.Red("no")
			}
			changedAt := configuration.ChangedAt[0:19]
			fmt.Printf("%4d %4d %4s       %-20s %-20s %-10s %-12s %s\n", i, configuration.ID, configuration.Configuration, configuration.Cluster, changedAt, configuration.ChangedBy, active, configuration.Reason)
		}
	}
}

// EnableClusterConfiguration enables the selected cluster configuration in the controller service
func EnableClusterConfiguration(api restapi.API, configurationID string) {
	err := api.EnableClusterConfiguration(configurationID)
	if err != nil {
		fmt.Println(colorizer.Red("Error communicating with the service"))
		fmt.Println(err)
		return
	}

	// everything's ok
	fmt.Println(colorizer.Blue("Configuration "+configurationID+" has been"), colorizer.Green("enabled"))
}

// DisableClusterConfiguration disables the selected cluster configuration in the controller service
func DisableClusterConfiguration(api restapi.API, configurationID string) {
	err := api.DisableClusterConfiguration(configurationID)
	if err != nil {
		fmt.Println(colorizer.Red("Error communicating with the service"))
		fmt.Println(err)
		return
	}

	// everything's ok
	fmt.Println(colorizer.Blue("Configuration "+configurationID+" has been"), colorizer.Red("disabled"))
}

// DescribeConfiguration displays additional information about selected configuration
func DescribeConfiguration(api restapi.API, clusterID string) {
	configuration, err := api.ReadClusterConfigurationByID(clusterID)
	if err != nil {
		fmt.Println(colorizer.Red("Error reading cluster configuration"))
		fmt.Println(err)
		return
	}

	fmt.Println(colorizer.Magenta("Configuration for cluster " + clusterID))
	fmt.Println(*configuration)
}

// DeleteClusterConfiguration deletes selected cluster configuration from database
func DeleteClusterConfiguration(api restapi.API, configurationID string) {
	err := api.DeleteClusterConfiguration(configurationID)
	if err != nil {
		fmt.Println(colorizer.Red("Error communicating with the service"))
		fmt.Println(err)
		return
	}

	// everything's ok, configuration has been deleted
	fmt.Println(colorizer.Blue("Configuration "+configurationID+" has been"), colorizer.Red("deleted"))
}

// AddClusterConfiguration creates a new cluster configuration
func AddClusterConfiguration(api restapi.API, username string) {
	if username == "" {
		fmt.Println(colorizer.Red("Not logged in"))
		return
	}

	cluster := prompt.Input("cluster: ", LoginCompleter)
	if cluster == "" {
		fmt.Println(colorizer.Red("Cancelled"))
		return
	}

	reason := prompt.Input("reason: ", LoginCompleter)
	if reason == "" {
		fmt.Println(colorizer.Red("Cancelled"))
		return
	}

	description := prompt.Input("description: ", LoginCompleter)
	if description == "" {
		fmt.Println(colorizer.Red("Cancelled"))
		return
	}

	// TODO: make the directory fully configurable
	err := FillInConfigurationList("configurations")
	if err != nil {
		fmt.Println(colorizer.Red("Cannot read any configuration file"))
		fmt.Println(err)
	}

	configurationFileName := prompt.Input("configuration file (TAB to complete): ", ConfigFileCompleter)
	if configurationFileName == "" {
		fmt.Println(colorizer.Red("Cancelled"))
		return
	}

	// TODO: make the directory fully configurable
	configuration, err := ioutil.ReadFile("configurations/" + configurationFileName)
	if err != nil {
		fmt.Println(colorizer.Red("Cannot read configuration file"))
		fmt.Println(err)
	}

	err = api.AddClusterConfiguration(username, cluster, reason, description, configuration)
	if err != nil {
		fmt.Println(colorizer.Red("Error communicating with the service"))
		fmt.Println(err)
		return
	}

	// everything's ok, configuration has been created
	fmt.Println(colorizer.Blue("Configuration has been created"))
}
