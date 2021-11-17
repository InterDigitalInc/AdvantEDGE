#!/bin/bash
set -e

{   
    printf "mode: advantedge \n"
    printf "sandbox: \n"
    printf "mecplatform: ${MEEP_MEP_NAME} \n"
    printf "appid:  \n"
    printf "localurl: ${MEEP_POD_NAME}  \n"
    printf "port: "
} <~/>app_instance.yaml

# Start service
exec /demo-server ./app_instance.yaml
