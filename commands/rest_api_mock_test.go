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

package commands_test

import (
	"github.com/redhatinsighs/insights-operator-cli/types"
)

// RestAPIMockErrors is an implementation of mocked REST API
type RestAPIMock struct {
}

// ReadListOfClusters reads mocked list of clusters via the REST API
func (api RestAPIMock) ReadListOfClusters() ([]types.Cluster, error) {
	clusters := []types.Cluster{
		types.Cluster{
			ID:   0,
			Name: "c8590f31-e97e-4b85-b506-c45ce1911a12"},
		types.Cluster{
			ID:   0,
			Name: "ffffffff-ffff-ffff-ffff-ffffffffffff"},
		types.Cluster{
			ID:   1,
			Name: "00000000-0000-0000-0000-000000000000"}}
	return clusters, nil
}

// ReadListOfTriggers reads mocked list of triggers via the REST API
func (api RestAPIMock) ReadListOfTriggers() ([]types.Trigger, error) {
	triggers := []types.Trigger{
		types.Trigger{
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
		types.Trigger{
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
		types.Trigger{
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
		types.Trigger{
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
	return triggers, nil
}

// ReadTriggerByID reads trigger identified by its ID via the REST API
func (api RestAPIMock) ReadTriggerByID(triggerID string) (*types.Trigger, error) {
	switch triggerID {
	case "0":
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
	return &trigger.Trigger, nil
}

// ReadListOfConfigurationProfiles reads mocked list of configuration profiles via the REST API
func (api RestAPIMock) ReadListOfConfigurationProfiles() ([]types.ConfigurationProfile, error) {
	profiles := []types.ConfigurationProfile{
		types.ConfigurationProfile{
			ID:            0,
			Configuration: "",
			ChangedAt:     "2020-01-01T00:00:00",
			ChangedBy:     "tester",
			Description:   "default configuration profile"},
		types.ConfigurationProfile{
			ID:            1,
			Configuration: "",
			ChangedAt:     "2020-01-01T00:00:00",
			ChangedBy:     "tester",
			Description:   "another configuration profile"},
	}
	return profiles, nil
}

// ReadListOfConfigurations reads mocked list of configuration via the REST API
func (api RestAPIMock) ReadListOfConfigurations() ([]types.ClusterConfiguration, error) {
	configurations := []types.ClusterConfiguration{
		types.ClusterConfiguration{
			ID:            0,
			Cluster:       "0",
			Configuration: "0",
			ChangedAt:     "2020-01-01T00:00:00",
			ChangedBy:     "tester",
			Active:        "1",
			Reason:        "configuration1"},
		types.ClusterConfiguration{
			ID:            1,
			Cluster:       "0",
			Configuration: "1",
			ChangedAt:     "2020-01-01T00:00:00",
			ChangedBy:     "tester",
			Active:        "1",
			Reason:        "configuration2"},
		types.ClusterConfiguration{
			ID:            2,
			Cluster:       "0",
			Configuration: "2",
			ChangedAt:     "2020-01-01T00:00:00",
			ChangedBy:     "tester",
			Active:        "0",
			Reason:        "configuration3"},
	}
	return configurations, nil
}

// ReadConfigurationProfile access the REST API endpoint to read selected configuration profile
func (api RestAPIMock) ReadConfigurationProfile(profileID string) (*types.ConfigurationProfile, error) {
	if profileID == "0" {
		profile := types.ConfigurationProfile{
			ID:            0,
			Configuration: "*configuration*",
			ChangedAt:     "2020-01-01T00:00:00",
			ChangedBy:     "tester",
			Description:   "empty configuration"}
		return &profile, nil
	}
	profile := types.ConfigurationProfile{}
	return &profile, nil
}

// ReadClusterConfigurationByID access the REST API endpoint to read cluster configuration for cluster defined by its ID
func (api RestAPIMock) ReadClusterConfigurationByID(configurationID string) (*string, error) {
	if configurationID == "0" {
		configuration := "configuration#0"
		return &configuration, nil
	}
	return nil, nil
}

// EnableClusterConfiguration access the REST API endpoint to enable existing cluster configuration
func (api RestAPIMock) EnableClusterConfiguration(configurationID string) error {
	return nil
}

// DisableClusterConfiguration access the REST API endpoint to disable existing cluster configuration
func (api RestAPIMock) DisableClusterConfiguration(configurationID string) error {
	return nil
}

// DeleteClusterConfiguration access the REST API endpoint to delete existing cluster configuration
func (api RestAPIMock) DeleteClusterConfiguration(configurationID string) error {
	return nil
}

// DeleteCluster access the REST API endpoint to delete/deregister existing cluster
func (api RestAPIMock) DeleteCluster(clusterID string) error {
	return nil
}

// DeleteConfigurationProfile access the REST API endpoint to delete existing configuration profile
func (api RestAPIMock) DeleteConfigurationProfile(profileID string) error {
	return nil
}

// AddCluster access the REST API endpoint to add/register new cluster
func (api RestAPIMock) AddCluster(name string) error {
	return nil
}

// AddConfigurationProfile access the REST API endpoint to add new configuration profile
func (api RestAPIMock) AddConfigurationProfile(username string, description string, configuration []byte) error {
	return nil
}

// AddClusterConfiguration access the REST API endpoint to add new cluster configuration
func (api RestAPIMock) AddClusterConfiguration(username string, cluster string, reason string, description string, configuration []byte) error {
	return nil
}

// AddTrigger access the REST API endpoint to add/register new trigger
func (api RestAPIMock) AddTrigger(username string, clusterName string, reason string, link string) error {
	return nil
}

// DeleteTrigger access the REST API endpoint to delete the selected trigger
func (api RestAPIMock) DeleteTrigger(triggerID string) error {
	return nil
}

// ActivateTrigger access the REST API endpoint to activate the selected trigger
func (api RestAPIMock) ActivateTrigger(triggerID string) error {
	return nil
}

// DeactivateTrigger access the REST API endpoint to deactivate the selected trigger
func (api RestAPIMock) DeactivateTrigger(triggerID string) error {
	return nil
}
