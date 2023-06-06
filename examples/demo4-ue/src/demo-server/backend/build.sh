#!/bin/sh

#set -vx

# Get full path to script directory
SCRIPT=$(readlink -f "$0")
BASEDIR=$(dirname "$SCRIPT")

# Create build directory
if [ "$#" -ne 1 ]; then
    echo "Missing bin directory"
    exit
fi
BINDIR=$1
mkdir -p $BINDIR $BINDIR/user-api

# Build demo App server
cd $BASEDIR

#find . -name "*.go" -type f -exec golangci-lint run {} \;
go build -o $BINDIR/demo-server .

cp -Rp ./api/  $BINDIR
mv $BINDIR/api/swagger.yaml $BINDIR/api/MEC\ Demo\ 4\ API
cp $BINDIR/api/* $BINDIR/user-api