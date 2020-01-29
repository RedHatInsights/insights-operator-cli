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

	"testing"
)

const (
	RESTApiPrefix   = "/api/v1/client/"
	ReadClustersURL = RESTApiPrefix + "cluster"
	ReadTriggersURL = RESTApiPrefix + "trigger"
	ReadProfilesURL = RESTApiPrefix + "profile"

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

// writeBody writes a given text into the response that is to be send to receiver
func writeBody(responseWriter http.ResponseWriter, body string) {
	responseWriter.Write([]byte(body))
}

// standardHandlerImpl is an implementation of handler that checks URL and when it's expected send a response
// with body that contains a body filled with given response string
func standardHandlerImpl(t *testing.T, expectedURL, responseStr string) func(responseWriter http.ResponseWriter, request *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
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
	triggers, err := api.ReadListOfConfigurationProfiles()
	if err != nil {
		t.Fatal(err)
	}

	if len(triggers) != 0 {
		t.Fatal("Expected empty list of triggers")
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
	triggers, err := api.ReadListOfConfigurationProfiles()
	if err != nil {
		t.Fatal(err)
	}

	if len(triggers) != 1 {
		t.Fatal("Expected list with one trigger only")
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
