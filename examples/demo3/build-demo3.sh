#!/bin/bash

# Get full path to script directory
SCRIPT=$(readlink -f "$0")
BASEDIR=$(dirname "$SCRIPT")

echo ""
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ">>> Building Demo Service Frontend"
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ""
$BASEDIR/src/client/build.sh
$BASEDIR/src/frontend/build.sh

DEMOBIN=$BASEDIR/bin/demo-server

echo ""
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ">>> Building Demo 3 Go Server"
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ""
$BASEDIR/src/backend/build.sh $DEMOBIN

echo ""
echo ">>> Demo Service build completed"
