#!/usr/bin/env bash

# TODO: find a better way to test this scenario. see the test file for context

### grep exits w/ a non-zero exit code if it doesn't find anything.
{ go test -v ./internal/failing || true ; } | grep "at:" | grep "internal/failing/verify_any__awesomeness__fils_calls_test.go" > /dev/null
