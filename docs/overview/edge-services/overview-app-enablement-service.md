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
- Edge Platform Application Enablement Service is implemented as a single sandbox pod within AdvantEDGE, providing service for all applications running as part of that sandbox
