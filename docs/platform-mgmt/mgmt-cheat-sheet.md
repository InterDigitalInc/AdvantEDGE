---
layout: default
title: Deployment Cheatsheet
parent: Platform Management
nav_order: 2
---

Topic | Abstract
------|------
[Common Deployment Tasks](#common-deployment-tasks) | Common platform deployment tasks
[Deployment Configuration](#deployment-configuration) | AdvantEDGE deployment configuration
[Deployment files & paths](#deployment-files--paths) | Common platform resource installation locations
[Deployment ports & ingress](#deployment-ports--ingress) | Exposed platform ports & paths 

---
## Common Deployment Tasks
The following list is a summary of steps to follow to perform specific AdvantEDGE tasks.<br>
Requires _Runtime & Build environment setup_ + `git clone https://github.com/InterDigitalInc/AdvantEDGE.git`

Bootstrapping (first time & after repo upgrade)

No  | Description | Command
--- | --- | ---
1   | `./go-apps/meepctl/install.sh`        | Build & install meepctl CLI tool
2   | `meepctl config ip <your-ip>`         | Configure meepctl
3   | `meepctl config gitdir <path-to-git>` | Configure meepctl
4   | `meepctl deploy dep`                  | Deploy dependencies pods
5   | `meepctl build all`                   | Build core pods
6   | `meepctl dockerize all`               | Generate core pod containers
7   | `meepctl deploy core`                 | Deploy platform core pods

- _skip #2/3 on version change_
- _skip #1/2/3 on restart all_
- _skip #1/2/3/4 on recompile & restart without changing deps_

Shutting down

No  | Description | Command
--- | --- | ---
1   | `meepctl delete core`         | Deletes platform core pods
2   | `meepctl delete dep`          | Deletes dependencies pods

Manually cleaning un-terminated sandbox and scenario pods after deleting platform
```
helm ls --short -n <sandbox-name> | xargs -L1 helm uninstall -n <sandbox-name>
```
Manually cleaning only un-terminated scenario pods after deleting platform
```
helm ls -A --short | grep <scenario-name> | xargs -L1 helm uninstall -n <sandbox-name>
```

---
## Deployment Configuration
AdvantEDGE deployments may be configured by changing the default deployment settings in the `~/AdvantEDGE/.meepctl-repocfg.yaml` configuration file.

The following table describes the configurable deployment parameters and their default values.

parameter | description | default
-|-|-
`permissions.uid` | User identifier for platform permissions | `1000`
`permissions.gid` | Group identifier for platform permissions | `1000`
`user.frontend`  | Use user-supplied frontend<br> -> frontend UI located @ .meep/user/frontend | `false`
`user.swagger`  | Use user-supplied platform & sandbox swagger-ui<br> -> platform swagger UI located @ .meep/user/swagger<br> -> sandbox swagger UI located @ .meep/user/sandbox-swagger | `false`
`ingress.host`  | Host fully-qualified domain name or IP address | `my-platform-fqdn`
`ingress.https-only`  | Use secure connection (HTTPS) only | `false`
`ingress.host-ports`  | Use host ports (true) or node ports (false) | `true`
`ingress.http-port`  | HTTP port number | `80`
`ingress.https-port`  | HTTP port number | `443`
`ingress.ca`  | Certificate authority (none\|self-signed\|lets-encrypt) | `self-signed`
`ingress.le-server-prod`  | lets-encrypt production server (true) or staging server (false) | `false`
`auth.enabled` | Enable authentication & authorization | `true`
`auth.provider-mode` | Provider-authenticated users allowed (open) or limited to user database (secure) | `open`
`auth.session.key-secret`  | Session encryption key k8s secret<br> Data:<br> -> _encryption-key_: encryption key  | `meep-session`
`auth.session.max-sessions`  | Maximum simultaneous sessions | `10`
`auth.github.enabled`  | Use GitHub OAuth | `true`
`auth.github.auth-url`  | GitHub authorization URL | `https://github.com/login/oauth/authorize`
`auth.github.token-url`  | GitHub access token URL | `https://github.com/login/oauth/access_token`
`auth.github.redirect-uri`  | GitHub OAuth redirect URI | `https://my-platform-fqdn/platform-ctrl/v1/authorize`
`auth.github.secret`  | GitHub OAuth k8s secret<br> Data:<br> -> _client-id_: GitHub OAuth app client ID<br> -> _secret_: GitHub OAuth app secret | `meep-oauth-github`
`auth.gitlab.enabled`  | Use GitLab OAuth | `true`
`auth.gitlab.auth-url`  | GitLab authorization URL | `https://gitlab.com/oauth/authorize`
`auth.gitlab.token-url`  | GitLab access token URL | `https://gitlab.com/oauth/token`
`auth.gitlab.redirect-uri`  | GitLab OAuth redirect URI | `https://my-platform-fqdn/platform-ctrl/v1/authorize`
`auth.gitlab.secret`  | GitLab OAuth k8s secret<br> Data:<br> -> _client-id_: GitLab OAuth app client ID<br> -> _secret_: GitLab OAuth app secret | `meep-oauth-gitlab`
`metrics.influx.enabled`  | Enable influx data backups | `false`
`metrics.influx.url`  | Object store url | `my-object-store-fqdn`
`metrics.influx.secret`  | Object store configuration secret | `meep-influx-objstore-config`
`metrics.influx.retention`  | Number of days to retain daily data backups | `7`
`metrics.prometheus.external-labels.region`  | Deployment region (geographic or logical) | `idcc`
`metrics.prometheus.external-labels.monitor`  | Function being monitored | `advantedge`
`metrics.prometheus.external-labels.promenv`  | Prometheus environment (_dev_ or _prod_) | `prod`
`metrics.prometheus.external-labels.replica`  | Unique deployment identifier | `platform-ip`
`metrics.thanos.enabled`  | Enable Thanos | `false`
`metrics.thanos.secret`  | Object store configuration secret | `meep-thanos-objstore-config`
`metrics.thanos.query.enabled`  | Enable querier | `true`
`metrics.thanos.query-frontend.enabled`  | Enable query frontend | `true`
`metrics.thanos.store-gateway.enabled`  | Enable store gateway | `true`
`metrics.thanos.compactor.enabled`  | Enable compactor | `false`
`metrics.thanos.compactor.retention.resolution-raw`  | Raw data retention | `30d`
`metrics.thanos.compactor.retention.resolution-5m`  | 5m downsampled data retention | `60d`
`metrics.thanos.compactor.retention.resolution-1h`  | 1h downsampled data retention | `10y`
`metrics.thanos.thanos-archive.enabled`  | Enable Thanos archive | `false`
`metrics.thanos.thanos-archive.secret`  | Archive object store configuration secret | `meep-thanos-archive-objstore-config`

Dependency microservices and Core & Sandbox Subsystem microservices may also be modified using this configuration file. For example, specific microservices can be excluded from a platform build or deployment if not required.

_**NOTE:** Modifying microservice configuration is only recommended for advanced platform users_

---
## Deployment files & paths
The following are common locations where AdvantEDGE resources are installed; path may vary depending on your environment

Location | Usage
--- | ---
`~/.meepctl.yaml`                     | meepctl CLI tool main configuration file
`~/AdvantEDGE/.meepctl-repocfg.yaml`  | meepctl CLI tool repo & deployment configuration file
`~/gocode/bin/meepctl`                | meepctl installation folder
`~/.meep/`                            | platform files, pods persistent storage
`~/.meep/certs/`                      | certificates
`~/.meep/codecov/`                    | code coverage reports
`~/.meep/codecov-bak/`                | previous code coverage reports
`~/.meep/couchdb/`                    | couchdb persistent storage
`~/.meep/docker-registry/`            | docker registry persistent storage
`~/.meep/grafana/`                    | grafana persistent storage
`~/.meep/influxdb/`                   | influxdb persistent storage
`~/.meep/omt/`                        | open-map-tiles persistent storage
`~/.meep/postgis/`                    | postgis persistent storage
`~/.meep/prometheus/alertmanager`     | prometheus alert manager persistent storage
`~/.meep/prometheus/server`           | prometheus server persistent storage
`~/.meep/tmp/`                        | temporary meepctl work directory
`~/.meep/user/`                       | user provided resources (fe, swagger, chart-values)
`~/.meep/virt-engine/`                | virt-engine persistent storage

---
## Deployment ports & ingress
AdvantEDGE platform API and frontend are served via an ingress controller on port 80 & 443.

The following tables present a summary of the service exposure.

### DEPENDENCY

Module | type | default
-|-|-
couchdb             |internal | -
docker-registry     |port     | 30001
grafana             |ingress  | `/grafana`
influxdb            |internal | -
kube-state-metrics  |internal | -
nginx-ingress       |port     | 80/443
open-map-tiles      |ingress  | `/map`<br>`/styles`<br>`/images`<br>`/fonts`<br>`/data`<br>`/leaflet-hash`
postgis             |internal | -
prometheus          |internal | -
redis               |internal | -

### PLATFORM

Module | type | default
-|-|-
meep-ams            |ingress  | `/<sandbox-name>/amsi`
meep-app-enablement |ingress  | `/<sandbox-name>/mec_app_support`<br>`/<sandbox-name>/mec_service_mgmt`
meep-auth-svc       |ingress  | `/auth`
meep-gis-engine     |ingress  | `/<sandbox-name>/gis`
meep-loc-serv       |ingress  | `/<sandbox-name>/location`
meep-metrics-engine |ingress  | `/<sandbox-name>/metrics`
meep-mg-manager     |ingress  | `/<sandbox-name>/mgm`
meep-mon-engine     |ingress  | `/mon-engine`
meep-platform-ctrl  |ingress  | `/`<br>`/api`<br>`/platform-ctrl`
meep-rnis           |ingress  | `/<sandbox-name>/rni`
meep-sandbox-ctrl   |ingress  | `/<sandbox-name>/api`<br>`/<sandbox-name>/sandbox-ctrl`
meep-tc-engine      |internal | -
meep-tc-sidecar     |internal | -
meep-virt-engine    |internal | -
meep-wais           |ingress  | `/<sandbox-name>/wai`
meep-webhook        |internal | -
