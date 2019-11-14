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
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func performReadRequest(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Communication error with the server %v", err)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Expected HTTP status 200 OK, got %d", response.StatusCode)
	}
	body, readErr := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if readErr != nil {
		return nil, fmt.Errorf("Unable to read response body")
	}

	return body, nil
}

func performWriteRequest(url string, method string, payload io.Reader) error {
	var client http.Client

	request, err := http.NewRequest(method, url, payload)
	if err != nil {
		return fmt.Errorf("Error creating request %v", err)
	}

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("Communication error with the server %v", err)
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Expected HTTP status 200 OK, 201 Created or 202 Accepted, got %d", response.StatusCode)
	}
	return nil
}
