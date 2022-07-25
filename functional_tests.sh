#!/bin/bash
# Copyright 2020 Red Hat, Inc
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


go build -race

# shellcheck disable=SC2181
if [ $? -eq 0 ]
then
    echo "CLI client build ok"
else
    echo "Build failed"
    exit 1
fi

go test -c -o ./functional-tests tests/*.go

# shellcheck disable=SC2181
if [ $? -eq 0 ]
then
    echo "Functional tests build ok"
else
    echo "Build failed"
    exit 1
fi

./functional-tests -test.v
EXIT_VALUE=$?
exit $EXIT_VALUE
