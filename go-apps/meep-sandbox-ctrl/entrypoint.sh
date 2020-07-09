#!/bin/bash
set -e

echo "MEEP_SANDBOX_NAME: ${MEEP_SANDBOX_NAME}"
echo "MEEP_HOST_URL: ${MEEP_HOST_URL}"
echo "USER_SWAGGER: ${USER_SWAGGER}"
echo "USER_SWAGGER_SANDBOX: ${USER_SWAGGER_SANDBOX}"

# Prepend sandbox name to REST API yaml files
for file in /swagger/*-api.yaml; do
    echo "Prepending [${MEEP_SANDBOX_NAME}] to basepath in: $file"
    sed -i 's,basePath: \"/\?,basePath: \"/'${MEEP_SANDBOX_NAME}'/,' $file;
done

# Copy user-swagger & adapt basepath to sandbox
if [[ ! -z "${USER_SWAGGER}" ]]; then
    cp -r ${USER_SWAGGER} ${USER_SWAGGER_SANDBOX}
    for file in ${USER_SWAGGER_SANDBOX}/*-api.yaml; do
        echo "Prepending [${MEEP_SANDBOX_NAME}] to basepath in: $file"
        sed -i 's,basePath: \"/\?,basePath: \"/'${MEEP_SANDBOX_NAME}'/,' $file;
    done
fi

# Update spec links in index.html
# sed -i 's,"url": "\([^"]*\)","url": "'${MEEP_HOST_URL}'/'${MEEP_SANDBOX_NAME}'/api/\1",g' /swagger/index.html

# Start virt engine
exec /meep-sandbox-ctrl
