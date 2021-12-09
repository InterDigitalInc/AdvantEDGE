#!/bin/sh
SCRIPT=$(readlink -f "$0")
BASEDIR=$(dirname "$SCRIPT")

# Make sure package-lock.json does not contain a bad version of 'warning' package
cd $BASEDIR
npm ci
