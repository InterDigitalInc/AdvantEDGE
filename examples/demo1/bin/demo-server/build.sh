#!/bin/sh

# Get full path to script directory
SCRIPT=$(readlink -f "$0")
BASEDIR=$(dirname "$SCRIPT")

# Set GO env variables
GOOS=linux

# Create vendor folder
cd $BASEDIR/../../src/demo-server
go mod vendor

# Build demo App server
echo "Building Demo Service Go Server"
cd $BASEDIR/../../src/demo-server
go build -o $BASEDIR/demo-server .
