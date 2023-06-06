#!/bin/bash
set -e

echo "MEEP_HOST_URL: ${MEEP_HOST_URL}"
echo "MEEP_SANDBOX_NAME: ${MEEP_SANDBOX_NAME}"
echo "MEEP_MEP_NAME: ${MEEP_MEP_NAME}"
echo "MEEP_CODECOV: ${MEEP_CODECOV}"
echo "MEEP_BROKER: ${MEEP_BROKER}" # E.g. mqtt://test.mosquito.org:1338 or amqp://guest:guest@localhost:5672
echo "MEEP_POA_LIST: ${MEEP_POA_LIST}" # E.g. poa-5g1;poa-5g2

if [[ ! -z "${MEEP_MEP_NAME}" ]]; then
    svcPath="${MEEP_SANDBOX_NAME}/${MEEP_MEP_NAME}"
else
    svcPath="${MEEP_SANDBOX_NAME}"
fi

# Update API yaml basepaths to enable "Try-it-out" feature
# OAS2: Set relative path to sandbox name + endpoint path (origin will be derived from browser URL)
# OAS3: Set full path to provided Host URL + sandbox name + endpoint path
setBasepath() {
    # OAS3
    hostName=$(echo "${MEEP_HOST_URL}" | sed -E 's/^\s*.*:\/\///g')
    echo "Replacing [localhost] with ${hostName} to url in: '$1'"
    sed -i "s,localhost,${hostName},g" "$1";

    # OAS2 and OAS3
    echo "Replacing [sandboxname] with ${svcPath} to basepath or url in: '$1'"
    sed -i "s,sandboxname,${svcPath},g" "$1";
}

# Set basepath for API files
for file in /api/*; do
    if [[ ! -e "$file" ]]; then continue; fi
    setBasepath "$file"
done

# Set basepath for user-supplied API files
for file in /user-api/*; do
    if [[ ! -e "$file" ]]; then continue; fi
    setBasepath "$file"
done

# Start service
currenttime=`date "+%Y%m%d-%H%M%S"`
filepath="/codecov/codecov-meep-vis-"
filename=$filepath$currenttime".out"
if [ "$MEEP_CODECOV" = 'true' ]; then
  exec /meep-vis -test.coverprofile=$filename __DEVEL--code-cov
else
  exec /meep-vis
fi
