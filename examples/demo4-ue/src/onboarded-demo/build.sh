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
mkdir -p $BINDIR

# Build demo App server
cd $BASEDIR

go build -o $BINDIR/onboarded-demo4 .
