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

package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	expect "github.com/Netflix/go-expect"
)

// get the console for already running application
func getConsole() *expect.Console {
	console, err := expect.NewConsole(expect.WithStdout(os.Stdout))
	if err != nil {
		log.Fatal(err)
	}
	return console
}

// start CLI application w/o color output and w/o command-line completer
func startCLI(console *expect.Console) *exec.Cmd {
	command := exec.Command("./insights-operator-cli", "--colors=false", "--completer=false")
	command.Stdin = console.Tty()
	command.Stdout = console.Tty()
	command.Stderr = console.Tty()

	err := command.Start()
	if err != nil {
		log.Fatal(err)
	}

	return command
}

// pause test for a given amount of time
func pause() {
	time.Sleep(time.Second)
}

// check whether the 'quit' command works as expected
func checkQuitCommand() {
	console := getConsole()
	defer console.Close()

	command := startCLI(console)
	pause()
	console.Send("quit\n")
	pause()

	go func() {
		console.ExpectEOF()
	}()

	err := command.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

// check whether the 'version' command works as expected
func checkVersionCommand() {
	console := getConsole()
	defer console.Close()

	command := startCLI(console)
	pause()
	console.Send("version\n")
	pause()
	console.ExpectString("Insights operator CLI client")
	console.ExpectString("version")
	console.ExpectString("compiled")
	console.Send("quit\n")
	pause()

	go func() {
		console.ExpectEOF()
	}()

	err := command.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	checkQuitCommand()
	checkVersionCommand()
}
