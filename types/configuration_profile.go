/*
Copyright © 2019, 2020, 2021 Red Hat, Inc.

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

package types

// Generated documentation is available at:
// https://pkg.go.dev/github.com/RedHatInsights/insights-operator-cli/types
//
// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-cli/packages/types/configuration_profile.html

// ConfigurationProfile structure represents configuration profile record in
// the controller service.
//
//	ID: unique key
//	Configuration: a JSON structure stored in a string
//	ChangeAt: username of admin that created or updated the configuration
//	ChangeBy: timestamp of the last configuration change
//	Description: a string with any comment(s) about the configuration
type ConfigurationProfile struct {
	ID            int    `json:"id"`
	Configuration string `json:"configuration"`
	ChangedAt     string `json:"changed_at"`
	ChangedBy     string `json:"changed_by"`
	Description   string `json:"description"`
}

// ConfigurationProfilesResponse structure represents response of controller
// service to configuration profiles request.
//
//	Status: status of response
//	Profiles: list of configuration profiles
type ConfigurationProfilesResponse struct {
	Status   string                 `json:"status"`
	Profiles []ConfigurationProfile `json:"profiles"`
}

// ConfigurationProfileResponse structure represents response of controller
// service to single configuration profile request.
//
//	Status: status of response
//	Profile: single configuration profile
type ConfigurationProfileResponse struct {
	Status  string               `json:"status"`
	Profile ConfigurationProfile `json:"profile"`
}
