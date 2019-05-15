#!/bin/sh

# Get full path to script directory
SCRIPT=$(readlink -f "$0")
BASEDIR=$(dirname "$SCRIPT")

# Create build directory
if [ "$#" -ne 1 ]; then
    echo "Missing bin directory"
    exit
fi
BINDIR=$1
mkdir -p $BINDIR

# Set GO env variables
GOOS=linux

# Create vendor folder
cd $BASEDIR
go mod vendor

# Build demo App server
cd $BASEDIR
go build -o $BINDIR/demo-server .
