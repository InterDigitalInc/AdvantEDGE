#!/bin/bash
set -e

# Configure & update helm repo
helm repo add incubator https://charts.helm.sh/incubator
helm repo update

# Start virt engine
exec /meep-virt-engine
