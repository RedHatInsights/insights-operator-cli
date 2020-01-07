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
	"errors"
	"github.com/redhatinsighs/insights-operator-cli/types"
)

// RestAPIMockErrors is an implementation of mocked REST API that returns errors in all cases.
type RestAPIMockErrors struct {
}

// ReadListOfClusters returns an error as its last return value
func (api RestAPIMockErrors) ReadListOfClusters() ([]types.Cluster, error) {
	clusters := []types.Cluster{}
	return clusters, errors.New("ReadListOfClusters error")
}

// ReadListOfTriggers returns an error as its last return value
func (api RestAPIMockErrors) ReadListOfTriggers() ([]types.Trigger, error) {
	triggers := []types.Trigger{}
	return triggers, errors.New("ReadListofTrigger error")
}

// ReadTriggerByID returns an error as its last return value
func (api RestAPIMockErrors) ReadTriggerByID(triggerID string) (*types.Trigger, error) {
	trigger := types.Trigger{}
	return &trigger, errors.New("ReadTriggerByID error")
}

// ReadListOfConfigurationProfiles returns an error as its last return value
func (api RestAPIMockErrors) ReadListOfConfigurationProfiles() ([]types.ConfigurationProfile, error) {
	profiles := []types.ConfigurationProfile{}
	return profiles, errors.New("ReadListOfConfigurationProfiles error")
}

// ReadListOfConfigurations returns an error as its last return value
func (api RestAPIMockErrors) ReadListOfConfigurations() ([]types.ClusterConfiguration, error) {
	configurations := []types.ClusterConfiguration{}
	return configurations, errors.New("ReadListOfConfigurations error")
}

// ReadConfigurationProfile returns an error as its last return value
func (api RestAPIMockErrors) ReadConfigurationProfile(profileID string) (*types.ConfigurationProfile, error) {
	profile := types.ConfigurationProfile{}
	return &profile, errors.New("ReadConfigurationProfile error")
}

// ReadClusterConfigurationByID returns an error as its last return value
func (api RestAPIMockErrors) ReadClusterConfigurationByID(configurationID string) (*string, error) {
	configuration := ""
	return &configuration, errors.New("ReadClusterConfigurationByID error")
}

// EnableClusterConfiguration returns an error as its last return value
func (api RestAPIMockErrors) EnableClusterConfiguration(configurationID string) error {
	return errors.New("EnableClusterConfiguration error")
}

// DisableClusterConfiguration returns an error as its last return value
func (api RestAPIMockErrors) DisableClusterConfiguration(configurationID string) error {
	return errors.New("DisableClusterConfiguration error")
}

// DeleteClusterConfiguration returns an error as its last return value
func (api RestAPIMockErrors) DeleteClusterConfiguration(configurationID string) error {
	return errors.New("DeleteClusterConfiguration error")
}

// DeleteCluster returns an error as its last return value
func (api RestAPIMockErrors) DeleteCluster(clusterID string) error {
	return errors.New("DeleteCluster error")
}

// DeleteConfigurationProfile returns an error as its last return value
func (api RestAPIMockErrors) DeleteConfigurationProfile(profileID string) error {
	return errors.New("DeleteConfigurationProfile error")
}

// AddCluster returns an error as its last return value
func (api RestAPIMockErrors) AddCluster(name string) error {
	return errors.New("AddCluster error")
}

// AddConfigurationProfile access the REST API endpoint to add new configuration profile
func (api RestAPIMockErrors) AddConfigurationProfile(username string, description string, configuration []byte) error {
	return errors.New("AddConfigurationProfile error")
}

// AddClusterConfiguration returns an error as its last return value
func (api RestAPIMockErrors) AddClusterConfiguration(username string, cluster string, reason string, description string, configuration []byte) error {
	return errors.New("AddClusterConfiguration error")
}

// AddTrigger returns an error as its last return value
func (api RestAPIMockErrors) AddTrigger(username string, clusterName string, reason string, link string) error {
	return errors.New("AddTrigger error")
}

// DeleteTrigger returns an error as its last return value
func (api RestAPIMockErrors) DeleteTrigger(triggerID string) error {
	return errors.New("DeleteTrigger error")
}

// ActivateTrigger returns an error as its last return value
func (api RestAPIMockErrors) ActivateTrigger(triggerID string) error {
	return errors.New("ActivateTrigger error")
}

// DeactivateTrigger returns an error as its last return value
func (api RestAPIMockErrors) DeactivateTrigger(triggerID string) error {
	return errors.New("DeactivateTrigger error")
}
