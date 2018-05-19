#!/usr/bin/env bash
set -e -o pipefail

. build-conf.sh

echo "$VERSION_FILE ---- version file"

versionData=`grep version package.json | sed  's/[,]//g'`
# versionData='{ "version":"0.12.1" }';

echo "versionData='{ $versionData }';" > spo/current-spo.json

echo "versionData='{ $versionData }';" > ../src/gui/static/dist/current-spo.json

echo "versionData='{ $versionData }';" > ../src/gui/static/src/current-spo.json
