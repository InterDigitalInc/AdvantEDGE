#!/bin/bash
set -e

echo "MEEP_HOST_URL: ${MEEP_HOST_URL}"
echo "MEEP_SANDBOX_NAME: ${MEEP_SANDBOX_NAME}"
echo "USER_SWAGGER: ${USER_SWAGGER}"
echo "USER_SWAGGER_SANDBOX: ${USER_SWAGGER_SANDBOX}"

# Update API yaml basepaths to enable "Try-it-out" feature
# OAS2: Set relative path to sandbox name + endpoint path (origin will be derived from browser URL)
# OAS3: Set full path to provided Host URL + sandbox name + endpoint path
setBasepath() {
    # OAS3
    hostName=$(echo "${MEEP_HOST_URL}" | sed -E 's/^\s*.*:\/\///g')
    #newHostName=${hostName}/${MEEP_SANDBOX_NAME}
    echo "Replacing [localhost] with ${hostName} to url in: $1"
    sed -i "s,localhost,${hostName},g" $1;

    # OAS2 and OAS3
    echo "Replacing [sandboxname] with ${MEEP_SANDBOX_NAME} to basepath or url in: $1"
    #sed -i 's,basePath: \"/\?,basePath: \"/'${MEEP_SANDBOX_NAME}'/,' $1;
    sed -i "s,sandboxname,${MEEP_SANDBOX_NAME},g" $1;
}

# Set baspath for AdvantEDGE Swagger API files
for file in /swagger/*-api.yaml; do
    setBasepath $file
done

# Set baspath for User-provided Swagger API files
if [[ ! -z "${USER_SWAGGER}" ]]; then
    cp -r ${USER_SWAGGER} ${USER_SWAGGER_SANDBOX}
    shopt -s nullglob
    for file in ${USER_SWAGGER_SANDBOX}/*-api.yaml; do
        setBasepath $file
    done
fi

# Start virt engine
exec /meep-sandbox-ctrl
