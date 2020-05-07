#!/bin/bash
set -e

echo "MEEP_SANDBOX_NAME: ${MEEP_SANDBOX_NAME}"
echo "MEEP_HOST_URL: ${MEEP_HOST_URL}"

# Prepend sandbox name to REST API yaml files
for file in /swagger/*-api.yaml; do
    echo "Prepending [${MEEP_SANDBOX_NAME}] to basepath in: $file"
    sed -i 's,basePath: \"/\?,basePath: \"/'${MEEP_SANDBOX_NAME}'/,' $file;
done

# Update spec links in index.html
# sed -i 's,"url": "\([^"]*\)","url": "'${MEEP_HOST_URL}'/'${MEEP_SANDBOX_NAME}'/api/\1",g' /swagger/index.html

# Start virt engine
exec /meep-sandbox-ctrl