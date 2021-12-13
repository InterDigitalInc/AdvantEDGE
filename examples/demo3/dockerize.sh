#!/bin/bash

# Get full path to script directory
SCRIPT=$(readlink -f "$0")
BASEDIR=$(dirname "$SCRIPT")

DEMOBIN=$BASEDIR/bin/demo-server

echo ""
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ">>> Dockerizing Demo3 Server"
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ""

# Copy frontend to demo-server folder
rm -r $DEMOBIN/static
mkdir -p $DEMOBIN/static

# Copy frontend to static folder
echo ">>> Copying Demo Server"
cp -Rf $BASEDIR/bin/demo-frontend/* $DEMOBIN/static

# Copy Dockerfile & config to bin folder
cp $BASEDIR/src/backend/Dockerfile $DEMOBIN
# cp $BASEDIR/src/backend/app_instance.yaml $DEMOBIN
cp $BASEDIR/entrypoint.sh $DEMOBIN

echo ">>> Dockerizing"
cd $DEMOBIN
docker build --no-cache --rm -t meep-docker-registry:30001/demo3 .
docker push meep-docker-registry:30001/demo3
cd $BASEDIR

echo ""
echo ">>> Done"
