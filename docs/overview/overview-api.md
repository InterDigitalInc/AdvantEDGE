---
layout: default
title: API
parent: Overview
nav_order: 4
---

Topic | Abstract
------|------
[Getting started with APIs](#getting-started-with-apis) | Finding/Viewing/Experimenting/Interfacing  AdvantEDGE APIs
[Platform  API](#platform-api) | [OAS2] Scenario and Sandbox Management endpoints
[Sandbox  API](#sandbox-api) | [OAS2] Scenario activation, Events, Connectivity & Application instance ID management endpoints
[Monitoring API](#monitoring-api) | [OAS2] Microservice & scenario deployment status information endpoints
[GIS API](#gis-api) | [OAS2] Geospatial information and automation endpoints
[Metrics API](#metrics-api) | [OAS2] Network metrics and Events query and subscription endpoints
[Metrics Notification API](#metrics-notification-api) | [OAS2] Client side metrics notification endpoints
[Location Service API](#location-service-api) | [OAS3] Location Service northbound API endpoints
[RNIS API](#rnis-api) | [OAS3] RNIS northbound API endpoints
[WAIS API](#wais-api) | [OAS3] WAIS northbound API endpoints
[App Support API](#app-support-api) | [OAS3] App Support northbound API endpoints
[Service Management API](#service-management-api) | [OAS3] Service Management northbound API endpoints
[AMS API](#ams-api) | [OAS3] Application Mobility Service northbound API endpoints
[Application State Transfer Service API](#application-state-transfer-service-api) | [OAS2] Mobility Group Manager membership and state transfer endpoints
[Application State Transfer Notification API](#application-state-transfer-notification-api) | [OAS2] Client side Mobility Group Manager state transfer endpoints
[VIS API](#vis-api) | [OAS3] VIS northbound API endpoints
NEXT STEP: [Recommended hardware](#next-step) | |

----
## Getting Started with APIs
AdvantEDGE backend offers a series of REST APIs that can be used to achieve different tasks.
- AdvantEDGE frontend uses certain APIs to enable user interactions with the backend platform from the browser environment
- User software may also use these APIs for different reasons, some examples follow:
  - to trigger events according to their scenario applications (ex. mobility/network characteristics)
  - to create PDU Sessions (ex, when experimenting cellular scenarios)
  - to reads metrics (ex: plot graphs, experiment with smart network algorithms)

AdvantEDGE APIs follow the [OpenAPI Specification](https://github.com/OAI/OpenAPI-Specification) (OAS) standard to define the platform APIs. Some of our APIs are aligned on OAS 2.0 version of the standard (formerly known as Swagger 2.0) while others are aligned with OAS 3.0.

Following the OAS standard, allows AdvantEDGE APIs to benefit from a rich eco-system of tools to facilitate integration as described below.

### Finding API Specification
In the GitHub repository, each application exposing a REST API has an associated `swagger.yaml` file.
- For Golang micro-services it can be found under `go-apps/<micro-service>/api/`
- For Golang packages it can be found under `go-packages/<package>/api/`
- For Javascript packages ir can be found under `js-packages/<package>/src/api/`

_**NOTE:** that not every component implements a REST API._

### Viewing API Specification
While YAML format is a convenient way to define an API, it is not a user friendly one to look at the API and understand its structure and details.

For that reason, AdvantEDGE provides Swagger-UI, a web browser viewer. After deploying AdvantEDGE platform, access the `your-ip-address/api/` or `your-ip-address/your-sandbox-name/api/` to see respectively the plarform level APIs or the Sandbox level APIs.

Swagger-UI (shown below) provides a convenient point-and-click interface to view API details (by clicking an operation) and a drop down menu (top right) to select the desired API.

![swagger-ui]({{site.baseurl}}/assets/images/swagger-ui.png)

### Experimenting with API
There are various ways to experiment with AdvantEDGE APIs.

One convenient way is to do it directly from Swagger-UI.
As shown below, when you expand an operation to view details, the Swagger-UI offers a _Try It Out_ button (right side) that allows to send the corresponding REST request to the live running system.

![swagger-try-it-out]({{site.baseurl}}/assets/images/swagger-try-it-out.png)

This is useful to quickly ramp-up on API functionality.

Alternatively, CLI curl or numerous browser integrated HTTP probing solutions can be used to achieve the same result.

_**IMPORTANT NOTE:**<br>
**Try It Now** functionality is not supported for Notification APIs since these have  to be implemented by the service user.<br>
For example, **Location Service Subscription Notification Callback REST API** is an API implemented by a user of the **Location Service REST API**<br>
Therefore, **Try It Now** works for **Location Service** but not for **Location Service Subscription Notification Callback**_

### Interfacing with API
Interfacing with AdvantEDGE APIs can be done following different methods.

Using the specifications, a user could invoke AdvantEDGE HTTP endpoints direclty from software (using a HTTP library) or from a script (using curl or similar tool). This method works but is not optimal.

A convenient alternative is using a code-generator to create a client package from the API specification. The package is subsequently integrated in your user codebase. As maintainers of this project, this what we do. Below, we detail the approach to help achiving this task for your user application.

For OAS2 specifications, we use swagger-codegen v2.4.16 from [this repo](https://github.com/swagger-api/swagger-codegen)
```
wget http://central.maven.org/maven2/io/swagger/swagger-codegen-cli/2.4.16/swagger-codegen-cli-2.4.16.jar -O swagger-codegen-cli.jar
```

For OAS3 specifications, we use swagger-codegen v3.0.22 from [this repo](https://github.com/swagger-api/swagger-codegen)
```
wget http://central.maven.org/maven2/io/swagger/swagger-codegen-cli/3.0.22/swagger-codegen-cli-3.0.22.jar -O swagger-codegen-cli.jar
```

_**NOTE v2 vs v3 swagger-codegen**<br>As explained [here](https://github.com/swagger-api/swagger-codegen#versioning), both version support OAS 2.0 but differently<br>- v2 swagger-codegen supports natively OAS 2.0 specs<br>- v3 swagger-codegen supports natively OAS 3.0 specs and supports OAS 2.0 via spec conversion prior to code generation<br>_


Swagger-codegen can generate client packages for many different languages:
```
java -jar swagger-codegen-cli.jar
Available languages: [ada, ada-server, akka-scala, android, apache2, apex, aspnetcore, bash, csharp, clojure, cwiki, cpprest, csharp-dotnet2, dart, dart-jaguar, elixir, elm, eiffel, erlang-client, erlang-server, finch, flash, python-flask, go, go-server, groovy, haskell-http-client, haskell, jmeter, jaxrs-cxf-client, jaxrs-cxf, java, inflector, jaxrs-cxf-cdi, jaxrs-spec, jaxrs, msf4j, java-pkmst, java-play-framework, jaxrs-resteasy-eap, jaxrs-resteasy, javascript, javascript-closure-angular, java-vertx, kotlin, lua, lumen, nancyfx, nodejs-server, objc, perl, php, powershell, pistache-server, python, qt5cpp, r, rails5, restbed, ruby, rust, rust-server, scala, scala-gatling, scala-lagom-server, scalatra, scalaz, php-silex, sinatra, slim, spring, dynamic-html, html2, html, swagger, swagger-yaml, swift4, swift3, swift, php-symfony, tizen, typescript-aurelia, typescript-angular, typescript-inversify, typescript-angularjs, typescript-fetch, typescript-jquery, typescript-node, undertow, ze-ph, kotlin-server]
```

Help from the tool can be obtained via
```
# general help
java -jar swagger-codegen-cli.jar help

# generate command help
java -jar swagger-codegen-cli.jar help generate
```

Finally, here are some code-gen examples for generating the AdvantEDGE Controller client package for various languages
```
# Generating Controller Client API (go)'
java -jar swagger-codegen-cli.jar generate -i meep-platform-ctrl-api.yaml -l go -o ./platform-ctrl-client/go -DpackageName=client
# Generating Controller Engine Client API (python)'
java -jar swagger-codegen-cli.jar generate -i meep-platform-ctrl-api.yaml -l python -o ./platform-ctrl-client/python -DpackageName=client
# Generating Controller Engine Client API (javascript)'
java -jar swagger-codegen-cli.jar generate -i meep-platform-ctrl-api.yaml -l javascript -o ./platform-ctrl-client/js/
```
How to use the generated package from your codebase is described in the generated README file.

----
## Platform API
This API allows to perform CRUD operations on the scenarios & sandboxes

API:
- From repository: [meep-platform-ctrl (yaml)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/go-apps/meep-platform-ctrl/api/swagger.yaml)
- From wiki: [Platform Controller (markdown)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/docs/api-platform-ctrl/README.md)
- From browser: `https://<your-advantedge-ip>/api`

----
## Sandbox API
This API allows to control scenario activation, send events, control PDU session connectivity and manage Application instance IDs

API:
- From repository: [meep-sandbox-ctrl (yaml)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/go-apps/meep-sandbox-ctrl/api/swagger.yaml)
- From wiki: [Sandbox Controller (markdown)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/docs/api-sandbox-ctrl/README.md)
- From browser: `https://<your-advantedge-ip>/<your-sandbox>/api`

----
## Monitoring API
This API allows to obtain micro-service deployment and status information

API:
- From repository: [meep-mon-engine (yaml)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/go-apps/meep-mon-engine/api/swagger.yaml)
- From wiki: [Monitoring Engine (markdown)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/docs/api-mon-engine/README.md)
- From browser: `https://<your-advantedge-ip>/api`

----
### GIS API
This API allows to obtain geospatial information and control geospatial automations

API:
- From repository: [meep-gis-engine (yaml)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/go-apps/meep-gis-engine/api/swagger.yaml)
- From wiki: [GIS Engine (markdown)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/docs/api-gis-engine/README.md)
- From browser: `https://<your-advantedge-ip>/<your-sandbox>/api`

----
### Metrics API
This API allows to obtain network metrics, event metrics by querying or subscribing.

API:
- From repository: [meep-metrics-engine (yaml)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/go-apps/meep-metrics-engine/api/v2/swagger.yaml)
- From wiki: [Metrics Engine (markdown)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/docs/api-metric-engine-v2/README.md)
- From browser: `https://<your-advantedge-ip>/<your-sandbox>/api`

----
### Metrics Notification API
This API must be implemented by a client subscribing to the Metrics API.

API:
- From repository: [meep-metrics-engine-notification-client (yaml)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/go-packages/meep-metrics-engine-notification-client/api/swagger.yaml)
- From wiki: [Metrics Engine Notification (markdown)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/docs/api-metric-engine-notif/README.md)
- From browser: `https://<your-advantedge-ip>/<your-sandbox>/api`

----
### Location Service API
This API allows to obtain location information by querying or subscribing; see [Location Service]({{site.baseurl}}{% link docs/overview/edge-services/overview-edge-services.md %}#location-service) for service description.

API:
- From repository: [meep-loc-serv (yaml)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/go-apps/meep-loc-serv/api/swagger.yaml)
- From wiki: [Location Service (markdown)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/docs/api-location/README.md)
- From browser: `https://<your-advantedge-ip>/<your-sandbox>/api`

----
### RNIS API
This API allows to obtain radio network information by querying or subscribing; see [RNIS]({{site.baseurl}}{% link docs/overview/edge-services/overview-edge-services.md %}#radio-network-information-service) for service description.

API:
- From repository: [meep-rnis (yaml)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/go-apps/meep-rnis/api/swagger.yaml)
- From wiki: [RNIS (markdown)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/docs/api-rnis/README.md)
- From browser: `https://<your-advantedge-ip>/<your-sandbox>/api`

----
### WAIS API
This API allows to obtain WLAN network information by querying or subscribing; see [WAIS]({{site.baseurl}}{% link docs/overview/edge-services/overview-edge-services.md %}#wireless-access-information-service) for service description.

API:
- From repository: [meep-wais (yaml)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/go-apps/meep-wais/api/swagger.yaml)
- From wiki: [WAIS (markdown)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/docs/api-wais/README.md)
- From browser: `https://<your-advantedge-ip>/<your-sandbox>/api`

----
### App Support API
This API enables application start-up/termination assistance by querying, subscribing and posting; see [Edge Platform Application Enablement Service]({{site.baseurl}}{% link docs/overview/edge-services/overview-edge-services.md %}#edge-platform-application-enablement-service) for service description.

API:
- From repository: [meep-app-support (yaml)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/go-apps/meep-app-enablement/api/app-support/swagger.yaml)
- From wiki: [App Support (markdown)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/docs/api-app-support/README.md)
- From browser: `https://<your-advantedge-ip>/<your-sandbox>/api`

----
### Service Management API
This API allows edge applications to discover, advertise, consume and offer MEC services by querying, subscribing and posting; see [Edge Platform Application Enablement Service]({{site.baseurl}}{% link docs/overview/edge-services/overview-edge-services.md %}#edge-platform-application-enablement-service) for service description.

API:
- From repository: [meep-service-mgmt (yaml)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/go-apps/meep-app-enablement/api/service-mgmt/swagger.yaml)
- From wiki: [Service Mgmt (markdown)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/docs/api-service-mgmt/README.md)
- From browser: `https://<your-advantedge-ip>/<your-sandbox>/api`

----
### AMS API
This API allows edge applications to relocate user context and/or application instance across MEC platforms by querying, subscribing and posting; see [Application Mobility Service]({{site.baseurl}}{% link docs/overview/edge-services/overview-edge-services.md %}#application-mobility-service) for service description.

API:
- From repository: [meep-ams (yaml)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/go-apps/meep-ams/api/swagger.yaml)
- From wiki: [AMS (markdown)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/docs/api-ams/README.md)
- From browser: `https://<your-advantedge-ip>/<your-sandbox>/api`

----
### Application State Transfer Service API
This API allows to manage group membership and control application state transfer; see [Application State Transfer Service]({{site.baseurl}}{% link docs/overview/edge-services/overview-edge-services.md %}#application-sate-transfer-service) for service description.

API:
- From repository: [meep-mg-manager (yaml)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/go-apps/meep-mg-manager/api/swagger.yaml)
- From wiki: [MG Manager (markdown)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/docs/api-mg-manager/README.md)
- From browser: `https://<your-advantedge-ip>/<your-sandbox>/api`

----
### Application State Transfer Notification API
This API must be implemented by a client subscribing to the Application State Transfer API.

API:
- From repository: [meep-mg-manager-client (yaml)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/go-packages/meep-mg-app-client/api/swagger.yaml)
- From wiki: [MG Manager Notification (markdown)](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/docs/api-mg-manager-notif/README.md)
- From browser: `https://<your-advantedge-ip>/<your-sandbox>/api`


## Next Step
Learn about the [recommended hardware]({{site.baseurl}}{% link docs/setup/env-hw.md %}):
