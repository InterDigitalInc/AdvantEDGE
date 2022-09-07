#!/bin/bash

# Get full path to script directory
SCRIPT=$(readlink -f "$0")
BASEDIR=$(dirname "$SCRIPT")

GOAPPS=(
    meep-loc-serv/server
    meep-rnis/server
    meep-wais/server
    meep-ams/server
    #meep-vis/server
)

GOPKGS=(
    meep-couch
    meep-gis-asset-mgr
    meep-vis-traffic-mgr
    meep-metrics
    meep-model
    meep-mq
    meep-net-char-mgr
    meep-pdu-session-store
    meep-subscriptions
    meep-users
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
