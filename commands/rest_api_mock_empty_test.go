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

// Mock object used by unit tests for REST API. This mocked object
// returns empty data from every method.

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-cli/packages/commands/rest_api_mock_empty_test.html

import (
	"github.com/RedHatInsights/insights-operator-cli/types"
)

// RestAPIMockEmpty is an implementation of mocked REST API that returns empty
// data structures in all cases.
type RestAPIMockEmpty struct {
}

// ReadListOfClusters reads mocked empty list of clusters via the REST API.
// This is a mock implementation of original method.
func (api RestAPIMockEmpty) ReadListOfClusters() ([]types.Cluster, error) {
	// data structure to be returned
	clusters := []types.Cluster{}

	// return mocked response with empty data structure
	return clusters, nil
}

// ReadListOfTriggers reads mocked empty list of triggers via the REST API.
// This is a mock implementation of original method.
func (api RestAPIMockEmpty) ReadListOfTriggers() ([]types.Trigger, error) {
	// data structure to be returned
	triggers := []types.Trigger{}

	// return mocked response with empty data structure
	return triggers, nil
}

// ReadTriggerByID reads trigger identified by its ID via the REST API.
// This is a mock implementation of original method.
func (api RestAPIMockEmpty) ReadTriggerByID(triggerID string) (*types.Trigger, error) {
	// data structure to be returned
	trigger := types.Trigger{}

	// return mocked response with empty data structure
	return &trigger, nil
}

// ReadListOfConfigurationProfiles reads mocked empty list of configuration
// profiles via the REST API.
// This is a mock implementation of original method.
func (api RestAPIMockEmpty) ReadListOfConfigurationProfiles() ([]types.ConfigurationProfile, error) {
	// data structure to be returned
	profiles := []types.ConfigurationProfile{}

	// return mocked response with empty data structure
	return profiles, nil
}

// ReadListOfConfigurations reads mocked empty list of configuration via the
// REST API.
// This is a mock implementation of original method.
func (api RestAPIMockEmpty) ReadListOfConfigurations() ([]types.ClusterConfiguration, error) {
	// data structure to be returned
	configurations := []types.ClusterConfiguration{}

	// return mocked response with empty data structure
	return configurations, nil
}

// ReadConfigurationProfile access the REST API endpoint to read selected
// configuration profile.
// This is a mock implementation of original method.
func (api RestAPIMockEmpty) ReadConfigurationProfile(profileID string) (*types.ConfigurationProfile, error) {
	// data structure to be returned
	profile := types.ConfigurationProfile{}

	// return mocked response with empty data structure
	return &profile, nil
}

// ReadClusterConfigurationByID access the REST API endpoint to read cluster
// configuration for cluster defined by its ID.
// This is a mock implementation of original method.
func (api RestAPIMockEmpty) ReadClusterConfigurationByID(configurationID string) (*string, error) {
	// data structure to be returned
	configuration := ""

	// return mocked response with empty data structure
	return &configuration, nil
}

// EnableClusterConfiguration access the REST API endpoint to enable existing
// cluster configuration.
// This is a mock implementation of original method.
func (api RestAPIMockEmpty) EnableClusterConfiguration(configurationID string) error {
	// return mocked response with empty data structure
	return nil
}

// DisableClusterConfiguration access the REST API endpoint to disable existing
// cluster configuration.
// This is a mock implementation of original method.
func (api RestAPIMockEmpty) DisableClusterConfiguration(configurationID string) error {
	// return mocked response with empty data structure
	return nil
}

// DeleteClusterConfiguration access the REST API endpoint to delete existing
// cluster configuration.
// This is a mock implementation of original method.
func (api RestAPIMockEmpty) DeleteClusterConfiguration(configurationID string) error {
	// return mocked response with empty data structure
	return nil
}

// DeleteCluster access the REST API endpoint to delete/deregister existing
// cluster.
// This is a mock implementation of original method.
func (api RestAPIMockEmpty) DeleteCluster(clusterID string) error {
	// return mocked response with empty data structure
	return nil
}

// DeleteConfigurationProfile access the REST API endpoint to delete existing
// configuration profile.
// This is a mock implementation of original method.
func (api RestAPIMockEmpty) DeleteConfigurationProfile(profileID string) error {
	// return mocked response with empty data structure
	return nil
}

// AddCluster access the REST API endpoint to add/register new cluster.
// This is a mock implementation of original method.
func (api RestAPIMockEmpty) AddCluster(name string) error {
	// return mocked response with empty data structure
	return nil
}

// AddConfigurationProfile access the REST API endpoint to add new
// configuration profile.
// This is a mock implementation of original method.
func (api RestAPIMockEmpty) AddConfigurationProfile(username string, description string, configuration []byte) error {
	// return mocked response with empty data structure
	return nil
}

// AddClusterConfiguration access the REST API endpoint to add new cluster
// configuration.
// This is a mock implementation of original method.
func (api RestAPIMockEmpty) AddClusterConfiguration(username string, cluster string, reason string, description string, configuration []byte) error {
	// return mocked response with empty data structure
	return nil
}

// AddTrigger access the REST API endpoint to add/register new trigger.
// This is a mock implementation of original method.
func (api RestAPIMockEmpty) AddTrigger(username string, clusterName string, reason string, link string) error {
	// return mocked response with empty data structure
	return nil
}

// DeleteTrigger access the REST API endpoint to delete the selected trigger.
// This is a mock implementation of original method.
func (api RestAPIMockEmpty) DeleteTrigger(triggerID string) error {
	// return mocked response with empty data structure
	return nil
}

// ActivateTrigger access the REST API endpoint to activate the selected trigger.
// This is a mock implementation of original method.
func (api RestAPIMockEmpty) ActivateTrigger(triggerID string) error {
	// return mocked response with empty data structure
	return nil
}

// DeactivateTrigger access the REST API endpoint to deactivate the selected
// trigger.
// This is a mock implementation of original method.
func (api RestAPIMockEmpty) DeactivateTrigger(triggerID string) error {
	// return mocked response with empty data structure
	return nil
}
