go build

if [ $? -eq 0 ]
then
    echo "CLI client build ok"
else
    echo "Build failed"
    exit 1
fi

go build -o functional-tests tests/functional_tests.go

if [ $? -eq 0 ]
then
    echo "Functional tests build ok"
else
    echo "Build failed"
    exit 1
fi

./functional-tests
EXIT_VALUE=$?
exit $EXIT_VALUE
