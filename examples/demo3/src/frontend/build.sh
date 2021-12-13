#!/bin/sh

# Get full path to script directory
SCRIPT=$(readlink -f "$0")
BASEDIR=$(dirname "$SCRIPT")

# Make sure package-lock.json does not contain a bad version of 'warning' package
cd $BASEDIR
if grep -q '\"warning\": \"^3.' "package-lock.json"; then
    echo "ERROR: package-lock.json erroneously includes dependency on warning v3.x"
    exit
fi

# Build Web UI Distribution using Webpack config
cd $BASEDIR
npm ci
rm -r $BASEDIR/dist
npm run build -- --output-path=$BASEDIR/../../bin/demo-frontend
