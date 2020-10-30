#!/bin/bash
set -e

echo "MEEP_SANDBOX_NAME: ${MEEP_SANDBOX_NAME}"
echo "USER_SWAGGER: ${USER_SWAGGER}"
echo "USER_SWAGGER_SANDBOX: ${USER_SWAGGER_SANDBOX}"

# Prepend sandbox name to REST API yaml files
for file in /swagger/*-api.yaml; do
    echo "Prepending [${MEEP_SANDBOX_NAME}] to basepath in: $file"
    sed -i 's,basePath: \"/\?,basePath: \"/'${MEEP_SANDBOX_NAME}'/,' $file;
    echo "Replacing {apiRoot} with ${MEEP_SANDBOX_NAME} to url in: $file"
    sed -i 's/{apiRoot}/'${MEEP_SANDBOX_NAME}'/g' $file;
done

# Copy user-swagger & adapt basepath to sandbox
if [[ ! -z "${USER_SWAGGER}" ]]; then
    cp -r ${USER_SWAGGER} ${USER_SWAGGER_SANDBOX}
    shopt -s nullglob
    for file in ${USER_SWAGGER_SANDBOX}/*-api.yaml; do
        echo "Prepending [${MEEP_SANDBOX_NAME}] to basepath in: $file"
        sed -i 's,basePath: \"/\?,basePath: \"/'${MEEP_SANDBOX_NAME}'/,' $file;
        echo "Replacing {apiRoot} with ${MEEP_SANDBOX_NAME} to url in: $file"
        newValue=${MEEP_HOST_URL}/${MEEP_SANDBOX_NAME}
        httpPrefixToRemove='http://'
        httpsPrefixToRemove='https://'
        echo $newValue
        newValueSuffix="${newValue/$httpsPrefixToRemove/}"
        newValueSuffix="${newValueSuffix/$httpPrefixToRemove/}"
        echo $newValueSuffix
        sed -i "s@{apiRoot}@$newValueSuffix@g" $file;
    done
fi

# Start virt engine
exec /meep-sandbox-ctrl
