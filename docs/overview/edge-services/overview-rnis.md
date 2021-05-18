---
layout: default
title: Radio Network Information Service
parent: EDGE Services
grand_parent: Overview
nav_order: 2
permalink: docs/overview/edge-services/rnis/
---

## Service Overview
AdvantEDGE provides a built-in RNIS implementation that integrates with scenarios.

This service provides the following capabilities:
- _Learning radio network conditions_
- _Performing user plane measurements_
- _Learning about devices connected to the radio node(s) associated with the mobile edge host and their radio access bearers_
- _Getting real-time updates on devices radio conditions as they move across the network_

### Micro-Services
  - _RNIS:_ Implements ETSI MEC012 northbound API with a custom integration with AdvantEDGE APIs

## API Version
- Radio Network Information Service is compliant with the ETSI MEC012 Specification, v2.1.1:
  - [ETSI GS MEC 012 V2.1.1](https://www.etsi.org/deliver/etsi_gs/MEC/001_099/012/02.01.01_60/gs_mec012v020101p.pdf)
  - [ETSI Forge - Radio Network Information API repository](https://forge.etsi.org/rep/mec/gs012-rnis-api)
- API
  - [API Definition](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/docs/api-rnis)
  - Based on OpenAPI Specification (OAS) 3.0

### AdvantEDGE Integration
- RNIS is implemented as a single sandbox pod within AdvantEDGE, providing service for all applications running as part of that sandbox

- 3 components:
  - Northbound Interface (NBI) & Core - (tightly coupled) implements the RNIS API & internal service logic
  - Southbound Interface (SBI) - (decoupled from NBI/Core) implements glue logic between AdvantEDGE & the NBI/Core

  - Threads:
    - Main thread      - (NBI/Core) Handles requests to the RNIS API (server) from users (i.e. scenario pods)
    - NBI event thread - (NBI/Core) Handles event channel from SBI (for RNIS subscriptions)
    - SBI event thread - (SBI) Handle events from AdvantEDGE (scenario updates, mobility events, etc.) & updates Loc. Service database)

  - Supports hot-restart
    - User / app subscriptions with and without a duration parameter survive RNIS pod restarts
