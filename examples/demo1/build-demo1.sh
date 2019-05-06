#!/bin/bash

# Get full path to script directory
SCRIPT=$(readlink -f "$0")
BASEDIR=$(dirname "$SCRIPT")

# Build Iperf Proxy
$BASEDIR/bin/iperf-proxy/build.sh

#Build frontend for demo svc app
$BASEDIR/src/demo-frontend/build.sh
#Build client for demo svc app
$BASEDIR/src/demo-client/js/build.sh

rm -r $BASEDIR/bin/demo-server/static
mkdir $BASEDIR/bin/demo-server/static
cp -Rf $BASEDIR/src/demo-frontend/dist/* $BASEDIR/bin/demo-server/static

# Build Demo App Server
echo "Building Demo Service REST API Go Server"
$BASEDIR/bin/demo-server/build.sh

# Build docker image
echo "Building Demo Service Docker image"
./dockerize.sh

echo "Demo Service build completed"




