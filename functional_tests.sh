go build -race

if [ $? -eq 0 ]
then
    echo "CLI client build ok"
else
    echo "Build failed"
    exit 1
fi

go test -c -o ./functional-tests tests/functional_test.go

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
