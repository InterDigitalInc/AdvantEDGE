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

# Copy frontend to static folder
echo ">>> Copying Demo Server"
cp -Rf $BASEDIR/bin/demo-frontend/* $DEMOBIN/static
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
cp $BASEDIR/src/backend/Dockerfile $DEMOBIN
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"

# Dockerize demo 
# cd $DEMOBIN
# docker build --no-cache --rm -t meep-docker-registry:30001/demo-server3 .
# docker push meep-docker-registry:30001/demo-server3

echo ""
echo ">>> Demo Service dockerize completed"




