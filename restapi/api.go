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

type Api interface {
	// cluster related commands
	ReadListOfClusters() ([]types.Cluster, error)
	AddCluster(id string, name string) error
	DeleteCluster(clusterId string) error

	// configuration profiles related commands
	ReadListOfConfigurationProfiles() ([]types.ConfigurationProfile, error)
	ReadConfigurationProfile(profileId string) (*types.ConfigurationProfile, error)
	AddConfigurationProfile(username string, description string, configuration []byte) error
	DeleteConfigurationProfile(profileId string) error

	// configuration related commands
	ReadListOfConfigurations() ([]types.ClusterConfiguration, error)
	ReadClusterConfigurationById(configurationId string) (*string, error)
	AddClusterConfiguration(username string, cluster string, reason string, description string, configuration []byte) error
	EnableClusterConfiguration(configurationId string) error
	DisableClusterConfiguration(configurationId string) error
	DeleteClusterConfiguration(configurationId string) error

	// trigger related commands
	ReadListOfTriggers() ([]types.Trigger, error)
	ReadTriggerById(triggerId string) (*types.Trigger, error)
	AddTrigger(username string, clusterName string, reason string, link string) error
	DeleteTrigger(triggerId string) error
	ActivateTrigger(triggerId string) error
	DeactivateTrigger(triggerId string) error
}
