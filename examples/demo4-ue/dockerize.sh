#!/bin/bash

#set -vx

# Get full path to script directory
SCRIPT=$(readlink -f "$0")
BASEDIR=$(dirname "$SCRIPT")

echo ""
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ">>> Dockerizing onboarded-demo4 Server"
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ""

DEMOBIN=$BASEDIR/bin/onboarded-demo

# Copy Dockerfile & config to bin folder
cp $BASEDIR/src/onboarded-demo/Dockerfile $DEMOBIN
cp $BASEDIR/src/onboarded-demo/entrypoint.sh $DEMOBIN
# cp $BASEDIR/src/onboarded-demo/onboarded-demo.yaml $DEMOBIN

echo ">>> Dockerizing"
cd $DEMOBIN
docker build --no-cache --rm -t meep-docker-registry:30001/onboarded-demo4 .
docker push meep-docker-registry:30001/onboarded-demo4
cd $BASEDIR

echo ""
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ">>> Dockerizing demo4-ue Server"
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ""

DEMOBIN=$BASEDIR/bin/demo-server

# Copy frontend to static folder
#echo ">>> Copying Demo Server"
#cp -Rf $BASEDIR/bin/demo-frontend/* $DEMOBIN/static

# Copy Dockerfile & config to bin folder
cp $BASEDIR/src/demo-server/backend/Dockerfile $DEMOBIN
# cp $BASEDIR/src/demo-server/backend/app_instance.yaml $DEMOBIN
cp $BASEDIR/src/demo-server/entrypoint.sh $DEMOBIN

echo ">>> Dockerizing"
cd $DEMOBIN
docker build --no-cache --rm -t meep-docker-registry:30001/demo4-ue .
docker push meep-docker-registry:30001/demo4-ue
cd $BASEDIR

echo ""
echo ">>> Done"
