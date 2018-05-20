#!/usr/bin/env bash

# Tests a single module
# Use it like this:
#   cd src/testme
#   SPOPATH=/path/to/sporepo ../../scripts/test-coverage.sh
# Coverage will open up in html
# You don't need to do this if you are not symlinking the repo into $GOPATH
# I recommend you put SPOPATH in ~/.bashrc
# Example:
# If spo repo is located at
#   /home/user/repos/spo
# Then $SPOPATH will be
#   /home/user/repos

MODULE="$1"
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

if [ ! -z "$MODULE" ]; then
    pushd "$DIR" >/dev/null
    pushd "../src/${MODULE}" > /dev/null
fi

go test -v -cover -coverprofile=coverage.out
if [ $? -eq 0 ]; then
    sed -i 's|_'${SPOPATH}'|github.com/spo|g' coverage.out
    go tool cover -html=coverage.out
    rm coverage.out
fi

if [ ! -z "$MODULE" ]; then
    popd >/dev/null
    popd >/dev/null
fi
