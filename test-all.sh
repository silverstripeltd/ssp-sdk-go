#!/bin/bash

PACKAGES=$(go list ./...)
EXIT_CODE=0

function setExit {
	local previousExitCode=$?
	if [ "$previousExitCode" -ne "0" ]; then
    	EXIT_CODE=${previousExitCode}
  	fi
}

# we need to hack around the coverage report to include several packes so we manually craft the coverage report
# and append each package coverage information to the overage.out file
echo "mode: atomic" > coverage.out
for PKG in $PACKAGES; do
	go fmt ${PKG}; setExit
	go vet ${PKG}; setExit
	go test -v -race -coverprofile=./profile.out $PKG; setExit
	if [ -f profile.out ]; then
    	tail -n +2 profile.out >> coverage.out; rm profile.out
  	fi
done

exit ${EXIT_CODE}
