---
layout: default
title: Deployment Workflow
parent: Platform Management
nav_order: 1
---

Topic | Abstract
------|------
[Deployment Workflow](#deployment-workflow) | An overview of the workflow to deploy the platform
[Build](#build) | Build AdvantEDGE
[Deploy](#deploy) | Deploy AdvantEDGE
[Upgrade](#upgrade) | Upgrade AdvantEDGE
[Test](#test) | Test AdvantEDGE
NEXT STEP: [Platform usage](#next-step) |

----
## Deployment Workflow
![mgmt-flow]({{site.baseurl}}/assets/images/mgmt-flow.jpg)

To help getting started, the above figure presents typical workflows of how to manage the AdvantEDGE platform.
- The middle flow present steps required on the first install to get the platform up & running
- The top flow presents steps to stop/start AdvantEDGE in-between uses
- The bottom flow presents steps to upgrade AdvantEDGE

Details on each of these steps can be found either in Environment Setup pages or here, in the Platform Management page.

Workflow on AdvantEDGE usage is available in the Platform Usage section.

----
## Build
This procedure
- clones repository
- builds & installs meepctl tool
- builds frontend & micro-services
- _pre-requisites_
  - _AdvantEDGE Development Environment_

### Clone repository
AdvantEDGE repository follows the [Gitflow Workflow](https://nvie.com/posts/a-successful-git-branching-model/) branching model for sharing official platform releases and development updates. Key branches are:
- **_master:_** Official platform release branch
- **_develop:_** Development branch with latest features
- **_feature:_** Long-lived feature branch

We recommend using the master branch

```
git clone https://github.com/InterDigitalInc/AdvantEDGE.git
```

### Build & install meepctl
The bash script below buids & installs [_meepctl CLI tool_](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/docs/meepctl/meepctl.md)

```
cd ~/AdvantEDGE/go-apps/meepctl
./install.sh
```

On first install, meepctl needs to be configured.
```
meepctl config ip <your-node-ip>
meepctl config gitdir /home/<user>/AdvantEDGE

# To view current meepctl configuration
meepctl config
```

### Build frontend & micro-services
meepctl CLI tool is used to build frontend and micro-services.

```
meepctl build all

# To build a single micro-service:
# meepctl build <micro-service name>

# Linters are executed by default on every build.
# To disable linting use the `--nolint` flag.
```

----
## Deploy
AdvantEDGE micro-services are classified in two groups: _core_ & _dependencies_; behavior is undefined if the _dependencies_ group is absent/deleted when core containers are deployed.

This procedure
- configures deployment (optional)
- deploys the _dependencies_
- containerize _core_ micro services
- deploys _core_ micro-services
- _pre-requisites_
  - _AdvantEDGE Runtime Environment_
  - _meepctl CLI tool installed_

### Configure deployment (optional)
AdvantEDGE comes with a [default configuration](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/.meepctl-repocfg.yaml) that can be used out-of-the-box for private deployments.

For custom deployments, the configuration file can be edited to control deployment parameters - see [Deployment Configuration]({{site.baseurl}}{% link docs/platform-mgmt/mgmt-cheat-sheet.md %}#deployment-configuration)

### Deploy dependencies
meepctl CLI tool is used to deploy AdvantEDGE dependencies on Kubernetes.

```
meepctl deploy dep

# To delete dependencies
# meepctl delete dep
```

### Containerize core micro-services
meepctl CLI toolis used to containerize AdvantEDGE micro-services.

```
meepctl dockerize all

# To verify that it worked
# docker images | grep meep
```

_**NOTE:**_
- _this command first creates the containers in the local docker registry_
- _then pushes the images in the K8s registry; therefore make sure dependencies are running_

### Deploy core micro-services
meepctl CLI tool is used to deploy AdvantEDGE core micro-services on Kubernetes.

```
meepctl deploy core

# To delete only core micro-services
# meepctl delete core

# To delete and deploy in one operation
# meepctl deploy core --force

# When done using
# meepctl delete core
# meepctl delete dep
```

Our favorite command to verify if everything is running
```
kubectl get pods --all-namespaces | grep meep
default       meep-auth-svc-68fb4dbffd-lh959                      1/1     Running   0          26h
default       meep-couchdb-couchdb-0                              1/1     Running   0          29h
default       meep-docker-registry-65b77797cb-tp665               1/1     Running   0          29h
default       meep-grafana-667984f74b-8k577                       1/1     Running   0          29h
default       meep-influxdb-0                                     1/1     Running   0          29h
default       meep-ingress-controller-shxrs                       1/1     Running   0          29h
default       meep-ingress-defaultbackend-5c57d5cd58-4ktwg        1/1     Running   0          29h
default       meep-kube-state-metrics-868576f6d4-hq7lb            1/1     Running   0          29h
default       meep-mon-engine-6b75855c74-4vcj5                    1/1     Running   0          26h
default       meep-open-map-tiles-7d99b886f-k5ndg                 1/1     Running   0          29h
default       meep-platform-ctrl-5994bb5868-9gl9d                 1/1     Running   0          26h
default       meep-postgis-0                                      2/2     Running   0          29h
default       meep-prometheus-couchdb-exporter-795d6b6dc5-csvfr   1/1     Running   0          29h
default       meep-prometheus-node-exporter-62qbw                 1/1     Running   0          29h
default       meep-prometheus-operator-c8b8896d7-vnlrb            1/1     Running   0          29h
default       meep-redis-master-0                                 2/2     Running   0          29h
default       meep-redis-slave-0                                  2/2     Running   0          29h
default       meep-virt-engine-6f44488b54-sk86p                   1/1     Running   0          26h
default       meep-webhook-6865678784-4ntnb                       1/1     Running   0          26h
```

Alternatively - a green indicator on the top right corner of the frontend indicates that all the pods are running.<br>
If some pods are not, indicator will be red & hovering indicates which pods are missing.

----
## Upgrade
This procedure
- uninstalls AdvantEDGE
- follow [Build procedure](#build)
- follow [Deploy procedure](#deploy)

### Uninstall AdvantEDGE
meepctl CLI tool is used to uninstall AdvantEDGE.

```
meepctl delete core
meepctl delete dep
```
_**NOTE:** meepctl CLI tool performs a version check in the local **.meepctl-repocfg.yaml** file. For this reason, it is recommended to uninstall AdvantEDGE before fetching the latest release._

----
## Test
AdvantEDGE currently supports end-to-end testing using [Cypress](https://www.cypress.io/). This Node-based JavaScript testing tool simulates user interactions with the frontend and validates expected UI updates.

This procedure
- installs Cypress
- runs unit tests
- runs Cypress CLI
- alternatively, runs Cypress GUI

### Install Cypress
To install Cypress run the following commands:

```
cd ~/AdvantEDGE/test
npm ci
```

### Run unit tests
```
cd AdvantEDGE/tests
./start-ut-env.sh
./run-ut.sh
./stop-ut-env.sh
```

### Run Cypress CLI
```
# Run Cypress tests using CLI
cd ~/AdvantEDGE/test
./run-cypress.sh

# Default AdvantEDGE URL used by cypress is http://127.0.0.1
# To run tests using another deployment:
# npm run cy:run -- --env meep_url="http://<Node IP>"
```

### Cypress GUI
```
# Run/Debug Cypress tests using GUI
cd ~/AdvantEDGE/test
npm run cy:open

# Default AdvantEDGE URL used by cypress is http://127.0.0.1
# To run tests using another deployment:
# npm run cy:open -- --env meep_url="http://<Node IP>"

```

_**NOTE:** Cypress may crash if max inotify watchers is too low. To fix this run the command:<br>
`echo fs.inotify.max_user_watches=524288 | sudo tee -a /etc/sysctl.conf && sudo sysctl -p`<br>
See details [here](https://github.com/guard/listen/wiki/Increasing-the-amount-of-inotify-watchers)_

### Code Coverage
**(THIS SECTION IS OUTDATED - CODE COVERAGE NEEDS TO BE REWORKED)**

AdvantEDGE core micro-services can be instrumented with code coverage instrumentation; when used in conjunction with Cypress or other system tests (manual or proprietary), it will provide an overview of the code coverage.

The following is a summary of the steps to enable code coverage measurement in AdvantEDGE:
- Build for code coverage: `meepctl build all --codecov`
- Dockerize: `meepctl dockerize all`
- Deploy for code coverage: `meepctl deploy core --codecov`
- Execute testing - use Cypress, manual or any desired system test

Once testing is completed
- **Stop the micro-services gracefully**: `meepctl delete core`
  _Build, dockerize  & deploy will instrument and execute core micro-services so they measure code coverage._
  _When terminated gracefully, the core micro-services will store the code coverage result in the following location: `~/.meep/codecov/<micro-service-name>/codecov-<micro-service-name>.out`_
- For convenience, code coverage reports can be generated using `meepctl test`


----
## Next Step
[Use AdvantEDGE platform]({{site.baseurl}}{% link docs/usage/usage-workflow.md %}):
- Create scenarios
- Execute scenarios
- Observe application beahvior
- etc.
