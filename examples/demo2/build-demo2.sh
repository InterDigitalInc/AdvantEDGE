#!/bin/bash

# Get full path to script directory
SCRIPT=$(readlink -f "$0")
BASEDIR=$(dirname "$SCRIPT")
cd $BASEDIR

# Retrieve workdir
WORKDIR=`grep workdir $HOME/.meepctl.yaml | sed 's/^.*workdir:[ \t]*//'`
DEMO2DIR=${WORKDIR}/virt-engine/user-charts/demo2

# Copy demo charts to workdir
mkdir -p ${DEMO2DIR}
cp -r ./charts ${DEMO2DIR}
cp -r ./values ${DEMO2DIR}

echo ">>> Demo2 build completed"
