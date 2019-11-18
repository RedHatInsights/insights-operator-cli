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
	"bytes"
	"encoding/json"
	"github.com/redhatinsighs/insights-operator-cli/types"
	"net/url"
)

type RestApi struct {
	controllerUrl string
}

func NewRestApi(controllerUrl string) RestApi {
	return RestApi{
		controllerUrl: controllerUrl,
	}
}

func (api RestApi) ReadListOfClusters() ([]types.Cluster, error) {
	clusters := []types.Cluster{}

	url := api.controllerUrl + API_PREFIX + "client/cluster"
	body, err := performReadRequest(url)

	err = json.Unmarshal(body, &clusters)
	if err != nil {
		return nil, err
	}
	return clusters, nil
}

func (api RestApi) ReadListOfTriggers() ([]types.Trigger, error) {
	var triggers []types.Trigger
	url := api.controllerUrl + API_PREFIX + "client/trigger"
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &triggers)
	if err != nil {
		return nil, err
	}
	return triggers, nil
}

func (api RestApi) ReadTriggerById(triggerId string) (*types.Trigger, error) {
	var trigger types.Trigger
	url := api.controllerUrl + API_PREFIX + "client/trigger/" + triggerId
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &trigger)
	if err != nil {
		return nil, err
	}
	return &trigger, nil
}

func (api RestApi) ReadListOfConfigurationProfiles() ([]types.ConfigurationProfile, error) {
	profiles := []types.ConfigurationProfile{}

	url := api.controllerUrl + API_PREFIX + "client/profile"
	body, err := performReadRequest(url)

	err = json.Unmarshal(body, &profiles)
	if err != nil {
		return nil, err
	}
	return profiles, nil
}

func (api RestApi) ReadListOfConfigurations() ([]types.ClusterConfiguration, error) {
	configurations := []types.ClusterConfiguration{}

	url := api.controllerUrl + API_PREFIX + "client/configuration"
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &configurations)
	if err != nil {
		return nil, err
	}
	return configurations, nil
}

func (api RestApi) ReadConfigurationProfile(profileId string) (*types.ConfigurationProfile, error) {
	var profile types.ConfigurationProfile
	url := api.controllerUrl + API_PREFIX + "client/profile/" + profileId
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (api RestApi) ReadClusterConfigurationById(configurationId string) (*string, error) {
	url := api.controllerUrl + API_PREFIX + "client/configuration/" + configurationId
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	str := string(body)
	return &str, nil
}

func (api RestApi) EnableClusterConfiguration(configurationId string) error {
	url := api.controllerUrl + API_PREFIX + "client/configuration/" + configurationId + "/enable"
	err := performWriteRequest(url, "PUT", nil)
	return err
}

func (api RestApi) DisableClusterConfiguration(configurationId string) error {
	url := api.controllerUrl + API_PREFIX + "client/configuration/" + configurationId + "/disable"
	err := performWriteRequest(url, "PUT", nil)
	return err
}

func (api RestApi) DeleteClusterConfiguration(configurationId string) error {
	url := api.controllerUrl + API_PREFIX + "client/configuration/" + configurationId
	err := performWriteRequest(url, "DELETE", nil)
	return err
}

func (api RestApi) DeleteCluster(clusterId string) error {
	url := api.controllerUrl + API_PREFIX + "client/cluster/" + clusterId
	err := performWriteRequest(url, "DELETE", nil)
	return err
}

func (api RestApi) DeleteConfigurationProfile(profileId string) error {
	url := api.controllerUrl + API_PREFIX + "client/profile/" + profileId
	err := performWriteRequest(url, "DELETE", nil)
	return err
}

func (api RestApi) AddCluster(id string, name string) error {
	query := id + "/" + name
	url := api.controllerUrl + API_PREFIX + "client/cluster/" + query
	err := performWriteRequest(url, "POST", nil)
	return err
}

func (api RestApi) AddConfigurationProfile(username string, description string, configuration []byte) error {
	query := "username=" + url.QueryEscape(username) + "&description=" + url.QueryEscape(description)
	url := api.controllerUrl + API_PREFIX + "client/profile?" + query
	err := performWriteRequest(url, "POST", bytes.NewReader(configuration))
	return err
}

func (api RestApi) AddClusterConfiguration(username string, cluster string, reason string, description string, configuration []byte) error {
	query := "username=" + url.QueryEscape(username) + "&reason=" + url.QueryEscape(reason) + "&description=" + url.QueryEscape(description)
	url := api.controllerUrl + API_PREFIX + "client/cluster/" + url.PathEscape(cluster) + "/configuration?" + query
	err := performWriteRequest(url, "POST", bytes.NewReader(configuration))
	return err
}

func (api RestApi) AddTrigger(username string, clusterName string, reason string, link string) error {
	query := "username=" + url.QueryEscape(username) + "&reason=" + url.QueryEscape(reason) + "&link=" + url.QueryEscape(link)
	url := api.controllerUrl + API_PREFIX + "client/cluster/" + url.PathEscape(clusterName) + "/trigger/must-gather?" + query
	err := performWriteRequest(url, "POST", nil)
	return err
}

func (api RestApi) DeleteTrigger(triggerId string) error {
	url := api.controllerUrl + API_PREFIX + "client/trigger/" + triggerId
	err := performWriteRequest(url, "DELETE", nil)
	return err
}

func (api RestApi) ActivateTrigger(triggerId string) error {
	url := api.controllerUrl + API_PREFIX + "client/trigger/" + triggerId + "/activate"
	err := performWriteRequest(url, "PUT", nil)
	return err
}

func (api RestApi) DeactivateTrigger(triggerId string) error {
	url := api.controllerUrl + API_PREFIX + "client/trigger/" + triggerId + "/deactivate"
	err := performWriteRequest(url, "PUT", nil)
	return err
}
