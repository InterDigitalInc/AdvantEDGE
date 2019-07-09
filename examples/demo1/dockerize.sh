#!/bin/bash

# Get full path to script directory
SCRIPT=$(readlink -f "$0")
BASEDIR=$(dirname "$SCRIPT")

DEMOBIN=$BASEDIR/bin/demo-server

echo ""
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ">>> Dockerizing Demo Server"
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ""

# Copy frontend to bin folder
rm -r $DEMOBIN/static
mkdir -p $DEMOBIN/static
cp -Rf $BASEDIR/bin/demo-frontend/* $DEMOBIN/static

# Copy Dockerfile to bin folder
cp $BASEDIR/src/demo-server/Dockerfile $DEMOBIN

# Dockerize demo 
cd $DEMOBIN
docker build --no-cache --rm -t meep-docker-registry:30001/demo-server .
docker push meep-docker-registry:30001/demo-server

echo ""
echo ">>> Demo Service dockerize completed"




