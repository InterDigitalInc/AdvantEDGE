#!/bin/bash

# Get full path to script directory
SCRIPT=$(readlink -f "$0")
BASEDIR=$(dirname "$SCRIPT")

GOPKGS=(
    meep-model
    meep-net-char-mgr
    meep-metric-store
    meep-watchdog
    meep-couch
)

echo ""
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ">>> Running Unit Tests"
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ""

for pkg in "${GOPKGS[@]}" ; do
    echo "+ pkg: $pkg"
    cd $BASEDIR/../go-packages/$pkg
    go test -count=1 ./...
    echo ""
done
