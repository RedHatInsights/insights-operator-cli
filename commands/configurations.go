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

package commands

// Generated documentation is available at:
// https://pkg.go.dev/github.com/RedHatInsights/insights-operator-cli/commands
//
// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-cli/packages/commands/configurations.html

import (
	"fmt"
	"github.com/RedHatInsights/insights-operator-cli/restapi"
	"github.com/c-bata/go-prompt"
	"github.com/logrusorgru/aurora"
	"os"
	"strings"
)

// configFileDirectory constant contains path to directory containing all
// configuration files for this tool.
const configFileDirectory = "configurations/"

const configurationsDirectory = "configurations"

// ListOfConfigurations function displays list of all configurations gathered
// via REST API call to the Controller Service.
func ListOfConfigurations(api restapi.API, filter string) {
	// TODO: filter in query?
	// try to read list of configurations and display error if something
	// wrong happens
	configurations, err := api.ReadListOfConfigurations()
	if err != nil {
		fmt.Println(colorizer.Red(ErrorReadingListOfConfigurations))
		fmt.Println(err)
		return
	}

	// list all configurations returned in HTTP response
	fmt.Println(colorizer.Magenta("List of configurations for all clusters"))
	fmt.Printf("%4s %4s %4s    %-20s %-20s %-10s %-12s %s\n", "#", "ID", "Profile", clusterUUID, changedAt, changedBy, activeTrigger, "Reason")
	for i, configuration := range configurations {
		// perform poor man's filtering on client side
		if strings.Contains(configuration.Cluster, filter) {
			var active aurora.Value
			if configuration.Active == "1" {
				active = colorizer.Green(conditionSet)
			} else {
				active = colorizer.Red("no")
			}
			changedAt := configuration.ChangedAt[0:19]
			fmt.Printf("%4d %4d %4s       %-20s %-20s %-10s %-12s %s\n", i, configuration.ID, configuration.Configuration, configuration.Cluster, changedAt, configuration.ChangedBy, active, configuration.Reason)
		}
	}
}

// EnableClusterConfiguration function enables the selected cluster
// configuration in the controller service via REST API call.
func EnableClusterConfiguration(api restapi.API, configurationID string) {
	// try to enable cluster configuration and display error if something
	// wrong happens
	err := api.EnableClusterConfiguration(configurationID)
	if err != nil {
		fmt.Println(colorizer.Red(ErrorCommunicationWithServiceErrorMessage))
		fmt.Println(err)
		return
	}

	// everything's ok
	fmt.Println(colorizer.Blue("Configuration "+configurationID+" has been"), colorizer.Green("enabled"))
}

// DisableClusterConfiguration function disables the selected cluster
// configuration in the controller service via REST API call.
func DisableClusterConfiguration(api restapi.API, configurationID string) {
	// try to disable cluster configuration and display error if something
	// wrong happens
	err := api.DisableClusterConfiguration(configurationID)
	if err != nil {
		fmt.Println(colorizer.Red(ErrorCommunicationWithServiceErrorMessage))
		fmt.Println(err)
		return
	}

	// everything's ok
	fmt.Println(colorizer.Blue("Configuration "+configurationID+" has been"), colorizer.Red("disabled"))
}

// DescribeConfiguration function displays additional information about
// selected configuration read via REST API call.
func DescribeConfiguration(api restapi.API, clusterID string) {
	// try to read cluster configuration by using its ID and display error
	// if something wrong happens
	configuration, err := api.ReadClusterConfigurationByID(clusterID)
	if err != nil {
		fmt.Println(colorizer.Red("Error reading cluster configuration"))
		fmt.Println(err)
		return
	}

	fmt.Println(colorizer.Magenta("Configuration for cluster " + clusterID))
	fmt.Println(*configuration)
}

// DeleteClusterConfiguration function deletes selected cluster configuration
// from database via REST API call.
func DeleteClusterConfiguration(api restapi.API, configurationID string) {
	// try to delete cluster configuration and display error if something
	// wrong happens
	err := api.DeleteClusterConfiguration(configurationID)
	if err != nil {
		fmt.Println(colorizer.Red(ErrorCommunicationWithServiceErrorMessage))
		fmt.Println(err)
		return
	}

	// everything's ok, configuration has been deleted
	fmt.Println(colorizer.Blue("Configuration "+configurationID+" has been"), colorizer.Red(deleted))
}

// AddClusterConfiguration function asks for all information needed to create
// new cluster configuration, all done via REST API call.
func AddClusterConfiguration(api restapi.API, username string) {
	if username == "" {
		fmt.Println(colorizer.Red(notLoggedIn))
		return
	}

	// ask user about cluster ID
	cluster := prompt.Input("cluster: ", LoginCompleter)
	if cluster == "" {
		fmt.Println(colorizer.Red(operationCancelled))
		return
	}

	// ask user about reason
	reason := prompt.Input(reasonPrompt, LoginCompleter)
	if reason == "" {
		fmt.Println(colorizer.Red(operationCancelled))
		return
	}

	// ask user about description
	description := prompt.Input(descriptionPrompt, LoginCompleter)
	if description == "" {
		fmt.Println(colorizer.Red(operationCancelled))
		return
	}

	// TODO: make the directory fully configurable
	err := FillInConfigurationList(configurationsDirectory)
	if err != nil {
		fmt.Println(colorizer.Red(CannotReadAnyConfigurationFileErrorMessage))
		fmt.Println(err)
	}

	// user need to select the configuration file
	configurationFileName := prompt.Input(configurationFilePrompt, ConfigFileCompleter)
	if configurationFileName == "" {
		fmt.Println(colorizer.Red(operationCancelled))
		return
	}

	// try to add new cluster configuration
	AddClusterConfigurationImpl(api, username, cluster, reason, description, configurationFileName)
}

func pathToConfigFile(directory, filename string) string {
	return directory + filename
}

// AddClusterConfigurationImpl function creates a new cluster configuration.
func AddClusterConfigurationImpl(api restapi.API, username, cluster, reason, description, configurationFileName string) {
	// TODO: make the directory fully configurable
	configuration, err := os.ReadFile(pathToConfigFile(configFileDirectory, configurationFileName))
	if err != nil {
		fmt.Println(colorizer.Red(CannotReadConfigurationFileErrorMessage))
		fmt.Println(err)
		return
	}

	// try to add cluster configuration and display error if something
	// wrong happens
	err = api.AddClusterConfiguration(username, cluster, reason, description, configuration)
	if err != nil {
		fmt.Println(colorizer.Red(ErrorCommunicationWithServiceErrorMessage))
		fmt.Println(err)
		return
	}

	// everything's ok, configuration has been created
	fmt.Println(colorizer.Blue("Configuration has been created"))
}
