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

// Generated documentation is available at:
// https://godoc.org/github.com/RedHatInsights/insights-operator-cli/commands
//
// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-cli/packages/commands/profiles.html

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/redhatinsighs/insights-operator-cli/restapi"
	"io/ioutil"
)

// ListOfProfiles displays list of configuration profiles gathered via REST API
// call to controller service
func ListOfProfiles(api restapi.API) {
	// try to read list of configuration profiles and display error when
	// something wrong happens
	profiles, err := api.ReadListOfConfigurationProfiles()
	if err != nil {
		fmt.Println(colorizer.Red("Error reading list of configuration profiles"))
		fmt.Println(err)
		return
	}

	fmt.Println(colorizer.Magenta("List of configuration profiles"))
	fmt.Printf("%4s %4s %-20s %-20s %s\n", "#", "ID", changedAt, "Changed by", "Description")
	for i, profile := range profiles {
		changedAt := profile.ChangedAt[0:19]
		fmt.Printf("%4d %4d %-20s %-20s %-s\n", i, profile.ID, changedAt, profile.ChangedBy, profile.Description)
	}
}

// DescribeProfile displays additional information about selected profile
func DescribeProfile(api restapi.API, profileID string) {
	// try to read configuration profile identified by its ID and display
	// error when something wrong happens
	profile, err := api.ReadConfigurationProfile(profileID)
	if err != nil {
		fmt.Println(colorizer.Red("Error reading configuration profile"))
		fmt.Println(err)
		return
	}

	fmt.Println(colorizer.Magenta("Configuration profile"))
	fmt.Println(profile.Configuration)
}

// DeleteConfigurationProfileNoConfirm deletes the profile selected by its ID
// w/o asking for confirmation
func DeleteConfigurationProfileNoConfirm(api restapi.API, profileID string) {
	DeleteConfigurationProfile(api, profileID, false)
}

// DeleteConfigurationProfile deletes the profile selected by its ID
func DeleteConfigurationProfile(api restapi.API, profileID string, askForConfirmation bool) {
	if askForConfirmation {
		if !ProceedQuestion("All configurations based on this profile will be deleted") {
			return
		}
	}

	// try to delete configuration profile identified by its ID and display
	// error when something wrong happens
	err := api.DeleteConfigurationProfile(profileID)
	if err != nil {
		fmt.Println(colorizer.Red(errorCommunicationWithServiceErrorMessage))
		fmt.Println(err)
		return
	}

	// everything's ok, configuration profile has been deleted
	fmt.Println(colorizer.Blue("Configuration profile "+profileID+" has been"), colorizer.Red(deleted))
}

// AddConfigurationProfile adds the profile to database
func AddConfigurationProfile(api restapi.API, username string) {
	if username == "" {
		fmt.Println(colorizer.Red(notLoggedIn))
		return
	}

	description := prompt.Input("description: ", LoginCompleter)
	if description == "" {
		fmt.Println(colorizer.Red(operationCancelled))
		return
	}

	// TODO: make the directory fully configurable
	err := FillInConfigurationList("configurations")
	if err != nil {
		fmt.Println(colorizer.Red(cannotReadAnyConfigurationFileErrorMessage))
		fmt.Println(err)
	}

	configurationFileName := prompt.Input(configurationFilePrompt, ConfigFileCompleter)
	if configurationFileName == "" {
		fmt.Println(colorizer.Red(operationCancelled))
		return
	}

	// TODO: make the directory fully configurable
	configuration, err := ioutil.ReadFile(pathToConfigFile(configFileDirectory, configurationFileName))
	if err != nil {
		fmt.Println(colorizer.Red(cannotReadConfigurationFileErrorMessage))
		fmt.Println(err)
	}

	// try to add configuration profile and display error when something
	// wrong happens
	err = api.AddConfigurationProfile(username, description, configuration)
	if err != nil {
		fmt.Println(colorizer.Red(errorCommunicationWithServiceErrorMessage))
		fmt.Println(err)
		return
	}

	// everything's ok, configuration profile can be created
	AddConfigurationProfileImpl(api, username, description, configurationFileName)
}

// AddConfigurationProfileImpl adds the profile to database
func AddConfigurationProfileImpl(api restapi.API, username string, description string, configurationFileName string) {
	// TODO: make the directory fully configurable
	// disable "G304 (CWE-22): Potential file inclusion via variable"
	// #nosec G304
	configuration, err := ioutil.ReadFile("configurations/" + configurationFileName)
	if err != nil {
		fmt.Println(colorizer.Red(cannotReadConfigurationFileErrorMessage))
		fmt.Println(err)
	}

	// try to add configuration profile and display error when something
	// wrong happens
	err = api.AddConfigurationProfile(username, description, configuration)
	if err != nil {
		fmt.Println(colorizer.Red(errorCommunicationWithServiceErrorMessage))
		fmt.Println(err)
		return
	}

	// everything's ok, configuration profile has been created
	fmt.Println(colorizer.Blue("Configuration profile has been created"))
}
