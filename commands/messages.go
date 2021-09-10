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

package commands

// Generated documentation is available at:
// https://pkg.go.dev/github.com/RedHatInsights/insights-operator-cli/commands
//
// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-cli/packages/commands/messages.html

// various messages to be displayed to user via command-line interface
const (
	operationCancelled = "Cancelled"
	changedAt          = "Changed at"
	changedBy          = "Changed by"

	// user is not logged in
	notLoggedIn = "Not logged in"

	// any object (cluster, configuration, profile) has been deleted
	deleted = "deleted"

	// prompt displayed for any TAB-completable inputs (file selection
	// etc.)
	configurationFilePrompt = "configuration file (TAB to complete): "

	// (any) condition has been met
	conditionSet = "yes"

	// cluster UUID
	clusterUUID = "Cluster"
)
