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

echo ""
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ">> Configuring Object Store secrets"
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ""
kubectl create secret generic meep-thanos-objstore-config --from-file=objstore.yml=$BASEDIR/thanos.yaml
kubectl create secret generic meep-thanos-archive-objstore-config --from-file=objstore.yml=$BASEDIR/thanos-archive.yaml

echo ""
echo ">>> Object Store configuration completed"




