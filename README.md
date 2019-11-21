# Command-line client for managing Insights operator

[![Go Report Card](https://goreportcard.com/badge/github.com/RedHatInsights/insights-operator-cli)](https://goreportcard.com/report/github.com/RedHatInsights/insights-operator-cli)

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

Please make sure to run `./test.sh` to check all changes made in the source code.

## Testing

Unit tests can be started by the following command:

```
./test.sh
```

It is also possible to specify CLI options for Go test. For example, if you need to disable test results caching, use the following command:

```
./test -count=1
```
