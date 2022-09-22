#!/bin/bash
set -e

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
mv /meep-vis /templates/sandbox/meep-vis

mkdir -p /templates/scenario
mv /meep-virt-chart-templates /templates/scenario/meep-virt-chart-templates

# Configure & update helm repo
helm repo add incubator https://charts.helm.sh/incubator
helm repo update

# Start virt engine
exec /meep-virt-engine
