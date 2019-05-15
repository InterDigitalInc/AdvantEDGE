#!/bin/sh

# Get full path to script directory
SCRIPT=$(readlink -f "$0")
BASEDIR=$(dirname "$SCRIPT")

# Install node module dependencies
cd $BASEDIR
npm ci

