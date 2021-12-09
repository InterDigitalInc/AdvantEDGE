---
layout: default
title: Edge Platform Application Enablement Service
parent: EDGE Services
grand_parent: Overview
nav_order: 4
permalink: docs/overview/edge-services/app-enablement/
---

## Service Overview
AdvantEDGE provides a built-in Edge Platform Application Enablement Service implementation that integrates with scenarios.

Mp1 reference point provides two different APIs: _MEC Application Support_ and _MEC Service Management_

These APIs allow MEC Applications to interact with the MEC System, such as:
- _Application registration/deregistration_ (supported)
- _Service discovery & offering_ (supported)
- _Event notifications about service and application availability_ (supported)
- _Traffic rules, DNS_ (not supported)
- _Time of day_ (supported)

## Micro-Services
  - _App Enablement Service:_ Implements ETSI MEC011 northbound APIs with a custom integration with AdvantEDGE APIs

## API Version
- Edge Platform Application Enablement Service is compliant with the ETSI MEC011 Specification, v2.1.1:
  - [ETSI GS MEC 011 V2.1.1](https://www.etsi.org/deliver/etsi_gs/MEC/001_099/011/02.01.01_60/gs_mec011v020101p.pdf)
  - [ETSI Forge - MEC Application Support API and MEC Service Management API](https://forge.etsi.org/rep/mec/gs011-app-enablement-api)
- API
  - [Application Support API Definition](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/docs/api-app-support)
  - [Service Management API Definition](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/docs/api-service-mgmt)
  - Based on OpenAPI Specification (OAS) 3.0

## AdvantEDGE Integration
Edge Platform Application Enablement Service is implemented as a single sandbox pod within AdvantEDGE, providing service for all applications running as part of that sandbox.

To use this service, MEC applications must first obtain an application instance ID from the AdvantEDGE platform using one of the following methods:
- For MEC applications configured in the deployed scenario:
  - From the ```MEEP_APP_ID``` environment variable, or
  - From the process table in the frontend Execution page
- For external MEC applications not configured in the deployed scenario:
  - From the _Applications_ endpoints of the [Sandbox Controller API]({{site.baseurl}}{% link docs/overview/overview-api.md %}#sandbox-api)

After provisioning the application instance ID, MEC applications must use the _Application Support API_ to confirm application readiness using the following procedure:
- Confirm application readiness:
  - Send _READY_ indication using application instance ID
  - ```POST .../mec_app_support/v1/applications/{appInstanceId}/confirm_ready```
- Register for graceful application termination (optional)
  - Subscribe for _AppTermination_ notification
  - When application is terminated, a notification will be sent with a grace period before clearing application resources
  - ```POST .../mec_app_support/v1/applications/{appInstanceId}/subscriptions```

Once confirmed, MEC Applications may use the _Service Management API_ to discover and register services as described in the use cases below.

### Use case for service-consuming MEC Application
This use case applies to MEC Applications, with no prior knowledge of MEC services offered by MEC platforms, wishing to discover, monitor & use available MEC services.
- Discover MEC services:
  - Retrieve service list from MEC platform
  - ```GET .../mec_service_mgmt/v1/services```
- Monitor MEC services:
  - Subscribe for service availability change notifications
  - Provide filter criteria (e.g. names, categories, etc.) to specify which services to watch
  - ```POST .../mec_service_mgmt/v1/applications/{appInstanceId}/subscriptions```
- Use MEC Services:
  - When target MEC service becomes available, MEC application uses the service
  - When target MEC service becomes unavailable, MEC application continues to monitor the service until it returns or another instance becomes available

### Use case for service-offering MEC Application
This use case applies to MEC Applications wishing to offer MEC services through MEC platforms.
- Register MEC services:
  - Provide MEC service information to the MEC platform
  - ```POST .../mec_service_mgmt/v1/applications/{appInstanceId}/services```
- Handle MEC Service requests:
  - MEC Service is now discoverable & reachable (according to configured scope) by other MEC applications
  - While available, service must process & respond to MEC application requests
- Stop MEC Service:
  - When service is no longer needed, inform the MEC platform to remove the service
  - MEC Service is no longer discoverable & reachable
  - ```DELETE .../mec_service_mgmt/v1/applications/{appInstanceId}/services```

### Use case for MEC Application using "scope of locality"
This use case applies to MEC Applications with multiple instances running on different MEC platforms with different localities. Configured scope of locality is used to determine which MEC Service instances are discoverable & reachable to MEC applications running on a MEC platform.
- Start multiple MEC Application instances
  - Each instance must obtain a unique application instance ID & confirm application readiness
- Register MEC services:
  - Each instance provides MEC service information to the MEC platform
  - Service configuration determines locality & scope of MEC service
  - ```POST .../mec_service_mgmt/v1/applications/{appInstanceId}/services```
- Discover MEC services:
  - Retrieve service list from MEC platform
  - Only discoverable & reachable services are returned
  - ```GET .../mec_service_mgmt/v1/services```
  - Example:
    - MEC Service instance running on MEC platform 1
    - Scope of locality set to _MEC\_HOST_ & _consumedLocalOnly_ set to _true_
    - Only MEC applications running on MEC platform 1 will be able to obtain MEC Service instance information
    - **NOTE:** MEC applications are always local in AdvantEDGE because each sandbox deploys a single MEC platform

