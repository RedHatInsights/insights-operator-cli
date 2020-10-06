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
// https://redhatinsights.github.io/insights-operator-cli/packages/commands/triggers.html

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/logrusorgru/aurora"
	"github.com/redhatinsighs/insights-operator-cli/restapi"
)

// ListOfTriggers function displays list of triggers (including must-gather
// one) gathered via REST API call to controller service.
func ListOfTriggers(api restapi.API) {
	// TODO: filter in query?
	// try to read list of triggers and display error message if anything
	// wrong happens
	triggers, err := api.ReadListOfTriggers()
	if err != nil {
		fmt.Println(colorizer.Red("Error reading list of triggers"))
		fmt.Println(err)
		return
	}

	fmt.Println(colorizer.Magenta("List of triggers for all clusters"))
	fmt.Printf("%4s %4s %-16s    %-20s %-20s %-12s %-12s %s\n", "#", "ID", "Type", "Cluster", "Triggered at", "Triggered by", "Active", "Acked at")
	for i, trigger := range triggers {
		var active aurora.Value
		if trigger.Active == 1 {
			active = colorizer.Green(conditionSet)
		} else {
			active = colorizer.Red("no")
		}
		triggeredAt := trigger.TriggeredAt[0:19]
		ackedAt := trigger.AckedAt[0:19]
		fmt.Printf("%4d %4d %-16s    %-20s %-20s %-12s %-12s %s\n", i, trigger.ID, trigger.Type, trigger.Cluster, triggeredAt, trigger.TriggeredBy, active, ackedAt)
	}
}

// DescribeTrigger function displays additional information about selected
// trigger.
func DescribeTrigger(api restapi.API, triggerID string) {
	// try to read trigger idintified by its ID and display error message
	// if anything wrong happens
	trigger, err := api.ReadTriggerByID(triggerID)
	if err != nil {
		fmt.Println(colorizer.Red("Error reading selected trigger"))
		fmt.Println(err)
		return
	}

	var active aurora.Value
	if trigger.Active == 1 {
		active = colorizer.Green(conditionSet)
	} else {
		active = colorizer.Red("no")
	}

	triggeredAt := trigger.TriggeredAt[0:19]
	ackedAt := trigger.AckedAt[0:19]

	var ttype aurora.Value
	if trigger.Type == "must-gather" {
		ttype = colorizer.Blue(trigger.Type)
	} else {
		ttype = colorizer.Magenta(trigger.Type)
	}

	fmt.Println(colorizer.Magenta("Trigger info"))
	fmt.Printf("ID:            %d\n", trigger.ID)
	fmt.Printf("Type:          %s\n", ttype)
	fmt.Printf("Cluster:       %s\n", trigger.Cluster)
	fmt.Printf("Triggered at:  %s\n", triggeredAt)
	fmt.Printf("Triggered by:  %s\n", trigger.TriggeredBy)
	fmt.Printf("Active:        %s\n", active)
	fmt.Printf("Acked at:      %s\n", ackedAt)
}

// AddTrigger function adds new trigger for a cluster.
func AddTrigger(api restapi.API, username string) {
	if username == "" {
		fmt.Println(colorizer.Red(notLoggedIn))
		return
	}

	clusterName := prompt.Input("cluster name: ", LoginCompleter)
	reason := prompt.Input("reason: ", LoginCompleter)
	link := prompt.Input("link: ", LoginCompleter)

	AddTriggerImpl(api, username, clusterName, reason, link)
}

// AddTriggerImpl function calls REST API to add a new trigger into the
// database.
func AddTriggerImpl(api restapi.API, username string, clusterName string, reason string, link string) {
	// try to add a new trigger and display error message if anything wrong
	// happens
	err := api.AddTrigger(username, clusterName, reason, link)
	if err != nil {
		fmt.Println(errorCommunicationWithServiceErrorMessage)
		fmt.Println(err)
		return
	}

	// everything's ok, trigger has been created
	fmt.Println(colorizer.Blue("Trigger has been created"))
}

// DeleteTrigger function deletes specified trigger.
func DeleteTrigger(api restapi.API, triggerID string) {
	// try to delete trigger idintified by its ID and display error message
	// if anything wrong happens
	err := api.DeleteTrigger(triggerID)
	if err != nil {
		fmt.Println(colorizer.Red(errorCommunicationWithServiceErrorMessage))
		fmt.Println(err)
		return
	}

	// everything's ok, trigger has been deleted
	fmt.Println(colorizer.Blue("Trigger "+triggerID+" has been"), colorizer.Red(deleted))
}

// ActivateTrigger function activates specified trigger.
func ActivateTrigger(api restapi.API, triggerID string) {
	// try to activate trigger idintified by its ID and display error
	// message if anything wrong happens
	err := api.ActivateTrigger(triggerID)
	if err != nil {
		fmt.Println(colorizer.Red(errorCommunicationWithServiceErrorMessage))
		fmt.Println(err)
		return
	}

	// everything's ok, trigger has been activated
	fmt.Println(colorizer.Blue("Trigger "+triggerID+" has been"), colorizer.Green("activated"))
}

// DeactivateTrigger deactivates specified trigger
func DeactivateTrigger(api restapi.API, triggerID string) {
	// try to deactivate trigger idintified by its ID and display error
	// message if anything wrong happens
	err := api.DeactivateTrigger(triggerID)
	if err != nil {
		fmt.Println(colorizer.Red(errorCommunicationWithServiceErrorMessage))
		fmt.Println(err)
		return
	}

	// everything's ok, trigger has been deactivated
	fmt.Println(colorizer.Blue("Trigger "+triggerID+" has been"), colorizer.Green("deactivated"))
}
