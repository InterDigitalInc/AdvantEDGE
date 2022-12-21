---
layout: default
title: Location Service
parent: EDGE Services
grand_parent: Overview
nav_order: 1
permalink: docs/overview/edge-services/loc-service/
---

## Service Overview
AdvantEDGE provides a built-in Location Service implementation that integrates with scenarios.

This feature provides the following capabilities:
- _Learning location of a device within the network_
- _Learning information on all devices located within a zone or connected to a point-of-access_
- _Getting real-time updates on device location as they move across the network_

## Micro-Services
  - _Location Service:_ Implements ETSI MEC013 northbound API with a custom integration with AdvantEDGE APIs

## Northbound API
- Location Service is compliant with the ETSI MEC013 Specification, v2.2.1:
  - [ETSI GS MEC 013 V2.2.1](https://www.etsi.org/deliver/etsi_gs/MEC/001_099/013/02.02.01_60/gs_mec013v020201p.pdf)
  - [ETSI Forge - Location API repository](https://forge.etsi.org/gitlab/mec/gs013-location-api)
- MEC013 OMA references:
  - [RESTful Network API for Zonal Presence](http://www.openmobilealliance.org/release/REST_NetAPI_ZonalPresence/V1_0-20160308-C/OMA-TS-REST_NetAPI_ZonalPresence-V1_0-20160308-C.pdf)
  - [RESTful Network API for Terminal Location](https://www.openmobilealliance.org/release/TerminalLocationREST/V1_0_1-20151029-A/OMA-TS-REST_NetAPI_TerminalLocation-V1_0_1-20151029-A.pdf)
- API
  - [API Definition](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/docs/api-location)
  - Based on OpenAPI Specification (OAS) 3.0

## AdvantEDGE Integration
- Location service is implemented as a single sandbox pod within AdvantEDGE, providing service for all applications running as part of that sandbox

- 3 components:
  - Northbound Interface (NBI) & Core - (tightly coupled) implements the Location Service API & internal service logic
  - Southbound Interface (SBI) - (decoupled from NBI/Core) implements glue logic between AdvantEDGE & the NBI/Core

- Threads:
  - Main thread      - (NBI/Core) Handles requests to the Location Service API (server) from users (i.e. scenario pods)
  - NBI event thread - (NBI/Core) Handles event channel from SBI (for Location Service subscriptions)
  - SBI event thread - (SBI) Handle events from AdvantEDGE (scenario updates, mobility events, etc.) & updates Loc. Service database)

- Supports hot-restart
  - User / app subscriptions with and without a duration parameter survive location service pod restarts

Figure below presents an exemplary sequence diagram

![master]({{site.baseurl}}/assets/images/flow-mob-event.png)
