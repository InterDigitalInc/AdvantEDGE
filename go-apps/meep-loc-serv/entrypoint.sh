#!/bin/bash
set -e

echo "MEEP_HOST_URL: ${MEEP_HOST_URL}"
echo "MEEP_SANDBOX_NAME: ${MEEP_SANDBOX_NAME}"

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

# Set basepath for API files
for file in /api/*.yaml; do
    if [[ ! -e "$file" ]]; then continue; fi
    setBasepath $file
done

# Set basepath for user-supplied API files
for file in /user-api/*.yaml; do
    if [[ ! -e "$file" ]]; then continue; fi
    setBasepath $file
done

# Start service
exec /meep-loc-serv
