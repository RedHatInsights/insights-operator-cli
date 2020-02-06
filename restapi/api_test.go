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

// expectNoErrors checks if the error is not reported by REST API call
func expectNoErrors(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

// expectError checks if the error is reported by REST API call
func expectError(t *testing.T, err error) {
	if err == nil {
		t.Fatal("Error is expected to be returned")
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
	expectNoErrors(t, err)

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
	expectNoErrors(t, err)

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
	expectError(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectNoErrors(t, err)

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
	expectNoErrors(t, err)

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
	expectError(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectNoErrors(t, err)

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
	expectNoErrors(t, err)

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
	expectError(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectNoErrors(t, err)

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
	expectNoErrors(t, err)

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
	expectError(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectNoErrors(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectNoErrors(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectNoErrors(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectNoErrors(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectNoErrors(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectNoErrors(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectNoErrors(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectNoErrors(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectNoErrors(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectNoErrors(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectNoErrors(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
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
	expectError(t, err)
}

func TestAddClusterStandardResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/cluster/cluster2", "POST", StatusOKJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.AddCluster("cluster2")
	expectNoErrors(t, err)
}

func TestAddClusterErrorResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/cluster/cluster2", "POST", StatusErrorJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.AddCluster("cluster2")
	expectError(t, err)
}

func TestAddClusterImproperJSONResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, "/api/v1/client/cluster/cluster2", "POST", ImproperJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.AddCluster("cluster2")
	expectError(t, err)
}

func TestAddClusterNotFoundResponse(t *testing.T) {
	// start a local HTTP server
	server := httptest.NewServer(http.NotFoundHandler())
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.AddCluster("cluster2")
	expectError(t, err)
}

func TestAddTriggerStandardResponse(t *testing.T) {
	// start a local HTTP server
	URL := "/api/v1/client/cluster/cluster2/trigger/must-gather?username=name&reason=reason&link=link"
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, URL, "POST", StatusOKJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.AddTrigger("name", "cluster2", "reason", "link")
	expectNoErrors(t, err)
}

func TestAddTriggerErrorResponse(t *testing.T) {
	// start a local HTTP server
	URL := "/api/v1/client/cluster/cluster2/trigger/must-gather?username=name&reason=reason&link=link"
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, URL, "POST", StatusErrorJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.AddTrigger("name", "cluster2", "reason", "link")
	expectError(t, err)
}

func TestAddTriggerImproperJSONResponse(t *testing.T) {
	// start a local HTTP server
	URL := "/api/v1/client/cluster/cluster2/trigger/must-gather?username=name&reason=reason&link=link"
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, URL, "POST", ImproperJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.AddTrigger("name", "cluster2", "reason", "link")
	expectError(t, err)
}

func TestAddTriggerImproperNotFoundResponse(t *testing.T) {
	// start a local HTTP server
	server := httptest.NewServer(http.NotFoundHandler())
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.AddTrigger("name", "cluster2", "reason", "link")
	expectError(t, err)
}

func TestAddConfigurationProfileStandardResponse(t *testing.T) {
	// start a local HTTP server
	URL := "/api/v1/client/profile?username=name&description=description"
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, URL, "POST", StatusOKJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.AddConfigurationProfile("name", "description", []byte{1, 2, 3})
	expectNoErrors(t, err)
}

func TestAddConfigurationProfileErrorResponse(t *testing.T) {
	// start a local HTTP server
	URL := "/api/v1/client/profile?username=name&description=description"
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, URL, "POST", StatusErrorJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.AddConfigurationProfile("name", "description", []byte{1, 2, 3})
	expectError(t, err)
}

func TestAddConfigurationProfileImproperJSONResponse(t *testing.T) {
	// start a local HTTP server
	URL := "/api/v1/client/profile?username=name&description=description"
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, URL, "POST", ImproperJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.AddConfigurationProfile("name", "description", []byte{1, 2, 3})
	expectError(t, err)
}

func TestAddConfigurationProfileNotFoundResponse(t *testing.T) {
	// start a local HTTP server
	server := httptest.NewServer(http.NotFoundHandler())
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.AddConfigurationProfile("name", "description", []byte{1, 2, 3})
	expectError(t, err)
}

func TestAddClusterConfigurationStandardResponse(t *testing.T) {
	// start a local HTTP server
	URL := "/api/v1/client/cluster/cluster2/configuration/create?username=name&reason=reason&description=description"
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, URL, "POST", StatusOKJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.AddClusterConfiguration("name", "cluster2", "reason", "description", []byte{1, 2, 3})
	expectNoErrors(t, err)
}

func TestAddClusterConfigurationErrorResponse(t *testing.T) {
	// start a local HTTP server
	URL := "/api/v1/client/cluster/cluster2/configuration/create?username=name&reason=reason&description=description"
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, URL, "POST", StatusErrorJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.AddClusterConfiguration("name", "cluster2", "reason", "description", []byte{1, 2, 3})
	expectError(t, err)
}

func TestAddClusterConfigurationImproperJSONResponse(t *testing.T) {
	// start a local HTTP server
	URL := "/api/v1/client/cluster/cluster2/configuration/create?username=name&reason=reason&description=description"
	server := mockedHTTPServer(standardHandlerForMethodImpl(t, URL, "POST", ImproperJSON))
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.AddClusterConfiguration("name", "cluster2", "reason", "description", []byte{1, 2, 3})
	expectError(t, err)
}

func TestAddClusterConfigurationNotFoundResponse(t *testing.T) {
	// start a local HTTP server
	server := httptest.NewServer(http.NotFoundHandler())
	// close the server when test finishes
	defer server.Close()

	api := restapi.NewRestAPI(server.URL)

	// perform REST API call against mocked HTTP server
	err := api.AddClusterConfiguration("name", "cluster2", "reason", "description", []byte{1, 2, 3})
	expectError(t, err)
}
