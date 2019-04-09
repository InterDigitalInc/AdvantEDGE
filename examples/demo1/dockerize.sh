#!/bin/bash

cd bin/demo-server
docker build --no-cache --rm -t demo-server .
cd ../../
