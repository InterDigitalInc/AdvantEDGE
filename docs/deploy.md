# Deployment Procedure
## Goal
- Guidance on deploying AdvantEDGE

## Overview
- setup [AdventEDGE Runtime Environment](setup_runtime.md)
- install and configure [meepctl CLI tool](meepctl/meepctl.md)
- ensure AdvantEDGE Docker images are available in a Docker Registry (see note below)
- `meepctl deploy all`

###### Note
> As per the Docker/Kubernetes workflow, Docker images must be stored in a Docker registry prior to being deployed. Released versions of AdvantEDGE micro-services are available on [DockerHub](missing link).

## AdvantEDGE deployment:
AdvantEDGE deployment is achieved through the [_meepctl CLI tool_](meepctl/meepctl.md)

AdvantEDGE is composed of a collection of micro-services that are classified in two groups: _core_ & _dependencies_. [_meepctl_](meepctl/meepctl.md) tool is used to create & destroy these micro-services on the K8s cluster; this is achieved through the [_deploy_](meepctl/meepctl_deploy.md) & [_delete_](meepctl/meepctl_delete.md) commands.

Let's see how it's done with the following examples.

Initially, deploy both groups using:
```
meepctl deploy all
```
When new AdvantEDGE version becomes available, only _core_ components need to be updated.
This is achieved by deleting and deploying the core group:
```
meepctl delete core
meepctl deploy core
```
alternatively
`meepctl deploy core --force` would achieve the same result

When finished using AdvantEDGE:
```
meepctrl delete all
```
###### Note
> AdvantEDGE dependencies are a pre-requisite needed by the core group. Therefore behavior is undefined if the dependency group is absent/deleted when core containers are deployed
