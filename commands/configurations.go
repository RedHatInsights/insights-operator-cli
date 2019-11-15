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

func EnableClusterConfiguration(api restapi.Api, configurationId string) {
	err := api.EnableClusterConfiguration(configurationId)
	if err != nil {
		fmt.Println(Red("Error communicating with the service"))
		fmt.Println(err)
		return
	} else {
		fmt.Println(Blue("Configuration "+configurationId+" has been "), Green("enabled"))
	}
}

func DisableClusterConfiguration(api restapi.Api, configurationId string) {
	err := api.DisableClusterConfiguration(configurationId)
	if err != nil {
		fmt.Println(Red("Error communicating with the service"))
		fmt.Println(err)
		return
	} else {
		fmt.Println(Blue("Configuration "+configurationId+" has been "), Red("disabled"))
	}
}
