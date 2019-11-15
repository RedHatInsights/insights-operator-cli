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
	. "github.com/logrusorgru/aurora"
	"github.com/redhatinsighs/insights-operator-cli/restapi"
)

func ListOfTriggers(api restapi.Api) {
	// TODO: filter in query?
	triggers, err := api.ReadListOfTriggers()
	if err != nil {
		fmt.Println(Red("Error reading list of triggers"))
		fmt.Println(err)
		return
	}

	fmt.Println(Magenta("List of triggers for all clusters"))
	fmt.Printf("%4s %4s %-16s    %-20s %-20s %-12s %-12s %s\n", "#", "ID", "Type", "Cluster", "Triggered at", "Triggered by", "Active", "Acked at")
	for i, trigger := range triggers {
		var active Value
		if trigger.Active == 1 {
			active = Green("yes")
		} else {
			active = Red("no")
		}
		triggeredAt := trigger.TriggeredAt[0:19]
		ackedAt := trigger.AckedAt[0:19]
		fmt.Printf("%4d %4d %-16s    %-20s %-20s %-12s %-12s %s\n", i, trigger.Id, trigger.Type, trigger.Cluster, triggeredAt, trigger.TriggeredBy, active, ackedAt)
	}
}
