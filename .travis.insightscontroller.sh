#!/bin/bash

git clone https://github.com/RedHatInsights/insights-operator-controller.git
cd insights-operator-controller
go build
export ENV=test
./local_storage/create_database_sqlite.sh

set -m
./insights-operator-controller &bg
sleep 5
