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
	"github.com/logrusorgru/aurora"
	"github.com/redhatinsighs/insights-operator-cli/restapi"
)

// ListOfProfiles displays list of configuration profiles gathered via REST API call to controller service
func ListOfProfiles(api restapi.Api) {
	profiles, err := api.ReadListOfConfigurationProfiles()
	if err != nil {
		fmt.Println(aurora.Red("Error reading list of configuration profiles"))
		fmt.Println(err)
		return
	}

	fmt.Println(aurora.Magenta("List of configuration profiles"))
	fmt.Printf("%4s %4s %-20s %-20s %s\n", "#", "ID", "Changed at", "Changed by", "Description")
	for i, profile := range profiles {
		changedAt := profile.ChangedAt[0:19]
		fmt.Printf("%4d %4d %-20s %-20s %-s\n", i, profile.Id, changedAt, profile.ChangedBy, profile.Description)
	}
}
