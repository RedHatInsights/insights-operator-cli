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

package types

// Generated documentation is available at:
// https://godoc.org/github.com/RedHatInsights/insights-operator-cli/types
//
// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-cli/packages/types/cluster.html

// Cluster structure represents cluster record in the controller service.
//     ID: unique key
//     Name: cluster GUID in the following format:
//         c8590f31-e97e-4b85-b506-c45ce1911a12
type Cluster struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ClustersResponse structure represents response of controller service to
// cluster request.
//     Status: status of response
//     Clusters: list of clusters
type ClustersResponse struct {
	Status   string    `json:"status"`
	Clusters []Cluster `json:"clusters"`
}
