#!/bin/bash

#set -vx

function usage() {
    echo "Usage: $0 {--rebuild_dai}" >&2
    echo "       --rebuild_dai To force MEEP-DAI micor-service to be rebuilt"
    exit 1
}

# Get full path to script directory
SCRIPT=$(readlink -f "$0")
BASEDIR=$(dirname "$SCRIPT")

# Get Arguments
REBUILD_DAI=$1
if [ "$REBUILD_DAI" != "" ]; then
    if [ "$REBUILD_DAI" != "--rebuild_dai" ]; then
        usage
    fi
fi

echo ""
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ">>> onboarded-demo"
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ""
DEMOBIN=$BASEDIR/bin/onboarded-demo

$BASEDIR/src/onboarded-demo/build.sh $DEMOBIN

# Update meep-dai binary folder with the MEC application
MEEPDAI_ONBOARDEDAPP_PATH=../../bin/meep-dai/onboardedapp
echo ">>> Updating DAI micro-service with onboarded MEC applications"
if [ -d $MEEPDAI_ONBOARDEDAPP_PATH ]; then
    rm -fr $MEEPDAI_ONBOARDEDAPP_PATH
fi
mkdir -p $MEEPDAI_ONBOARDEDAPP_PATH
cp -Rp $DEMOBIN $MEEPDAI_ONBOARDEDAPP_PATH

if [ "$REBUILD_DAI" != "" ]; then
    echo ">>> Building DAI micro-service"
    meepctl build meep-dai #--nolint
    echo ">>> Dockerizing DAI micro-service with new onboarded MEC applications"
    meepctl dockerize meep-dai
fi

echo ""
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ">>> Building Demo 4 UE application Go Server"
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ""
DEMOBIN=$BASEDIR/bin/demo-server

$BASEDIR/src/demo-server/backend/build.sh $DEMOBIN

echo ""
echo ">>> Demo4-ue Service build completed"
