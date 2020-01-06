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

package commands_test

import (
	"bytes"
	"github.com/logrusorgru/aurora"
	"github.com/redhatinsighs/insights-operator-cli/commands"
	"io"
	"os"
	"sync"
)

// configureColorizer configures the Aurora colorizer. For tests
// it is preferred to turn-off colorization.
func configureColorizer() {
	colorizer := aurora.NewAurora(false)
	commands.SetColorizer(colorizer)
}

// captureStandardOutput captures the standard output for specified code
// block and then returns the captured output.
// Please see https://medium.com/@hau12a1/golang-capturing-log-println-and-fmt-println-output-770209c791b4
// for further explanation how it works under the hood.
//
// TODO: put it into a dependent module
func captureStandardOutput(function func()) (string, error) {
	// backup of the real stdout
	stdout := os.Stdout

	// temporary replacement for stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		return "", err
	}

	// temporarily replace real Stdout by the mocked one
	defer func() {
		os.Stdout = stdout
	}()
	os.Stdout = writer

	// channel with captured standard output
	captured := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, reader)
		captured <- buf.String()
	}()
	wg.Wait()
	// provided function that (probably) prints something to standard output
	function()
	writer.Close()
	return <-captured, nil
}
