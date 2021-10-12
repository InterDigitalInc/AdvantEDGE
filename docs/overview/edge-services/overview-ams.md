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
- _Application self-controlled_
  - Application triggers and executes the context transfer
  - Context is transferred from source to target application
  - MEC system's role is to enable connectivity
- _Device assisted_
  - Device triggers and executes the context transfer
  - Context is kept on the device
  - MEC system's role is to decide if application mobility is required
- _MEC assisted_
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
- Application Mobility Service is implemented as a single sandbox pod within AdvantEDGE, providing service for all applications running as part of that sandbox
