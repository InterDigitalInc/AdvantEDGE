#!/bin/sh

# Get full path to script directory
SCRIPT=$(readlink -f "$0")
BASEDIR=$(dirname "$SCRIPT")

# Install node module dependencies
echo "Building Demo Service JS REST Client"
cd $BASEDIR
npm ci

