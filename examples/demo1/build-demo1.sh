#!/bin/bash

# Get full path to script directory
SCRIPT=$(readlink -f "$0")
BASEDIR=$(dirname "$SCRIPT")

DEMOBIN=$BASEDIR/bin/demo-server
PROXYBIN=$BASEDIR/bin/iperf-proxy

echo ""
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ">> Building Iperf Proxy JS Client"
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ""
$BASEDIR/src/iperf-proxy-client/js/build.sh

echo ""
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ">>> Building Iperf Proxy Go Server"
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ""
$BASEDIR/src/iperf-proxy/build.sh $PROXYBIN

echo ""
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ">>> Building Demo Service JS Client"
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ""
$BASEDIR/src/demo-client/js/build.sh

echo ""
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ">>> Building Demo Service Frontend"
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ""
$BASEDIR/src/demo-frontend/build.sh

echo ""
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ">>> Building Demo Service Go Server"
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ""
$BASEDIR/src/demo-server/build.sh $DEMOBIN

echo ""
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ">>> Dockerizing Demo Server"
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ""

# Copy frontend to bin folder
rm -r $DEMOBIN/static
mkdir -p $DEMOBIN/static
cp -Rf $BASEDIR/src/demo-frontend/dist/* $DEMOBIN/static

# Copy Dockerfile to bin folder
cp $BASEDIR/src/demo-server/Dockerfile $DEMOBIN

# Dockerize demo 
cd $DEMOBIN
docker build --no-cache --rm -t demo-server .

echo ""
echo ">>> Demo Service build completed"




