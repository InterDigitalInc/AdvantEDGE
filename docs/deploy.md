# Deployment Procedure
## Goal
- Guidance on deploying AdvantEDGE

## Pre-requisites
- setup [AdventEDGE Runtime Environment](setup_runtime.md)
- install and configure [meepctl CLI tool](meepctl/meepctl.md)

## Summary
- [`meepctl dockerize`](#dockerization)
- [`meepctl deploy all`](#deployment)

###### Note
> As per the Docker/Kubernetes workflow, Docker images must be stored in a Docker registry prior to being deployed.

## Dockerization
Prior to deploying AdvantEDGE, an intermediate step is needed.
Container images of AdvantEDGE micro-services need to be generated and stored in the local Docker registry.

[_meepctl_](meepctl/meepctl.md) tool is used to dockerize these micro-services and store them in the local docker registry; this is achieved through the [_dockerize_](meepctl/meepctl_dockerize.md) command.

```
meepctl dockerize
```

To verify that the operation was successful, you can list the local registry content and verify the creation time of the AdvantEDGE containers
```
docker images | grep meep
```

## Deployment
AdvantEDGE deployment is achieved through the [_meepctl CLI tool_](meepctl/meepctl.md)

AdvantEDGE micro-services are classified in two groups: _core_ & _dependencies_. [_meepctl_](meepctl/meepctl.md) tool is used to create & destroy these micro-services on the K8s cluster; this is achieved through the [_deploy_](meepctl/meepctl_deploy.md) & [_delete_](meepctl/meepctl_delete.md) commands.

Let's see how it's done with the following examples.

Initially, deploy both groups using:
```
meepctl deploy all
```
Typically, when new AdvantEDGE version become available, only _core_ components need to be updated.
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
