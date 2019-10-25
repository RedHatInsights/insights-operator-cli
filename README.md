# Command-line client for managing Insights operator

## Description

A simple CLI client for managing the Insights operator. Currently this client supports just basic operations to retrieve and change configuration of operator on selected cluster.

## Supported commands

* `exit`
* `quit`
* `bye`
  * exit from the client

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

Please make sure to run `make test` to check all changes made in the source code.

## Testing

Unit tests can be started by the following command:

    ./test.sh

It is also possible to specify CLI options for Go test. For example, if you need to disable test results caching, use the following command:

    ./test -count=1

