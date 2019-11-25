go test $(go list ./... | grep -v tests) $@
exit $?
