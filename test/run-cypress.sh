#!/bin/bash

# Get full path to script directory
SCRIPT=$(readlink -f "$0")
BASEDIR=$(dirname "$SCRIPT")
cd $BASEDIR

echo ""
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ">>> Running Cypress Tests"
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo ""

# Retrieve workdir
WORKDIR=`grep workdir $HOME/.meepctl.yaml | sed 's/^.*workdir:[ \t]*//'`
DEMO2DIR=${WORKDIR}/charts/demo2

# Copy demo charts to workdir
mkdir -p ${DEMO2DIR}
cp -r ../examples/demo2/charts ${DEMO2DIR}
cp -r ../examples/demo2/values ${DEMO2DIR}

# Install Cypress
npm ci

# Run cypress tests
npm run cy:run
