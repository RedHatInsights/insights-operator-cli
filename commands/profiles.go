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
// https://redhatinsights.github.io/insights-operator-cli/packages/commands/profiles.html

import (
	"fmt"
	"github.com/RedHatInsights/insights-operator-cli/restapi"
	"github.com/c-bata/go-prompt"
	"io/ioutil"
)

// ListOfProfiles function displays list of configuration profiles gathered via
// REST API call made to controller service.
func ListOfProfiles(api restapi.API) {
	// try to read list of configuration profiles and display error when
	// something wrong happens
	profiles, err := api.ReadListOfConfigurationProfiles()
	if err != nil {
		// in case of error just print the error message
		fmt.Println(colorizer.Red("Error reading list of configuration profiles"))
		fmt.Println(err)
		return
	}

	// REST API call returns data
	fmt.Println(colorizer.Magenta("List of configuration profiles"))
	fmt.Printf("%4s %4s %-20s %-20s %s\n", "#", "ID", changedAt, changedBy, "Description")

	// list all profiles
	for i, profile := range profiles {
		// update timestamps not to contain irrelevant parts
		changedAt := profile.ChangedAt[0:19]
		fmt.Printf("%4d %4d %-20s %-20s %-s\n", i, profile.ID, changedAt, profile.ChangedBy, profile.Description)
	}
}

// DescribeProfile function displays additional information about selected
// profile
func DescribeProfile(api restapi.API, profileID string) {
	// try to read configuration profile identified by its ID and display
	// error when something wrong happens
	profile, err := api.ReadConfigurationProfile(profileID)
	if err != nil {
		// in case of error just print the error message
		fmt.Println(colorizer.Red("Error reading configuration profile"))
		fmt.Println(err)
		return
	}

	// print the configuration profile
	fmt.Println(colorizer.Magenta("Configuration profile"))
	fmt.Println(profile.Configuration)
}

// DeleteConfigurationProfileNoConfirm function deletes the configuration
// profile selected by its ID w/o asking for confirmation
func DeleteConfigurationProfileNoConfirm(api restapi.API, profileID string) {
	// directly call REST API endpoint to delete configuration profile
	DeleteConfigurationProfile(api, profileID, false)
}

// DeleteConfigurationProfile function deletes the profile selected by its ID
func DeleteConfigurationProfile(api restapi.API, profileID string, askForConfirmation bool) {
	if askForConfirmation {
		// ask user whether he/she really want to delete configuration profile
		if !ProceedQuestion("All configurations based on this profile will be deleted") {
			return
		}
	}

	// try to delete configuration profile identified by its ID and display
	// error when something wrong happens
	err := api.DeleteConfigurationProfile(profileID)
	if err != nil {
		// in case of error just print the error message
		fmt.Println(colorizer.Red(errorCommunicationWithServiceErrorMessage))
		fmt.Println(err)
		return
	}

	// everything's ok, configuration profile has been deleted
	fmt.Println(colorizer.Blue("Configuration profile "+profileID+" has been"), colorizer.Red(deleted))
}

// AddConfigurationProfile function adds the profile to database
func AddConfigurationProfile(api restapi.API, username string) {
	// check if user is already loged in
	if username == "" {
		fmt.Println(colorizer.Red(notLoggedIn))
		return
	}

	// ask for description of configuration profile
	description := prompt.Input(descriptionPrompt, LoginCompleter)
	if description == "" {
		fmt.Println(colorizer.Red(operationCancelled))
		return
	}

	// TODO: make the directory fully configurable
	err := FillInConfigurationList(configurationsDirectory)
	if err != nil {
		fmt.Println(colorizer.Red(cannotReadAnyConfigurationFileErrorMessage))
		fmt.Println(err)
	}

	// let the user select the file
	configurationFileName := prompt.Input(configurationFilePrompt, ConfigFileCompleter)
	if configurationFileName == "" {
		// in case of error just print the error message
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
		// in case of error just print the error message
		fmt.Println(colorizer.Red(errorCommunicationWithServiceErrorMessage))
		fmt.Println(err)
		return
	}

	// everything's ok, configuration profile can be created
	AddConfigurationProfileImpl(api, username, description, configurationFileName)
}

// AddConfigurationProfileImpl function adds the profile to database
func AddConfigurationProfileImpl(api restapi.API, username, description, configurationFileName string) {
	// TODO: make the directory fully configurable
	// disable "G304 (CWE-22): Potential file inclusion via variable"
	configuration, err := ioutil.ReadFile("configurations/" + configurationFileName) // #nosec G304
	if err != nil {
		fmt.Println(colorizer.Red(cannotReadConfigurationFileErrorMessage))
		fmt.Println(err)
	}

	// try to add configuration profile and display error when something
	// wrong happens
	err = api.AddConfigurationProfile(username, description, configuration)
	if err != nil {
		// in case of error just print the error message
		fmt.Println(colorizer.Red(errorCommunicationWithServiceErrorMessage))
		fmt.Println(err)
		return
	}

	// everything's ok, configuration profile has been created
	fmt.Println(colorizer.Blue("Configuration profile has been created"))
}
