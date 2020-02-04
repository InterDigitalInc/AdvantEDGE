#!/bin/bash

# Get full path to script directory
SCRIPT=$(readlink -f "$0")
BASEDIR=$(dirname "$SCRIPT")

echo ""
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ">>> Installing redis DB for Unit Testing"
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ""

helm install --name meep-ut-redis --set master.service.nodePort=30380 $BASEDIR/../charts/redis/

echo ""
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ">>> Installing couch DB for Unit Testing"
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ""

helm install --name meep-ut-couchdb --set service.nodePort=30985 --set persistentVolume.enabled=false --set persistentVolumeClaim.enabled=false $BASEDIR/../charts/couchdb/

echo ""
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ">>> Installing influx DB for Unit Testing"
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ""

helm install --name meep-ut-influxdb --set service.apiNodePort=30986 --set service.rpcNodePort=30988 --set persistence.enabled=false $BASEDIR/../charts/influxdb/