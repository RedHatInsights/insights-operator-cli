/*
Copyright © 2020 Red Hat, Inc.

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

package restapi_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/redhatinsighs/insights-operator-cli/restapi"
	"github.com/redhatinsighs/insights-operator-cli/types"

	"testing"
)

const (
	RESTApiPrefix         = "/api/v1/client/"
	ReadClustersURL       = RESTApiPrefix + "cluster"
	ReadTriggersURL       = RESTApiPrefix + "trigger"
	ReadProfilesURL       = RESTApiPrefix + "profile"
	ReadConfigurationsURL = RESTApiPrefix + "configuration"

	StatusOKJSON    = `{"status":"ok"}`
	StatusErrorJSON = `{"status":"error"}`
	ImproperJSON    = `this is not proper JSON`
)

// mockedHTTPServer prepares new instance of testing HTTP server
func mockedHTTPServer(handler func(responseWriter http.ResponseWriter, request *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}

// checkURL checks that the request URL is expected for given usage of HTTP server
func checkURL(t *testing.T, request *http.Request, expectedURL string) {
	if request == nil {
		t.Error("Ptr. to request structure is nil")
	}
	// check the actual URL if it's the same as expected one
	if request.URL.String() != expectedURL {
		t.Error("Invalid URL:", request.URL.String(), "expected:", expectedURL)
	}

}

// checkMethod checks if the method in HTTP request is appropriate
func checkMethod(t *testing.T, request *http.Request, method string) {
	if request.Method != method {
		t.Error("Inapropriate: method used to call REST API:", request.Method)
	}
}

// writeBody writes a given text into the response that is to be send to receiver
func writeBody(responseWriter http.ResponseWriter, body string) {
	responseWriter.Write([]byte(body))
}

// standardHandlerImpl is an implementation of handler that checks URL and when it's expected send a response
// with body that contains a body filled with given response string
func standardHandlerImpl(t *testing.T, expectedURL, responseStr string) func(responseWriter http.ResponseWriter, request *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		checkMethod(t, request, "GET")
		// check if the URL is expected one
		checkURL(t, request, expectedURL)
		// send response to be tested later
		writeBody(responseWriter, responseStr)
	}
}

// standardHandlerForMethodImpl is an implementation of handler that checks URL and when it's expected send a response
// with body that contains a body filled with given response string. Additionally used method is checked as well.
func standardHandlerForMethodImpl(t *testing.T, expectedURL, method, responseStr string) func(responseWriter http.ResponseWriter, request *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		checkMethod(t, request, method)
		// check if the URL is expected one
		checkURL(t, request, expectedURL)
		// send response to be tested later
		writeBody(responseWriter, responseStr)
	}
}

// TestReadListOfClustersEmptyList check the method ReadListOfClusters
func TestReadListOfClustersEmptyList(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerImpl(t, ReadClustersURL, StatusOKJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	clusters, err := api.ReadListOfClusters()
	if err != nil {
		t.Fatal(err)
	}

	if len(clusters) != 0 {
		t.Fatal("Expected empty list of clusters")
	}
}

// TestReadListOfClustersOneCluster check the method ReadListOfClusters
func TestReadListOfClustersOneCluster(t *testing.T) {
	// start a local HTTP server
	const responseAsString = `
	{
		"clusters": [{"ID":0,"Name":"Name"}],
		"status":"ok"
	}`
	server := mockedHTTPServer(standardHandlerImpl(t, ReadClustersURL, responseAsString))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	clusters, err := api.ReadListOfClusters()
	if err != nil {
		t.Fatal(err)
	}

	if len(clusters) != 1 {
		t.Fatal("Expected list with one cluster only")
	}
}

// TestReadListOfClustersErrorStatus check the method ReadListOfClusters
func TestReadListOfClustersErrorStatus(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerImpl(t, ReadClustersURL, StatusErrorJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadListOfClusters()
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadListOfClustersEmptyResponse check the method ReadListOfClusters
func TestReadListOfClustersEmptyResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(func(responseWriter http.ResponseWriter, request *http.Request) {
		// just check the URL, don't send any body in the response
		checkURL(t, request, "/api/v1/client/cluster")
	})
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadListOfClusters()
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadListOfClustersWrongJSON check the method ReadListOfClusters
func TestReadListOfClustersWrongJSON(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerImpl(t, ReadClustersURL, ImproperJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadListOfClusters()
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadListOfClustersResponseError check the method ReadListOfClusters
func TestReadListOfClustersResponseError(t *testing.T) {
	// start a local HTTP server
	server := httptest.NewServer(http.NotFoundHandler())
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadListOfClusters()
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadListOfTriggersEmptyList check the method ReadListOfTriggers
func TestReadListOfTriggersEmptyList(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerImpl(t, ReadTriggersURL, StatusOKJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	triggers, err := api.ReadListOfTriggers()
	if err != nil {
		t.Fatal(err)
	}

	if len(triggers) != 0 {
		t.Fatal("Expected empty list of triggers")
	}
}

// TestReadListOfTriggersOneTrigger check the method ReadListOfTriggers
func TestReadListOfTriggersOneTrigger(t *testing.T) {
	// start a local HTTP server
	const responseAsString = `
	{
		"triggers": [{"ID":0,"Name":"Name","Type":"must-gather"}],
		"status":"ok"
	}`
	server := mockedHTTPServer(standardHandlerImpl(t, ReadTriggersURL, responseAsString))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	triggers, err := api.ReadListOfTriggers()
	if err != nil {
		t.Fatal(err)
	}

	if len(triggers) != 1 {
		t.Fatal("Expected list with one trigger only")
	}
}

// TestReadListOfTriggersErrorStatus check the method ReadListOfTriggers
func TestReadListOfTriggersErrorStatus(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerImpl(t, ReadTriggersURL, StatusErrorJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadListOfTriggers()
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadListOfTriggersEmptyResponse check the method ReadListOfTriggers
func TestReadListOfTriggersEmptyResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(func(responseWriter http.ResponseWriter, request *http.Request) {
		// just check the URL, don't send any body in the response
		checkURL(t, request, "/api/v1/client/trigger")
	})
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadListOfTriggers()
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadListOfTriggersWrongJSON check the method ReadListOfTriggers
func TestReadListOfTriggersWrongJSON(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerImpl(t, ReadTriggersURL, ImproperJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadListOfTriggers()
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadListOfTriggersResponseError check the method ReadListOfTriggers
func TestReadListOfTriggersResponseError(t *testing.T) {
	// start a local HTTP server
	server := httptest.NewServer(http.NotFoundHandler())
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadListOfTriggers()
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadListOfConfigurationProfilesEmptyList check the method ReadListOfConfigurationProfiles
func TestReadListOfConfigurationProfilesEmptyList(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerImpl(t, ReadProfilesURL, StatusOKJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	profiles, err := api.ReadListOfConfigurationProfiles()
	if err != nil {
		t.Fatal(err)
	}

	if len(profiles) != 0 {
		t.Fatal("Expected empty list of profiles")
	}
}

// TestReadListOfConfigurationProfilesOneProfile check the method ReadListOfConfigurationProfiles
func TestReadListOfConfigurationProfilesOneProfile(t *testing.T) {
	// start a local HTTP server
	const responseAsString = `
	{
		"profiles": [{}],
		"status":"ok"
	}`
	server := mockedHTTPServer(standardHandlerImpl(t, ReadProfilesURL, responseAsString))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	profiles, err := api.ReadListOfConfigurationProfiles()
	if err != nil {
		t.Fatal(err)
	}

	if len(profiles) != 1 {
		t.Fatal("Expected list with one profile only")
	}
}

// TestReadListOfConfigurationProfilesErrorStatus check the method ReadListOfConfigurationProfiles
func TestReadListOfConfigurationProfilesErrorStatus(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerImpl(t, ReadProfilesURL, StatusErrorJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadListOfConfigurationProfiles()
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadListOfConfigurationProfilesEmptyResponse check the method ReadListOfConfigurationProfiles
func TestReadListOfConfigurationProfilesEmptyResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(func(responseWriter http.ResponseWriter, request *http.Request) {
		// just check the URL, don't send any body in the response
		checkURL(t, request, "/api/v1/client/profile")
	})
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadListOfConfigurationProfiles()
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadListOfConfigurationProfilesWrongJSON check the method ReadListOfConfigurationProfiles
func TestReadListOfConfigurationProfilesWrongJSON(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerImpl(t, ReadProfilesURL, ImproperJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadListOfConfigurationProfiles()
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadListOfConfigurationProfilesResponseError check the method ReadListOfConfigurationProfiles
func TestReadListOfConfigurationProfilesResponseError(t *testing.T) {
	// start a local HTTP server
	server := httptest.NewServer(http.NotFoundHandler())
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadListOfConfigurationProfiles()
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadListOfConfigurationsEmptyList check the method ReadListOfConfigurations
func TestReadListOfConfigurationsEmptyList(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerImpl(t, ReadConfigurationsURL, StatusOKJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	configurations, err := api.ReadListOfConfigurations()
	if err != nil {
		t.Fatal(err)
	}

	if len(configurations) != 0 {
		t.Fatal("Expected empty list of configurations")
	}
}

// TestReadListOfConfigurationsOneConfiguration check the method ReadListOfConfigurations
func TestReadListOfConfigurationsOneConfiguration(t *testing.T) {
	// start a local HTTP server
	const responseAsString = `
	{
		"configuration": [{}],
		"status":"ok"
	}`
	server := mockedHTTPServer(standardHandlerImpl(t, ReadConfigurationsURL, responseAsString))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	configurations, err := api.ReadListOfConfigurations()
	if err != nil {
		t.Fatal(err)
	}

	if len(configurations) != 1 {
		t.Fatal("Expected list with one configuration only")
	}
}

// TestReadListOfConfigurationsErrorStatus check the method ReadListOfConfigurations
func TestReadListOfConfigurationsErrorStatus(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerImpl(t, ReadConfigurationsURL, StatusErrorJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadListOfConfigurations()
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadListOfConfigurationsEmptyResponse check the method ReadListOfConfigurations
func TestReadListOfConfigurationsEmptyResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(func(responseWriter http.ResponseWriter, request *http.Request) {
		// just check the URL, don't send any body in the response
		checkURL(t, request, ReadConfigurationsURL)
	})
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadListOfConfigurations()
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadListOfConfigurationsWrongJSON check the method ReadListOfConfigurations
func TestReadListOfConfigurationsWrongJSON(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerImpl(t, ReadConfigurationsURL, ImproperJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadListOfConfigurations()
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadListOfConfigurationsResponseError check the method ReadListOfConfigurations
func TestReadListOfConfigurationsResponseError(t *testing.T) {
	// start a local HTTP server
	server := httptest.NewServer(http.NotFoundHandler())
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadListOfConfigurations()
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadTriggerByIDStandardResponse check the method ReadTriggerByID
func TestReadTriggerByIDStandardResponse(t *testing.T) {
	const responseAsString = `
	{
		"trigger": {"ID":1,"Type":"must-gather","Cluster":"ffff"},
		"status":"ok"
	}`
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerImpl(t, "/api/v1/client/trigger/1", responseAsString))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	trigger, err := api.ReadTriggerByID("1")
	if err != nil {
		t.Fatal(err)
	}
	expected := types.Trigger{
		ID:      1,
		Type:    "must-gather",
		Cluster: "ffff",
	}
	if *trigger != expected {
		t.Fatal("Improper trigger returned: ", trigger)
	}
}

// TestReadTriggerByIDImproperJSON check the method ReadTriggerByID
func TestReadTriggerByIDImproperJSON(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerImpl(t, "/api/v1/client/trigger/1", ImproperJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadTriggerByID("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadTriggerByIDErrorResponse check the method ReadTriggerByID
func TestReadTriggerByIDErrorResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerImpl(t, "/api/v1/client/trigger/1", StatusErrorJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadTriggerByID("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadTriggerByIDEmptyResponse check the method ReadTriggerByID
func TestReadTriggerByIDEmptyResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(func(responseWriter http.ResponseWriter, request *http.Request) {
		// just check the URL, don't send any body in the response
		checkURL(t, request, "/api/v1/client/trigger/1")
	})
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadTriggerByID("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadTriggerByIDNotFoundResponse check the method ReadTriggerByID
func TestReadTriggerByIDNotFoundResponse(t *testing.T) {
	// start a local HTTP server
	server := httptest.NewServer(http.NotFoundHandler())
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadTriggerByID("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadConfigurationProfileStandardResponse check the method ReadConfigurationProfile
func TestReadConfigurationProfileStandardResponse(t *testing.T) {
	const responseAsString = `
	{
		"profile": {"id":1,"configuration":"","changed_at":"2020-01-01","changed_by":"tester","description":"description"},
		"status":"ok"
	}`
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerImpl(t, "/api/v1/client/profile/1", responseAsString))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	profile, err := api.ReadConfigurationProfile("1")
	if err != nil {
		t.Fatal(err)
	}
	expected := types.ConfigurationProfile{
		ID:            1,
		Configuration: "",
		ChangedAt:     "2020-01-01",
		ChangedBy:     "tester",
		Description:   "description",
	}
	if *profile != expected {
		t.Fatal("Improper configuration profile returned: ", *profile)
	}
}

// TestReadConfigurationProfileImproperJSON check the method ReadConfigurationProfile
func TestReadConfigurationProfileImproperJSON(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerImpl(t, "/api/v1/client/profile/1", ImproperJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadConfigurationProfile("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadConfigurationProfileErrorResponse check the method ReadConfigurationProfile
func TestReadConfigurationProfileErrorResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerImpl(t, "/api/v1/client/profile/1", StatusErrorJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadConfigurationProfile("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadConfigurationProfileEmptyResponse check the method ReadConfigurationProfile
func TestReadConfigurationProfileEmptyResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(func(responseWriter http.ResponseWriter, request *http.Request) {
		// just check the URL, don't send any body in the response
		checkURL(t, request, "/api/v1/client/profile/1")
	})
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadConfigurationProfile("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadConfigurationProfileNotFoundResponse check the method ReadConfigurationProfile
func TestReadConfigurationProfileNotFoundResponse(t *testing.T) {
	// start a local HTTP server
	server := httptest.NewServer(http.NotFoundHandler())
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadConfigurationProfile("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadClusterConfigurationByIDStandardResponse check the method ReadClusterConfigurationByID
func TestReadClusterConfigurationByIDStandardResponse(t *testing.T) {
	const responseAsString = `
	{
		"configuration": "config",
		"status":"ok"
	}`
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerImpl(t, "/api/v1/client/configuration/1", responseAsString))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	configuration, err := api.ReadClusterConfigurationByID("1")
	if err != nil {
		t.Fatal(err)
	}
	expected := "config"
	if *configuration != expected {
		t.Fatal("Improper cluster configuration returned: ", *configuration)
	}
}

// TestReadClusterConfigurationByIDImproperJSON check the method ReadClusterConfigurationByID
func TestReadClusterConfigurationByIDImproperJSON(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerImpl(t, "/api/v1/client/configuration/1", ImproperJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadClusterConfigurationByID("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadClusterConfigurationByIDErrorResponse check the method ReadClusterConfigurationByID
func TestReadClusterConfigurationByIDErrorResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerImpl(t, "/api/v1/client/configuration/1", StatusErrorJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadClusterConfigurationByID("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadClusterConfigurationByIDEmptyResponse check the method ReadClusterConfigurationByID
func TestReadClusterConfigurationByIDEmptyResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(func(responseWriter http.ResponseWriter, request *http.Request) {
		// just check the URL, don't send any body in the response
		checkURL(t, request, "/api/v1/client/configuration/1")
	})
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadClusterConfigurationByID("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestReadClusterConfigurationByIDNotFoundResponse check the method ReadClusterConfigurationByID
func TestReadClusterConfigurationByIDNotFoundResponse(t *testing.T) {
	// start a local HTTP server
	server := httptest.NewServer(http.NotFoundHandler())
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	_, err := api.ReadClusterConfigurationByID("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestEnableClusterConfigurationStandardResponse check the method EnableClusterConfiguration
func TestEnableClusterConfigurationStandardResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/configuration/1/enable", "PUT", StatusOKJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.EnableClusterConfiguration("1")
	if err != nil {
		t.Fatal(err)
	}
}

// TestEnableClusterConfigurationImproperJSON check the method EnableClusterConfiguration
func TestEnableClusterConfigurationImproperJSON(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/configuration/1/enable", "PUT", ImproperJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.EnableClusterConfiguration("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestEnableClusterConfigurationErrorResponse check the method EnableClusterConfiguration
func TestEnableClusterConfigurationErrorResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/configuration/1/enable", "PUT", StatusErrorJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.EnableClusterConfiguration("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestEnableClusterConfigurationNotFoundResponse check the method DeleteClusterConfiguration
func TestEnableClusterConfigurationNotFoundResponse(t *testing.T) {
	// start a local HTTP server
	server := httptest.NewServer(http.NotFoundHandler())
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.EnableClusterConfiguration("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestDisableClusterConfigurationStandardResponse check the method DisableClusterConfiguration
func TestDisableClusterConfigurationStandardResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/configuration/1/disable", "PUT", StatusOKJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.DisableClusterConfiguration("1")
	if err != nil {
		t.Fatal(err)
	}
}

// TestDisableClusterConfigurationImproperJSON check the method DisableClusterConfiguration
func TestDisableClusterConfigurationImproperJSON(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/configuration/1/disable", "PUT", ImproperJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.DisableClusterConfiguration("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestDisableClusterConfigurationErrorResponse check the method DisableClusterConfiguration
func TestDisableClusterConfigurationErrorResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/configuration/1/disable", "PUT", StatusErrorJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.DisableClusterConfiguration("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestDisableClusterConfigurationNotFoundResponse check the method DeleteClusterConfiguration
func TestDisableClusterConfigurationNotFoundResponse(t *testing.T) {
	// start a local HTTP server
	server := httptest.NewServer(http.NotFoundHandler())
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.DisableClusterConfiguration("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestDeleteClusterConfigurationStandardResponse check the method DeleteClusterConfiguration
func TestDeleteClusterConfigurationStandardResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/configuration/1", "DELETE", StatusOKJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.DeleteClusterConfiguration("1")
	if err != nil {
		t.Fatal(err)
	}
}

// TestDeleteClusterConfigurationImproperJSON check the method DeleteClusterConfiguration
func TestDeleteClusterConfigurationImproperJSON(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/configuration/1", "DELETE", ImproperJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.DeleteClusterConfiguration("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestDeleteClusterConfigurationErrorResponse check the method DeleteClusterConfiguration
func TestDeleteClusterConfigurationErrorResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/configuration/1", "DELETE", StatusErrorJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.DeleteClusterConfiguration("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestDeleteClusterConfigurationNotFoundResponse check the method DeleteClusterConfiguration
func TestDeleteClusterConfigurationNotFoundResponse(t *testing.T) {
	// start a local HTTP server
	server := httptest.NewServer(http.NotFoundHandler())
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.DeleteClusterConfiguration("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestDeleteClusterStandardResponse check the method DeleteCluster
func TestDeleteClusterStandardResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/cluster/1", "DELETE", StatusOKJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.DeleteCluster("1")
	if err != nil {
		t.Fatal(err)
	}
}

// TestDeleteClusterImproperJSON check the method DeleteCluster
func TestDeleteClusterImproperJSON(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/cluster/1", "DELETE", ImproperJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.DeleteCluster("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestDeleteClusterErrorResponse check the method DeleteCluster
func TestDeleteClusterErrorResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/cluster/1", "DELETE", StatusErrorJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.DeleteCluster("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestDeleteClusterNotFoundResponse check the method DeleteCluster
func TestDeleteClusterNotFoundResponse(t *testing.T) {
	// start a local HTTP server
	server := httptest.NewServer(http.NotFoundHandler())
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.DeleteCluster("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestDeleteConfigurationProfileStandardResponse check the method DeleteConfigurationProfile
func TestDeleteConfigurationProfileStandardResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/profile/1", "DELETE", StatusOKJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.DeleteConfigurationProfile("1")
	if err != nil {
		t.Fatal(err)
	}
}

// TestDeleteConfigurationProfileImproperJSON check the method DeleteConfigurationProfile
func TestDeleteConfigurationProfileImproperJSON(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/profile/1", "DELETE", ImproperJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.DeleteConfigurationProfile("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestDeleteConfigurationProfileErrorResponse check the method DeleteConfigurationProfile
func TestDeleteConfigurationProfileErrorResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/profile/1", "DELETE", StatusErrorJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.DeleteConfigurationProfile("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestDeleteConfigurationProfileNotFoundResponse check the method DeleteConfigurationProfile
func TestDeleteConfigurationProfileNotFoundResponse(t *testing.T) {
	// start a local HTTP server
	server := httptest.NewServer(http.NotFoundHandler())
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.DeleteConfigurationProfile("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestActivateTriggerStandardResponse check the method ActivateTrigger
func TestActivateTriggerStandardResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/trigger/1/activate", "PUT", StatusOKJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.ActivateTrigger("1")
	if err != nil {
		t.Fatal(err)
	}
}

// TestActivateTriggerImproperJSON check the method ActivateTrigger
func TestActivateTriggerImproperJSON(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/trigger/1/activate", "PUT", ImproperJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.ActivateTrigger("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestActivateTriggerErrorResponse check the method ActivateTrigger
func TestActivateTriggerErrorResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/trigger/1/activate", "PUT", StatusErrorJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.ActivateTrigger("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestActivateTriggerNotFoundResponse check the method ActivateTrigger
func TestActivateTriggerNotFoundResponse(t *testing.T) {
	// start a local HTTP server
	server := httptest.NewServer(http.NotFoundHandler())
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.ActivateTrigger("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestDeactivateTriggerStandardResponse check the method DeactivateTrigger
func TestDeactivateTriggerStandardResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/trigger/1/deactivate", "PUT", StatusOKJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.DeactivateTrigger("1")
	if err != nil {
		t.Fatal(err)
	}
}

// TestDeactivateTriggerImproperJSON check the method DeactivateTrigger
func TestDeactivateTriggerImproperJSON(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/trigger/1/deactivate", "PUT", ImproperJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.DeactivateTrigger("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestDeactivateTriggerErrorResponse check the method DeactivateTrigger
func TestDeactivateTriggerErrorResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/trigger/1/deactivate", "PUT", StatusErrorJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.DeactivateTrigger("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestDeactivateTriggerNotFoundResponse check the method DeactivateTrigger
func TestDeactivateTriggerNotFoundResponse(t *testing.T) {
	// start a local HTTP server
	server := httptest.NewServer(http.NotFoundHandler())
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.DeactivateTrigger("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestDeleteTriggerStandardResponse check the method DeleteTrigger
func TestDeleteTriggerStandardResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/trigger/1", "DELETE", StatusOKJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.DeleteTrigger("1")
	if err != nil {
		t.Fatal(err)
	}
}

// TestDeleteTriggerImproperJSON check the method DeleteTrigger
func TestDeleteTriggerImproperJSON(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/trigger/1", "DELETE", ImproperJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.DeleteTrigger("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestDeleteTriggerErrorResponse check the method DeleteTrigger
func TestDeleteTriggerErrorResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/trigger/1", "DELETE", StatusErrorJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.DeleteTrigger("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

// TestDeleteTriggerNotFoundResponse check the method DeleteTrigger
func TestDeleteTriggerNotFoundResponse(t *testing.T) {
	// start a local HTTP server
	server := httptest.NewServer(http.NotFoundHandler())
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.DeleteTrigger("1")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

func TestAddClusterStandardResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/cluster/cluster2", "POST", StatusOKJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.AddCluster("cluster2")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddClusterErrorResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/cluster/cluster2", "POST", StatusErrorJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.AddCluster("cluster2")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

func TestAddClusterImproperJSONResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/cluster/cluster2", "POST", ImproperJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.AddCluster("cluster2")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}

func TestAddClusterNotFoundResponse(t *testing.T) {
	// start a local HTTP server
	server := httptest.NewServer(http.NotFoundHandler())
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.AddCluster("cluster2")
	if err == nil {
		t.Fatal("Error is expected to be returned")
	}
}
