---
layout: default
title: Application Mobility Service
parent: EDGE Services
grand_parent: Overview
nav_order: 5
permalink: docs/overview/edge-services/ams/
---

## Service Overview
AdvantEDGE provides a built-in Application Mobility Service implementation that integrates with scenarios.

Application Mobility Service provides support for relocation of user context between MEC hosts; application instance relocation is not supported.

AMS defines three types of MEC application user-context transfer:
- _Application self-controlled_ (not supported)
  - Application triggers and executes the context transfer
  - Context is transferred from source to target application
  - MEC system's role is to enable connectivity
- _Device assisted_ (not supported)
  - Device triggers and executes the context transfer
  - Context is kept on the device
  - MEC system's role is to decide if application mobility is required
- _MEC assisted_ (supported)
  - MEC system triggers and assists the context transfer
  - Context is transferred from source to target application

## Micro-Services
  - _AMS:_ Implements ETSI MEC021 northbound APIs with a custom integration with AdvantEDGE APIs

## API Version
- Application Mobility Service is compliant with the ETSI MEC021 Specification, v2.1.1:
  - [ETSI GS MEC 021 V2.1.1](https://www.etsi.org/deliver/etsi_gs/MEC/001_099/021/02.01.01_60/gs_MEC021v020101p.pdf)
  - [ETSI Forge - Application Mobility Service API](https://forge.etsi.org/rep/mec/gs021-amsi-api)
- API
  - [API Definition](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/docs/api-ams)
  - Based on OpenAPI Specification (OAS) 3.0

## AdvantEDGE Integration
Application Mobility Service is implemented as a single sandbox pod within AdvantEDGE, providing service for all applications running as part of that sandbox.

To use this service, MEC applications must first obtain an application instance ID from the AdvantEDGE platform using one of the following methods:
- For MEC applications configured in the deployed scenario:
  - From the ```MEEP_APP_ID``` environment variable, or
  - From the process table in the frontend Execution page
- For external MEC applications not configured in the deployed scenario:
  - From the _Applications_ endpoints of the [Sandbox Controller API]({{site.baseurl}}{% link docs/overview/overview-api.md %}#sandbox-api)

After obtaining the application instance ID, MEC applications must use the [Application Support API]({{site.baseurl}}{% link docs/overview/overview-api.md %}#app-support-api) to confirm application readiness and optionally register for graceful termination. Once confirmed, MEC Applications may use the Application Mobility Service to perform MEC-assisted application mobility as described in the use case below.

Note that MEC application instances must be configured to run on an edge or fog node; this can be provisioned directly in the configured scenario by placing a MEC Application under an edge/fog node, or via the mandatory _nodeName_ parameter of the _Sandbox Controller Applications_ endpoints.

The Application Mobility Service implements a minimum hop-count MEC Application selection algorithm to determine which of the registered MEC application instances to use. The algorithm uses real-time terminal device locations to calculate the nearest registered MEC Application instance; it then triggers application context transfers when necessary.

### Use case for MEC-assisted application mobility
This use case applies to MEC Applications with multiple instances that wish to perform user context transfers with the assistance of the Application Mobility Service (AMS). MEC Applications register to AMS with a list of terminal devices to track, and subscribe for Mobility Procedure notifications to receive target application instance details on terminal mobility events.
- Configure & Deploy an AMS scenario: 
  - Must include terminal devices that move across POAs & edge/fog node coverage areas
  - Must include multiple MEC Application instances running on different edge/fog nodes
  - Each instance must obtain a unique application instance ID & confirm application readiness
- Register to AMS:
  - MEC Applications provide registration information (including application instance ID)
  - Device information should only be included by the MEC application instance with user context ownership
  - ```POST .../amsi/v1/app_mobility_services```
- Subscribe for Mobility notifications:
  - MEC applications with user context ownership must subscribe for _MobilityProcedure_ notifications
  - Tracked device information must also be provided in the subscription request
  - ```POST .../amsi/v1/subscriptions```
- Wait for Terminal device to trigger Mobility notification:
  - Enable terminal device mobility & movement
  - AMS algorithm monitors terminal locations and triggers _MobilityProcedure_ notifications when there is a change in the preferred/selected MEC application instance
  - Notification includes target application information
  - _NOTE:_ notifications are only send to MEC applications running on the source edge/fog node
- Perform user context transfer
  - Source MEC application instance transfers terminal device context to target MEC application instance
- Inform AMS about transfer complete:
  - The source MEC application should set the device context transfer state to _COMPLETE_
  - Optionally delete the device information until the user context returns
  - ```PUT .../app_mobility_services/{appMobilityServiceId}```
- Repeat procedure:
  - Target MEC Application now becomes the source for future context transfers
