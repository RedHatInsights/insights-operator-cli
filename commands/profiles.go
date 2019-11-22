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
	"io/ioutil"
)

// ListOfProfiles displays list of configuration profiles gathered via REST API call to controller service
func ListOfProfiles(api restapi.Api) {
	profiles, err := api.ReadListOfConfigurationProfiles()
	if err != nil {
		fmt.Println(colorizer.Red("Error reading list of configuration profiles"))
		fmt.Println(err)
		return
	}

	fmt.Println(colorizer.Magenta("List of configuration profiles"))
	fmt.Printf("%4s %4s %-20s %-20s %s\n", "#", "ID", "Changed at", "Changed by", "Description")
	for i, profile := range profiles {
		changedAt := profile.ChangedAt[0:19]
		fmt.Printf("%4d %4d %-20s %-20s %-s\n", i, profile.ID, changedAt, profile.ChangedBy, profile.Description)
	}
}

// DescribeProfile displays additional information about selected profile
func DescribeProfile(api restapi.Api, profileID string) {
	profile, err := api.ReadConfigurationProfile(profileID)
	if err != nil {
		fmt.Println(colorizer.Red("Error reading configuration profile"))
		fmt.Println(err)
		return
	}

	fmt.Println(colorizer.Magenta("Configuration profile"))
	fmt.Println(profile.Configuration)
}

// DeleteConfigurationProfile deletes the profile selected by its ID
func DeleteConfigurationProfile(api restapi.Api, profileID string) {
	if !ProceedQuestion("All configurations based on this profile will be deleted") {
		return
	}

	err := api.DeleteConfigurationProfile(profileID)
	if err != nil {
		fmt.Println(colorizer.Red("Error communicating with the service"))
		fmt.Println(err)
		return
	}

	// everything's ok, configuration profile has been deleted
	fmt.Println(colorizer.Blue("Configuration profile "+profileID+" has been "), colorizer.Red("deleted"))
}

// AddConfigurationProfile adds the profile to database
func AddConfigurationProfile(api restapi.Api, username string) {
	if username == "" {
		fmt.Println(colorizer.Red("Not logged in"))
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

	err = api.AddConfigurationProfile(username, description, configuration)
	if err != nil {
		fmt.Println(colorizer.Red("Error communicating with the service"))
		fmt.Println(err)
		return
	}

	// everything's ok, configuration profile has been created
	fmt.Println(colorizer.Blue("Configuration profile has been created"))
}
