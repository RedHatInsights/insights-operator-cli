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

package commands_test

// Mock object used by unit tests for REST API

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-cli/packages/commands/rest_api_mock_test.html

import (
	"github.com/RedHatInsights/insights-operator-cli/types"
)

// RestAPIMock structure is an implementation of mocked REST API
type RestAPIMock struct {
}

// ReadListOfClusters reads mocked list of clusters via the REST API.
// This is a mock implementation of original method.
func (api RestAPIMock) ReadListOfClusters() ([]types.Cluster, error) {
	// data structure to be returned
	clusters := []types.Cluster{
		{
			ID:   0,
			Name: "c8590f31-e97e-4b85-b506-c45ce1911a12"},
		{
			ID:   0,
			Name: "ffffffff-ffff-ffff-ffff-ffffffffffff"},
		{
			ID:   1,
			Name: "00000000-0000-0000-0000-000000000000"}}

	// return mocked response
	return clusters, nil
}

// ReadListOfTriggers reads mocked list of triggers via the REST API.
// This is a mock implementation of original method.
func (api RestAPIMock) ReadListOfTriggers() ([]types.Trigger, error) {
	// data structure to be returned
	triggers := []types.Trigger{
		{
			ID:          0,
			Type:        "must-gather",
			Cluster:     "ffffffff-ffff-ffff-ffff-ffffffffffff",
			Reason:      "we need to run must-gather",
			Link:        "https://www.webpagetest.org/",
			TriggeredAt: "2020-01-01T00:00:00",
			TriggeredBy: "tester",
			AckedAt:     "1970-01-01T00:00:00",
			Parameters:  "",
			Active:      1},
		{
			ID:          1,
			Type:        "must-gather",
			Cluster:     "ffffffff-ffff-ffff-ffff-ffffffffffff",
			Reason:      "we need to run must-gather",
			Link:        "https://www.webpagetest.org/",
			TriggeredAt: "2020-01-01T00:00:00",
			TriggeredBy: "tester",
			AckedAt:     "2020-01-01T00:00:00",
			Parameters:  "",
			Active:      0},
		{
			ID:          2,
			Type:        "must-gather",
			Cluster:     "00000000-0000-0000-0000-000000000000",
			Reason:      "we need to run must-gather",
			Link:        "https://www.webpagetest.org/",
			TriggeredAt: "2020-01-01T00:00:00",
			TriggeredBy: "tester",
			AckedAt:     "1970-01-01T00:00:00",
			Parameters:  "",
			Active:      1},
		{
			ID:          3,
			Type:        "different-trigger",
			Cluster:     "00000000-0000-0000-0000-000000000000",
			Reason:      "we need to run must-gather",
			Link:        "https://www.webpagetest.org/",
			TriggeredAt: "2020-01-01T00:00:00",
			TriggeredBy: "tester",
			AckedAt:     "2020-01-01T00:00:00",
			Parameters:  "",
			Active:      0},
	}

	// return mocked response
	return triggers, nil
}

// ReadTriggerByID reads trigger identified by its ID via the REST API.
// This is a mock implementation of original method.
func (api RestAPIMock) ReadTriggerByID(triggerID string) (*types.Trigger, error) {
	switch triggerID {
	case "0":
		// data structure to be returned
		trigger := types.Trigger{
			ID:          0,
			Type:        "must-gather",
			Cluster:     "ffffffff-ffff-ffff-ffff-ffffffffffff",
			Reason:      "we need to run must-gather",
			Link:        "https://www.webpagetest.org/",
			TriggeredAt: "2020-01-01T00:00:00",
			TriggeredBy: "tester",
			AckedAt:     "1970-01-01T00:00:00",
			Parameters:  "",
			Active:      1}
		return &trigger, nil
	case "1":
		// data structure to be returned
		trigger := types.Trigger{
			ID:          1,
			Type:        "must-gather",
			Cluster:     "ffffffff-ffff-ffff-ffff-ffffffffffff",
			Reason:      "we need to run must-gather",
			Link:        "https://www.webpagetest.org/",
			TriggeredAt: "2020-01-01T00:00:00",
			TriggeredBy: "tester",
			AckedAt:     "2020-01-02T00:00:00",
			Parameters:  "",
			Active:      0}
		return &trigger, nil
	case "2":
		// data structure to be returned
		trigger := types.Trigger{
			ID:          2,
			Type:        "another-one",
			Cluster:     "00000000-0000-0000-0000-000000000000",
			Reason:      "something else",
			Link:        "https://www.webpagetest.org/",
			TriggeredAt: "2020-01-01T00:00:00",
			TriggeredBy: "tester",
			AckedAt:     "2020-01-02T00:00:00",
			Parameters:  "-a -W",
			Active:      0}
		return &trigger, nil
	}
	trigger := types.TriggerResponse{}

	// return mocked response
	return &trigger.Trigger, nil
}

// ReadListOfConfigurationProfiles reads mocked list of configuration profiles
// via the REST API.
// This is a mock implementation of original method.
func (api RestAPIMock) ReadListOfConfigurationProfiles() ([]types.ConfigurationProfile, error) {
	// data structure to be returned
	profiles := []types.ConfigurationProfile{
		{
			ID:            0,
			Configuration: "",
			ChangedAt:     "2020-01-01T00:00:00",
			ChangedBy:     "tester",
			Description:   "default configuration profile"},
		{
			ID:            1,
			Configuration: "",
			ChangedAt:     "2020-01-01T00:00:00",
			ChangedBy:     "tester",
			Description:   "another configuration profile"},
	}

	// return mocked response
	return profiles, nil
}

// ReadListOfConfigurations reads mocked list of configuration via the REST API.
// This is a mock implementation of original method.
func (api RestAPIMock) ReadListOfConfigurations() ([]types.ClusterConfiguration, error) {
	// data structure to be returned
	configurations := []types.ClusterConfiguration{
		{
			ID:            0,
			Cluster:       "0",
			Configuration: "0",
			ChangedAt:     "2020-01-01T00:00:00",
			ChangedBy:     "tester",
			Active:        "1",
			Reason:        "configuration1"},
		{
			ID:            1,
			Cluster:       "0",
			Configuration: "1",
			ChangedAt:     "2020-01-01T00:00:00",
			ChangedBy:     "tester",
			Active:        "1",
			Reason:        "configuration2"},
		{
			ID:            2,
			Cluster:       "0",
			Configuration: "2",
			ChangedAt:     "2020-01-01T00:00:00",
			ChangedBy:     "tester",
			Active:        "0",
			Reason:        "configuration3"},
	}

	// return mocked response
	return configurations, nil
}

// ReadConfigurationProfile access the REST API endpoint to read selected
// configuration profile.
// This is a mock implementation of original method.
func (api RestAPIMock) ReadConfigurationProfile(profileID string) (*types.ConfigurationProfile, error) {
	if profileID == "0" {
		// data structure to be returned
		profile := types.ConfigurationProfile{
			ID:            0,
			Configuration: "*configuration*",
			ChangedAt:     "2020-01-01T00:00:00",
			ChangedBy:     "tester",
			Description:   "empty configuration"}
		return &profile, nil
	}
	profile := types.ConfigurationProfile{}

	// return mocked response
	return &profile, nil
}

// ReadClusterConfigurationByID access the REST API endpoint to read cluster
// configuration for cluster defined by its ID.
// This is a mock implementation of original method.
func (api RestAPIMock) ReadClusterConfigurationByID(configurationID string) (*string, error) {
	if configurationID == "0" {
		// data structure to be returned
		configuration := "configuration#0"

		// return mocked response
		return &configuration, nil
	}

	// return mocked response
	return nil, nil
}

// EnableClusterConfiguration access the REST API endpoint to enable existing
// cluster configuration.
// This is a mock implementation of original method.
func (api RestAPIMock) EnableClusterConfiguration(configurationID string) error {
	// return mocked response
	return nil
}

// DisableClusterConfiguration access the REST API endpoint to disable existing
// cluster configuration.
// This is a mock implementation of original method.
func (api RestAPIMock) DisableClusterConfiguration(configurationID string) error {
	// return mocked response
	return nil
}

// DeleteClusterConfiguration access the REST API endpoint to delete existing
// cluster configuration.
// This is a mock implementation of original method.
func (api RestAPIMock) DeleteClusterConfiguration(configurationID string) error {
	// return mocked response
	return nil
}

// DeleteCluster access the REST API endpoint to delete/deregister existing
// cluster.
// This is a mock implementation of original method.
func (api RestAPIMock) DeleteCluster(clusterID string) error {
	// return mocked response
	return nil
}

// DeleteConfigurationProfile access the REST API endpoint to delete existing
// configuration profile.
// This is a mock implementation of original method.
func (api RestAPIMock) DeleteConfigurationProfile(profileID string) error {
	// return mocked response
	return nil
}

// AddCluster access the REST API endpoint to add/register new cluster.
// This is a mock implementation of original method.
func (api RestAPIMock) AddCluster(name string) error {
	// return mocked response
	return nil
}

// AddConfigurationProfile access the REST API endpoint to add new
// configuration profile.
// This is a mock implementation of original method.
func (api RestAPIMock) AddConfigurationProfile(username, description string, configuration []byte) error {
	// return mocked response
	return nil
}

// AddClusterConfiguration access the REST API endpoint to add new cluster.
// configuration
// This is a mock implementation of original method.
func (api RestAPIMock) AddClusterConfiguration(username, cluster, reason, description string, configuration []byte) error {
	// return mocked response
	return nil
}

// AddTrigger access the REST API endpoint to add/register new trigger.
// This is a mock implementation of original method.
func (api RestAPIMock) AddTrigger(username, clusterName, reason, link string) error {
	// return mocked response
	return nil
}

// DeleteTrigger access the REST API endpoint to delete the selected trigger.
// This is a mock implementation of original method.
func (api RestAPIMock) DeleteTrigger(triggerID string) error {
	// return mocked response
	return nil
}

// ActivateTrigger access the REST API endpoint to activate the selected
// trigger.
// This is a mock implementation of original method.
func (api RestAPIMock) ActivateTrigger(triggerID string) error {
	// return mocked response
	return nil
}

// DeactivateTrigger access the REST API endpoint to deactivate the selected
// trigger.
// This is a mock implementation of original method.
func (api RestAPIMock) DeactivateTrigger(triggerID string) error {
	// return mocked response
	return nil
}
