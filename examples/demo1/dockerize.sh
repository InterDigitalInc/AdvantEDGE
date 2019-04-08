#!/bin/bash

cd bin/demo-server
docker build --no-cache --rm -t demo1-server .
cd ../../
