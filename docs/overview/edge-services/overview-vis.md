---
layout: default
title: V2X Information Service
parent: EDGE Services
grand_parent: Overview
nav_order: 1
permalink: docs/overview/edge-services/vis/
---

## Service Overview
AdvantEDGE provides a built-in V2X Information Service implementation that integrates with scenarios.

This feature provides the following capabilities:
- _TBD_
- _TBD_
- _TBD_

## Micro-Services
  - _V2X Information Service:_ Implements ETSI MEC030 northbound API with a custom integration with AdvantEDGE APIs

## Northbound API
- V2X Information Service is compliant with the ETSI MEC030 Specification, v2.2.1:
  - [ETSI GS MEC 030 V2.2.1](https://www.etsi.org/deliver/etsi_gs/MEC/001_099/030/02.02.01_60/gs_MEC030v020201p.pdf)
  - [ETSI Forge - V2X Information Service API repository](https://forge.etsi.org/rep/mec/gs030-vis-api)
- API
  - [API Definition](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/docs/api-location)
  - Based on OpenAPI Specification (OAS) 3.0

## AdvantEDGE Integration
- V2X Information service is implemented as a single sandbox pod within AdvantEDGE, providing service for all applications running as part of that sandbox

- 3 components:
  - Northbound Interface (NBI) & Core - (tightly coupled) implements the V2X Information Service API & internal service logic
  - Southbound Interface (SBI) - (decoupled from NBI/Core) implements glue logic between AdvantEDGE & the NBI/Core

- Threads:
  - Main thread      - (NBI/Core) Handles requests to the V2X Information Service API (server) from users (i.e. scenario pods)
  - NBI event thread - (NBI/Core) Handles event channel from SBI (for V2X Information Service subscriptions)
  - SBI event thread - (SBI) Handle events from AdvantEDGE (scenario updates, mobility events, etc.) & updates V2X Information Service database)

- Supports hot-restart
  - User / app subscriptions with and without a duration parameter survive location service pod restarts

Figure below presents an exemplary sequence diagram

