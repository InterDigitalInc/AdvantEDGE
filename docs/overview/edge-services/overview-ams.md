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

Application Mobility Service provides support for relocation of user context between MEC hosts; application instance relocation not supported.

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

To use this service, MEC applications must first obtain an application instance ID from the AdvantEDGE platform using the _Applications_ endpoints of the [Sandbox Controller API]({{site.baseurl}}{% link docs/overview/overview-api.md %}#sandbox-api). After provisioning the application instance ID, MEC applications must use the [Application Support API]({{site.baseurl}}{% link docs/overview/overview-api.md %}#app-support-api) to confirm application readiness and optionally register for graceful termination.

Once confirmed, MEC Applications may use the Application Mobility Service to perform MEC-assisted application mobility as described in the use case below.

### Use case for MEC-assisted application mobility
This use case applies to MEC Applications with multiple instances running on different MEC platforms wishing to perform user context transfers with the assistance of the Application Mobility Service. MEC applications inform the service about which terminal devices to track and subscribe for Mobility Procedure notifications with details about the target application instance.
- Start multiple MEC Application instances
  - Each instance must obtain a unique application instance ID & confirm application readiness
- Register to AMS:
  - MEC Applications provide registration information (including application instance ID)
  - Device information must be included only by the MEC application with the user context
  - ```POST .../amsi/v1/app_mobility_services```
- Subscribe for Mobility notifications:
  - MEC applications with the user context must subscribe for _MobilityProcedure_ notifications
  - Tracked device information must also be provided
  - ```POST .../amsi/v1/app_mobility_services```
- Wait for Terminal device to trigger Mobility notification:
  - AMS algorithm monitors terminal and trigger a mobility procedure when it transitions to the coverage area of a different MEC platform
  - Notification includes target application information
  - _NOTE:_ notifications are only send to MEC applications running on the source MEC platform
- Perform user context transfer
  - Source MEC application instance transfers terminal device context to target MEC application instance
- Inform AMS about transfer complete:
  - The source MEC application should set the device context transfer state to _COMPLETE_
  - Optionally delete the device information until the user context returns
  - PUT .../app_mobility_services/{appMobilityServiceId}
- Repeat procedure:
  - Target MEC Application now becomes the source for future context transfers
