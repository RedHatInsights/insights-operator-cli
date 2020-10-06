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

package restapi

import (
	"encoding/json"
	"fmt"
	"github.com/redhatinsighs/insights-operator-cli/types"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	communicationErrorWithServerErrorMessage = "Communication error with the server %v"
	unableToReadResponseBodyError            = "Unable to read response body"
)

// performReadRequest function try to perform HTTP request using the HTTP GET
// method and if the call is successful read the body of response.
func performReadRequest(url string) ([]byte, error) {
	// disable "G107 (CWE-88): Potential HTTP request made with variable url"
	// #nosec G107
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf(communicationErrorWithServerErrorMessage, err)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Expected HTTP status 200 OK, got %d", response.StatusCode)
	}
	body, readErr := ioutil.ReadAll(response.Body)
	defer closeResponseBody(response)

	if readErr != nil {
		return nil, fmt.Errorf(unableToReadResponseBodyError)
	}

	return body, nil
}

// performWriteRequest function try to perform HTTP request using the specified
// HTTP method (POST, PUT, DELETE) and if the call is successful read the body
// of response.
func performWriteRequest(url string, method string, payload io.Reader) error {
	var client http.Client

	request, err := http.NewRequest(method, url, payload)
	if err != nil {
		return fmt.Errorf("Error creating request %v", err)
	}

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf(communicationErrorWithServerErrorMessage, err)
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Expected HTTP status 200 OK, 201 Created or 202 Accepted, got %d", response.StatusCode)
	}
	body, readErr := ioutil.ReadAll(response.Body)
	defer closeResponseBody(response)

	if readErr != nil {
		return fmt.Errorf(unableToReadResponseBodyError)
	}
	return parseResponse(body)
}

// closeResponseBody function tries to close body of HTTP response with basic
// error check
func closeResponseBody(response *http.Response) {
	err := response.Body.Close()
	if err != nil {
		log.Println(err)
	}
}

// parseResponse function tries to parse the body of HTTP response into JSON
// structure that should contain at least one attribute stored under key
// "Status".
func parseResponse(body []byte) error {
	resp := types.Response{}
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return err
	}
	if resp.Status != "ok" {
		return fmt.Errorf("Error response: %s", resp.Status)
	}
	return nil
}
