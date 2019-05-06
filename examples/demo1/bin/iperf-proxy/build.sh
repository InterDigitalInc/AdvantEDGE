#!/bin/sh

# Get full path to script directory
SCRIPT=$(readlink -f "$0")
BASEDIR=$(dirname "$SCRIPT")

# Set GO env variables
GOOS=linux

# Create vendor folder
cd $BASEDIR/../../src/iperf-proxy
go mod vendor

# Build demo App iperf proxy server
echo "Building Iperf Proxy Go Server"
cd $BASEDIR/../../src/iperf-proxy
go build -o $BASEDIR/iperf-proxy .
