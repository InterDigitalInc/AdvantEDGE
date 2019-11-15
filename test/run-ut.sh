#!/bin/bash

# Get full path to script directory
SCRIPT=$(readlink -f "$0")
BASEDIR=$(dirname "$SCRIPT")

GOPKGS=(
    meep-model
    meep-net-char-mgr
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
    go test ./...
    echo ""
done
