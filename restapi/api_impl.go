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

package restapi

// Generated documentation is available at:
// https://godoc.org/github.com/RedHatInsights/insights-operator-cli/restapi
//
// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-cli/packages/restapi/api_impl.html

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/RedHatInsights/insights-operator-cli/types"
	"net/http"
	"net/url"
)

// RestAPI is a structure representing instance of REST API
type RestAPI struct {
	controllerURL string
}

// NewRestAPI function is a constructor to construct new instance of REST API
func NewRestAPI(controllerURL string) RestAPI {
	return RestAPI{
		controllerURL: controllerURL,
	}
}

// ReadListOfClusters method reads list of clusters via the REST API
func (api RestAPI) ReadListOfClusters() ([]types.Cluster, error) {
	// structure for deserialized response
	clusters := types.ClustersResponse{}

	// construct URL to be used to access REST API endpoint
	url := api.controllerURL + APIPrefix + "client/cluster"

	// perform REST API call and check the result
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	// try to deserialize payload
	err = json.Unmarshal(body, &clusters)
	if err != nil {
		return nil, err
	}
	// and check for the status message in payload
	if clusters.Status != "ok" {
		return nil, fmt.Errorf(clusters.Status)
	}
	return clusters.Clusters, nil
}

// ReadListOfTriggers method reads list of triggers via the REST API
func (api RestAPI) ReadListOfTriggers() ([]types.Trigger, error) {
	// structure for deserialized response
	triggers := types.TriggersResponse{}

	// construct URL to be used to access REST API endpoint
	url := api.controllerURL + APIPrefix + "client/trigger"

	// perform REST API call and check the result
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	// try to deserialize payload
	err = json.Unmarshal(body, &triggers)
	if err != nil {
		return nil, err
	}
	// and check for the status message in payload
	if triggers.Status != "ok" {
		return nil, fmt.Errorf(triggers.Status)
	}
	return triggers.Triggers, nil
}

// ReadTriggerByID method reads trigger identified by its ID via the REST API
func (api RestAPI) ReadTriggerByID(triggerID string) (*types.Trigger, error) {
	// structure for deserialized response
	trigger := types.TriggerResponse{}

	// construct URL to be used to access REST API endpoint
	url := api.controllerURL + APIPrefix + "client/trigger/" + triggerID

	// perform REST API call and check the result
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	// try to deserialize payload
	err = json.Unmarshal(body, &trigger)
	if err != nil {
		return nil, err
	}
	// and check for the status message in payload
	if trigger.Status != "ok" {
		return nil, fmt.Errorf(trigger.Status)
	}
	return &trigger.Trigger, nil
}

// ReadListOfConfigurationProfiles method reads list of configuration profiles
// via the REST API
func (api RestAPI) ReadListOfConfigurationProfiles() ([]types.ConfigurationProfile, error) {
	// structure for deserialized response
	profiles := types.ConfigurationProfilesResponse{}

	// construct URL to be used to access REST API endpoint
	url := api.controllerURL + APIPrefix + "client/profile"

	// perform REST API call and check the result
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	// try to deserialize payload
	err = json.Unmarshal(body, &profiles)
	if err != nil {
		return nil, err
	}
	// and check for the status message in payload
	if profiles.Status != "ok" {
		return nil, fmt.Errorf(profiles.Status)
	}
	return profiles.Profiles, nil
}

// ReadListOfConfigurations method reads list of configuration via the REST API
func (api RestAPI) ReadListOfConfigurations() ([]types.ClusterConfiguration, error) {
	// structure for deserialized response
	configurations := types.ClusterConfigurationsResponse{}

	// construct URL to be used to access REST API endpoint
	url := api.controllerURL + APIPrefix + "client/configuration"

	// perform REST API call and check the result
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	// try to deserialize payload
	err = json.Unmarshal(body, &configurations)
	if err != nil {
		return nil, err
	}
	// and check for the status message in payload
	if configurations.Status != "ok" {
		return nil, fmt.Errorf(configurations.Status)
	}
	return configurations.Configurations, nil
}

// ReadConfigurationProfile method access the REST API endpoint to read
// selected configuration profile
func (api RestAPI) ReadConfigurationProfile(profileID string) (*types.ConfigurationProfile, error) {
	// structure for deserialized response
	profile := types.ConfigurationProfileResponse{}

	// construct URL to be used to access REST API endpoint
	url := api.controllerURL + APIPrefix + "client/profile/" + profileID

	// perform REST API call and check the result
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	// try to deserialize payload
	err = json.Unmarshal(body, &profile)
	if err != nil {
		return nil, err
	}
	// and check for the status message in payload
	if profile.Status != "ok" {
		return nil, fmt.Errorf(profile.Status)
	}
	return &profile.Profile, nil
}

// ReadClusterConfigurationByID method access the REST API endpoint to read
// cluster configuration for cluster defined by its ID
func (api RestAPI) ReadClusterConfigurationByID(configurationID string) (*string, error) {
	// structure for deserialized response
	configuration := types.ConfigurationResponse{}

	// construct URL to be used to access REST API endpoint
	url := api.controllerURL + APIPrefix + "client/configuration/" + configurationID

	// perform REST API call and check the result
	body, err := performReadRequest(url)
	if err != nil {
		return nil, err
	}

	// try to deserialize payload
	err = json.Unmarshal(body, &configuration)
	if err != nil {
		return nil, err
	}
	// and check for the status message in payload
	if configuration.Status != "ok" {
		return nil, fmt.Errorf(configuration.Status)
	}
	return &configuration.Configuration, nil
}

// EnableClusterConfiguration access the REST API endpoint to enable existing
// cluster configuration
func (api RestAPI) EnableClusterConfiguration(configurationID string) error {
	// construct URL to be used to access REST API endpoint
	url := api.controllerURL + APIPrefix + "client/configuration/" + configurationID + "/enable"

	// perform REST API call and return error code
	err := performWriteRequest(url, http.MethodPut, nil)
	return err
}

// DisableClusterConfiguration access the REST API endpoint to disable existing
// cluster configuration
func (api RestAPI) DisableClusterConfiguration(configurationID string) error {
	// construct URL to be used to access REST API endpoint
	url := api.controllerURL + APIPrefix + "client/configuration/" + configurationID + "/disable"

	// perform REST API call and return error code
	err := performWriteRequest(url, http.MethodPut, nil)
	return err
}

// DeleteClusterConfiguration access the REST API endpoint to delete existing
// cluster configuration
func (api RestAPI) DeleteClusterConfiguration(configurationID string) error {
	// construct URL to be used to access REST API endpoint
	url := api.controllerURL + APIPrefix + "client/configuration/" + configurationID

	// perform REST API call and return error code
	err := performWriteRequest(url, http.MethodDelete, nil)
	return err
}

// DeleteCluster access the REST API endpoint to delete/deregister existing
// cluster
func (api RestAPI) DeleteCluster(clusterID string) error {
	// construct URL to be used to access REST API endpoint
	url := api.controllerURL + APIPrefix + "client/cluster/" + clusterID

	// perform REST API call and return error code
	err := performWriteRequest(url, http.MethodDelete, nil)
	return err
}

// DeleteConfigurationProfile access the REST API endpoint to delete existing
// configuration profile
func (api RestAPI) DeleteConfigurationProfile(profileID string) error {
	// construct URL to be used to access REST API endpoint
	url := api.controllerURL + APIPrefix + "client/profile/" + profileID

	// perform REST API call and return error code
	err := performWriteRequest(url, http.MethodDelete, nil)
	return err
}

// AddCluster access the REST API endpoint to add/register new cluster
func (api RestAPI) AddCluster(name string) error {
	query := name
	// construct URL to be used to access REST API endpoint
	url := api.controllerURL + APIPrefix + "client/cluster/" + query

	// perform REST API call and return error code
	err := performWriteRequest(url, http.MethodPost, nil)
	return err
}

// AddConfigurationProfile access the REST API endpoint to add new
// configuration profile
func (api RestAPI) AddConfigurationProfile(username string, description string, configuration []byte) error {
	query := "username=" + url.QueryEscape(username) + "&description=" + url.QueryEscape(description)
	// construct URL to be used to access REST API endpoint
	url := api.controllerURL + APIPrefix + "client/profile?" + query

	// perform REST API call and return error code
	err := performWriteRequest(url, http.MethodPost, bytes.NewReader(configuration))
	return err
}

// AddClusterConfiguration access the REST API endpoint to add new cluster
// configuration
func (api RestAPI) AddClusterConfiguration(username string, cluster string, reason string, description string, configuration []byte) error {
	query := "username=" + url.QueryEscape(username) + "&reason=" + url.QueryEscape(reason) + "&description=" + url.QueryEscape(description)
	// construct URL to be used to access REST API endpoint
	url := api.controllerURL + APIPrefix + "client/cluster/" + url.PathEscape(cluster) + "/configuration/create?" + query

	// perform REST API call and return error code
	err := performWriteRequest(url, http.MethodPost, bytes.NewReader(configuration))
	return err
}

// AddTrigger access the REST API endpoint to add/register new trigger
func (api RestAPI) AddTrigger(username string, clusterName string, reason string, link string) error {
	query := "username=" + url.QueryEscape(username) + "&reason=" + url.QueryEscape(reason) + "&link=" + url.QueryEscape(link)
	// construct URL to be used to access REST API endpoint
	url := api.controllerURL + APIPrefix + "client/cluster/" + url.PathEscape(clusterName) + "/trigger/must-gather?" + query

	// perform REST API call and return error code
	err := performWriteRequest(url, http.MethodPost, nil)
	return err
}

// DeleteTrigger access the REST API endpoint to delete the selected trigger
func (api RestAPI) DeleteTrigger(triggerID string) error {
	// construct URL to be used to access REST API endpoint
	url := api.controllerURL + APIPrefix + "client/trigger/" + triggerID

	// perform REST API call and return error code
	err := performWriteRequest(url, http.MethodDelete, nil)
	return err
}

// ActivateTrigger access the REST API endpoint to activate the selected
// trigger
func (api RestAPI) ActivateTrigger(triggerID string) error {
	// construct URL to be used to access REST API endpoint
	url := api.controllerURL + APIPrefix + "client/trigger/" + triggerID + "/activate"

	// perform REST API call and return error code
	err := performWriteRequest(url, http.MethodPut, nil)
	return err
}

// DeactivateTrigger access the REST API endpoint to deactivate the selected
// trigger
func (api RestAPI) DeactivateTrigger(triggerID string) error {
	// construct URL to be used to access REST API endpoint
	url := api.controllerURL + APIPrefix + "client/trigger/" + triggerID + "/deactivate"

	// perform REST API call and return error code
	err := performWriteRequest(url, http.MethodPut, nil)
	return err
}
