#!/bin/bash

# Get full path to script directory
SCRIPT=$(readlink -f "$0")
BASEDIR=$(dirname "$SCRIPT")

echo ""
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ">> Removing old secrets"
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ""
kubectl delete secret meep-thanos-objstore-config
kubectl delete secret meep-thanos-archive-objstore-config
kubectl delete secret meep-influx-objstore-config

echo ""
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ">> Configuring Object Store secrets"
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ""
kubectl create secret generic meep-thanos-objstore-config --from-file=objstore.yml=$BASEDIR/objstore-thanos.yaml
kubectl create secret generic meep-thanos-archive-objstore-config --from-file=objstore.yml=$BASEDIR/objstore-thanos-archive.yaml
kubectl create secret generic meep-influx-objstore-config --from-file=credentials=$BASEDIR/objstore-influx.cfg

echo ""
echo ">>> Object Store configuration completed"




