#!/bin/bash

# Get full path to script directory
SCRIPT=$(readlink -f "$0")
BASEDIR=$(dirname "$SCRIPT")

# Configure environment
GOOS=linux
IMAGE_NAME=meepctl
echo "$IMAGE_NAME"

cd $BASEDIR

# Clean build
echo "...clean"
go clean

# Create vendor folder
echo "...vendor"
go mod vendor

# Lint code
echo "...lint"
golangci-lint run

# Build
echo "...build"
go build -o ./$IMAGE_NAME .

# Install
echo "...install"
go install
