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

package restapi_test

import (
	"bytes"
	"github.com/redhatinsighs/insights-operator-cli/restapi"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// mockedHttpServer prepares new instance of testing HTTP server
func mockedHttpServer(handler func(responseWriter http.ResponseWriter, request *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}

// TestPerformReadRequestProperResponse check if body of response can be processed correctly by performReadRequest function.
func TestPerformReadRequestProperResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHttpServer(func(responseWriter http.ResponseWriter, request *http.Request) {
		// send response to be tested later
		responseWriter.Write([]byte(`OK`))
	})
	// close the server when test finishes
	defer server.Close()

	// try to read response from the server
	body, err := restapi.PerformReadRequest(server.URL)
	if err != nil {
		t.Fatal("Unable to perform request", err)
	}

	// check for proper body
	if !bytes.Equal(body, []byte(`OK`)) {
		t.Fatal("Improper response body", err)
	}
}

// TestPerformReadRequestStatusCode check how response can be processed by performReadRequest function.
func TestPerformReadRequestStatusCode(t *testing.T) {
	// start a local HTTP server
	server := mockedHttpServer(func(responseWriter http.ResponseWriter, request *http.Request) {
		// send response to be tested later
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte(`error`))
	})
	// close the server when test finishes
	defer server.Close()

	// try to read response from the server
	body, err := restapi.PerformReadRequest(server.URL)
	if err == nil {
		t.Fatal("Error is expected")
	}

	// body needs to be nil in case of any error
	if body != nil {
		t.Fatal("Body needs to be nil in case of any error")
	}
}

// TestPerformReadRequestEmptyBody check how response can be processed by performReadRequest function.
func TestPerformReadRequestEmptyBody(t *testing.T) {
	// start a local HTTP server
	server := mockedHttpServer(func(responseWriter http.ResponseWriter, request *http.Request) {
		// send response to be tested later
		responseWriter.WriteHeader(http.StatusOK)
	})
	// close the server when test finishes
	defer server.Close()

	// try to read response from the server
	body, err := restapi.PerformReadRequest(server.URL)
	if err != nil {
		t.Fatal("Unable to perform request", err)
	}

	// check for proper body - it needs to be empty this time
	if len(body) > 0 {
		t.Fatal("Improper response body", err)
	}
}

// TestPerformReadRequestErrorInCommunication check how response can be processed by performReadRequest function.
func TestPerformReadRequestErrorInCommunication(t *testing.T) {
	// try to read response from the server, but by using improper URL
	body, err := restapi.PerformReadRequest("")
	if err == nil {
		t.Fatal("Error is expected")
	}

	// body needs to be nil in case of any error
	if body != nil {
		t.Fatal("Body needs to be nil in case of any error")
	}
}

// TestPerformWriteRequestProperResponse check the behaviour of function performWriteRequest
func TestPerformWriteRequestProperResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHttpServer(func(responseWriter http.ResponseWriter, request *http.Request) {
		// send response to be tested later
		responseWriter.Write([]byte(`{"status":"ok"}`))
	})
	// close the server when test finishes
	defer server.Close()

	// try to read response from the server
	err := restapi.PerformWriteRequest(server.URL, "POST", nil)
	if err != nil {
		t.Fatal("Unable to perform request", err)
	}
}

// TestPerformWriteRequestErrorResponse check the behaviour of function performWriteRequest
func TestPerformWriteRequestErrorResponse(t *testing.T) {
	// start a local HTTP server
	server := mockedHttpServer(func(responseWriter http.ResponseWriter, request *http.Request) {
		// send response to be tested later
		responseWriter.Write([]byte(`{"status":"error"}`))
	})
	// close the server when test finishes
	defer server.Close()

	// try to read response from the server
	err := restapi.PerformWriteRequest(server.URL, "POST", nil)
	if err == nil {
		t.Fatal("Error should be reported")
	}
}

// TestPerformWriteRequestImproperStatusCode check the behaviour of function performWriteRequest
func TestPerformWriteRequestImproperStatusCode(t *testing.T) {
	// start a local HTTP server
	server := mockedHttpServer(func(responseWriter http.ResponseWriter, request *http.Request) {
		// send response to be tested later
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte(`{"status":"ok"}`))
	})
	// close the server when test finishes
	defer server.Close()

	// try to read response from the server
	err := restapi.PerformWriteRequest(server.URL, "POST", nil)
	if err == nil {
		t.Fatal("Error should be reported")
	}
	if !strings.HasPrefix(err.Error(), "Expected HTTP status 200 OK, 201 Created or 202 Accepted") {
		t.Fatal("Unexpected error message:", err.Error())
	}
}

// TestPerformWriteRequestEmptyBody check the behaviour of function performWriteRequest
func TestPerformWriteRequesttEmptyBody(t *testing.T) {
	// start a local HTTP server
	server := mockedHttpServer(func(responseWriter http.ResponseWriter, request *http.Request) {
		// send empty response to be tested later
	})
	// close the server when test finishes
	defer server.Close()

	// try to read response from the server
	err := restapi.PerformWriteRequest(server.URL, "POST", nil)
	if err == nil {
		t.Fatal("Error should be reported")
	}
}

// TestPerformWriteRequestErrorInCommunication check how response can be processed by performWriteRequest function.
func TestPerformWriteRequestErrorInCommunication(t *testing.T) {
	err := restapi.PerformWriteRequest("", "POST", nil)
	if err == nil {
		t.Fatal("Error is expected")
	}
	if !strings.HasPrefix(err.Error(), "Communication error with the server") {
		t.Fatal("Unexpected error message:", err.Error())
	}
}

// TestParseResponseOkStatus check the behaviour of function parseResponse
func TestParseResponseOkStatus(t *testing.T) {
	body := []byte(`{"status":"ok"}`)
	err := restapi.ParseResponse(body)
	if err != nil {
		t.Fatal(err)
	}
}

// TestParseResponseErrorStatus check the behaviour of function parseResponse
func TestParseResponseErrorStatus(t *testing.T) {
	body := []byte(`{"status":"error"}`)
	err := restapi.ParseResponse(body)
	if err == nil {
		t.Fatal("Error report is expected")
	}
	if err.Error() != "Error response: error" {
		t.Fatal("Invalid error response:", err)
	}
}

// TestParseResponseImproperJSON check the behaviour of function parseResponse
func TestParseResponseImproperJSON(t *testing.T) {
	body := []byte(`this is not JSON`)
	err := restapi.ParseResponse(body)
	if err == nil {
		t.Fatal("Error report is expected")
	}
}
