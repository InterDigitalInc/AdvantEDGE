#!/bin/bash

# Get full path to script directory
SCRIPT=$(readlink -f "$0")
BASEDIR=$(dirname "$SCRIPT")

GOAPPS=(
    meep-loc-serv/server
    meep-rnis/server
    meep-wais/server
)

GOPKGS=(
    meep-couch
    meep-metric-store
    meep-model
    meep-mq
    meep-net-char-mgr
    meep-postgis
    meep-watchdog
)

echo ""
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ">>> Running Unit Tests"
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ""

for pkg in "${GOPKGS[@]}" ; do
    echo "+ pkg: $pkg"
    cd $BASEDIR/../go-packages/$pkg
    go test -count=1 ./... -cover
    echo ""
done

for app in "${GOAPPS[@]}" ; do
    echo "+ app: $app"
    cd $BASEDIR/../go-apps/$app
    go test -count=1 ./... -cover
    echo ""
done
