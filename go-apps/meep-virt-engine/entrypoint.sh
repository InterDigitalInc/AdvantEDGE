#!/bin/bash
set -e
# echo "MEEP_HOST_URL: ${MEEP_HOST_URL}"
echo "MEEP_CODECOV: ${MEEP_CODECOV}"
echo "MEEP_CODECOV_LOCATION: ${MEEP_CODECOV_LOCATION}"

# Move helm charts from docker data to template folder
mkdir -p /templates/sandbox
mv /meep-gis-engine /templates/sandbox/meep-gis-engine
mv /meep-loc-serv /templates/sandbox/meep-loc-serv
mv /meep-metrics-engine /templates/sandbox/meep-metrics-engine
mv /meep-mg-manager /templates/sandbox/meep-mg-manager
mv /meep-rnis /templates/sandbox/meep-rnis
mv /meep-app-enablement /templates/sandbox/meep-app-enablement
mv /meep-wais /templates/sandbox/meep-wais
mv /meep-ams /templates/sandbox/meep-ams
mv /meep-sandbox-ctrl /templates/sandbox/meep-sandbox-ctrl
mv /meep-tc-engine /templates/sandbox/meep-tc-engine

mkdir -p /templates/scenario
mv /meep-virt-chart-templates /templates/scenario/meep-virt-chart-templates

# Configure & update helm repo
helm repo add incubator https://charts.helm.sh/incubator
helm repo update

# Start virt engine
currenttime=`date "+%Y%m%d-%H%M%S"`
filepath="/codecov/codecov-meep-virt-engine-"
filename=$filepath$currenttime".out"
if [ "$MEEP_CODECOV" = 'true' ]; then
  MEEP_CODECOV=${MEEP_CODECOV} MEEP_CODECOV_LOCATION=${MEEP_CODECOV_LOCATION} exec /meep-virt-engine -test.coverprofile=$filename __DEVEL--code-cov
else
  exec /meep-virt-engine
fi
