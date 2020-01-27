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

// mockedHttpServer prepares new instance of testing HTTP server
func mockedHttpServer(handler func(responseWriter http.ResponseWriter, request *http.Request)) *httptest.Server {
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
