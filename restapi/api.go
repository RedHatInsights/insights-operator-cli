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

package restapi

import (
	"github.com/redhatinsighs/insights-operator-cli/types"
)

// API represents API to the controller service. Normally it is implemented via REST API, but other methods can be used as well.
type API interface {
	// cluster related commands
	ReadListOfClusters() ([]types.Cluster, error)
	AddCluster(id string, name string) error
	DeleteCluster(clusterID string) error

	// configuration profiles related commands
	ReadListOfConfigurationProfiles() ([]types.ConfigurationProfile, error)
	ReadConfigurationProfile(profileID string) (*types.ConfigurationProfile, error)
	AddConfigurationProfile(username string, description string, configuration []byte) error
	DeleteConfigurationProfile(profileID string) error

	// configuration related commands
	ReadListOfConfigurations() ([]types.ClusterConfiguration, error)
	ReadClusterConfigurationByID(configurationID string) (*string, error)
	AddClusterConfiguration(username string, cluster string, reason string, description string, configuration []byte) error
	EnableClusterConfiguration(configurationID string) error
	DisableClusterConfiguration(configurationID string) error
	DeleteClusterConfiguration(configurationID string) error

	// trigger related commands
	ReadListOfTriggers() ([]types.Trigger, error)
	ReadTriggerByID(triggerID string) (*types.Trigger, error)
	AddTrigger(username string, clusterName string, reason string, link string) error
	DeleteTrigger(triggerID string) error
	ActivateTrigger(triggerID string) error
	DeactivateTrigger(triggerID string) error
}
