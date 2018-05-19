#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
echo "spo binary dir:" "$DIR"
pushd "$DIR" >/dev/null
go run cmd/spo/spo.go --gui-dir="${DIR}/src/gui/static/" $@
popd >/dev/null
