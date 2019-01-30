# AdvantEDGE Deployment Procedure
## Goals
- Guidance on deploying AdvantEDGE

## Overview
AdvantEDGE can be deployed using [this](link missing) script.

Before proceeding, make sure AdventEDGE environment has been setup [(details)](setup.md)

The script uses various Helm [charts](../charts) to deploy AdvantEDGE micro-services.

For help on script usage, type
```
deploy.sh
```

### Note
> As per the Docker/Kubernetes workflow, prior to deployment, AdvantEDGE micro-services must be stored in a Docker registry.

> This can be achieved either by using the DockerHub released images, by importing images in the local Docker registry or by building images from source (containerize in local registry).
