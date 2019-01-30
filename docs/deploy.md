# AdvantEDGE Deployment Procedure
## Goals
- Guidance on deploying AdvantEDGE

## Overview
AdvantEDGE deployment has been automated with [this](../deploy.sh) script.

The script uses [Helm](https://helm.sh/) to deploy AdvantEDGE micro-services using [these charts](../charts).

Before proceeding, make sure AdventEDGE environment has been [setup](setup.md)


Script usage:
```
deploy.sh
```

### Note
> As per the Docker/Kubernetes workflow, Docker images must be stored in a Docker registry prior to being deployed.<br/> Released versions of AdvantEDGE micro-services are available on [DockerHub](missing link).
