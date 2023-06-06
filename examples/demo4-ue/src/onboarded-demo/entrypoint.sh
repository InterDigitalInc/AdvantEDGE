#!/bin/bash
set -e

echo "mode: advantedge" >app_instance.yaml
echo "sandbox:" >>app_instance.yaml
echo "mecplatform: ${MEEP_MEP_NAME}" >>app_instance.yaml
echo "appid:" ${MEEP_APP_ID} >>app_instance.yaml
echo "localurl: ${MEEP_POD_NAME}" >>app_instance.yaml
echo "port:" >>app_instance.yaml

# Start service
exec /onboarded-demo4 ./app_instance.yaml
