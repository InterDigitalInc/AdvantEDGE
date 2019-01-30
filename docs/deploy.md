# AdvantEDGE Deployment Procedure
## Goals
- Guidance on deploying AdvantEDGE

## Overview
AdvantEDGE can be deployed using [this](../deploy.sh) script.

Before proceeding, make sure AdventEDGE environment has been [setup](setup.md)

The script uses [Helm](https://helm.sh/) charts to deploy AdvantEDGE micro-services; the charts are available [here](../charts).

Script usage:
```
deploy.sh
```

### Note
> As per the Docker/Kubernetes workflow, prior to deployment, AdvantEDGE micro-services must be stored in a Docker registry.<br/> This can be achieved either by using DockerHub released images, by importing images in the local Docker registry or by building images from source (*images in local registry*).
