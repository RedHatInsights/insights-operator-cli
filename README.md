# Command-line client for managing Insights operator

[![Go Report Card](https://goreportcard.com/badge/github.com/RedHatInsights/insights-operator-cli)](https://goreportcard.com/report/github.com/RedHatInsights/insights-operator-cli) [![Build Status](https://travis-ci.org/RedHatInsights/insights-operator-cli.svg?branch=master)](https://travis-ci.org/RedHatInsights/insights-operator-cli) [![codecov](https://codecov.io/gh/RedHatInsights/insights-operator-cli/branch/master/graph/badge.svg)](https://codecov.io/gh/RedHatInsights/insights-operator-cli)

## Description

A simple CLI client for managing the Insights operator. Currently this client supports just basic operations to retrieve and change configuration of operator on selected cluster.

## Supported commands

### Cluster operations:
* **list clusters** list all clusters known to the service
* **delete cluster ##**         delete selected cluster
* **add cluster**               create new cluster
* **new cluster**               alias for previous command

### Configuration profiles:
* **list profiles**             list all profiles known to the service
* **describe profile ##**       describe profile selected by its ID
* **delete profile ##**         delete profile selected by its ID

### Cluster configurations:
* **list configurations**       list all configurations known to the service
* **describe configuration ##** describe cluster configuration selected by its ID
* **add configuration**         add new configuration
* **new configuration**         alias for previous command
* **enable configuration ##**   enable cluster configuration selected by its ID
* **disable configuration ##**  disable cluster configuration selected by its ID
* **delete configuration ##**   delete configuration selected by its ID

### Must-gather trigger:       
* **list triggers**             list all triggers
* **describe trigger ##**       describe trigger selected by its ID
* **add trigger**               add new trigger
* **new trigger**               alias for previous command
* **activate trigger ##**       activate trigger selected by its ID
* **deactivate trigger ##**     deactivate trigger selected by its ID
* **delete trigger**            delete trigger

### Other commands:
* **version**                   print version information
* **quit**                      quit the application
* **exit**                      dtto
* **bye**                       dtto
* **help**                      this help
* **copyright**                 displays copyright notice
* **license**                   displays license used by this project
* **authors**                   displays list of authors


## How to build the CLI client

Use the standard Go command:

```
go build
```

This command should create an executable file named `insights-operator-cli`.

## Start

Just run the executable file created by `go build`:

```
./insights-operator-cli
```

## Configuration

Configuration are stored in a file `config.toml`.
At this moment, just `CONTROLLER_URL` needs to be specified.

## Contributing

Please look into document [CONTRIBUTING.md](CONTRIBUTING.md) that contains all information about how to contribute to this project.

Please look also at [Definitiot of Done](DoD.md) document with further informations.

Also make sure to run `./test.sh` to check all changes made in the source code.

## Testing

Unit tests can be started by the following command:

```
./test.sh
```

It is also possible to specify CLI options for Go test. For example, if you need to disable test results caching, use the following command:

```
./test -count=1
```

## CI

[Travis CI](https://travis-ci.com/) is configured for this repository. Several tests and checks are started for all pull requests:

* Unit tests that use the standard tool `go test`
* `go fmt` tool to check code formatting. That tool is run with `-s` flag to perform [following transformations](https://golang.org/cmd/gofmt/#hdr-The_simplify_command)
* `go vet` to report likely mistakes in source code, for example suspicious constructs, such as Printf calls whose arguments do not align with the format string.
* `golint` as a linter for all Go sources stored in this repository
* `gocyclo` to report all functions and methods with too high cyclomatic complexity. The cyclomatic complexity of a function is calculated according to the following rules: 1 is the base complexity of a function +1 for each 'if', 'for', 'case', '&&' or '||' Go Report Card warns on functions with cyclomatic complexity > 9
* `goconst` to find repeated strings that could be replaced by a constant
* `ineffassign` to detect and print all ineffectual assignments in Go code
* `errcheck` for checking for all unchecked errors in go programs
* `shellcheck` to perform static analysis for all shell scripts used in this repository
* `abcgo` to measure ABC metrics for Go source code and check if the metrics does not exceed specified threshold

History of checks done by CI is available at [RedHatInsights / insights-operator-cli](https://travis-ci.org/RedHatInsights/insights-operator-cli).

## Contribution

Please look into document [CONTRIBUTING.md](CONTRIBUTING.md) that contains all information about how to contribute to this project.
