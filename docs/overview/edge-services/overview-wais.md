---
layout: default
title: Wireless Access Information Service
parent: EDGE Services
grand_parent: Overview
nav_order: 3
permalink: docs/overview/edge-services/wais/
---

## Service Overview
AdvantEDGE provides a built-in WAIS implementation that integrates with scenarios.

This service provides the following capabilities:
- _Learning terminal (station or STA) information_
- _Learning access-point information_
- _Getting real-time updates on devices WLAN conditions as they move across the network_

## Micro-Services
  - _WAIS:_ Implements ETSI MEC028 northbound API with a custom integration with AdvantEDGE APIs

## Northbound API
- WLAN Access Information Service is compliant with the ETSI MEC028 Specification, v2.1.1:
  - [ETSI GS MEC 028 V2.1.1](https://www.etsi.org/deliver/etsi_gs/MEC/001_099/028/02.01.01_60/gs_mec028v020101p.pdf)
  - [ETSI Forge - WLAN Access Information API repository](https://forge.etsi.org/rep/mec/gs028-wai-api)
- API
  - [API Definition](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/docs/api-wais)
  - Based on OpenAPI Specification (OAS) 3.0

## AdvantEDGE Integration
- WAIS is implemented as a single sandbox pod within AdvantEDGE, providing service for all applications running as part of that sandbox

- 3 components:
  - Northbound Interface (NBI) & Core - (tightly coupled) implements the WAIS API & internal service logic
  - Southbound Interface (SBI) - (decoupled from NBI/Core) implements glue logic between AdvantEDGE & the NBI/Core

  - Threads:
    - Main thread      - (NBI/Core) Handles requests to the WAIS API (server) from users (i.e. scenario pods)
    - NBI event thread - (NBI/Core) Handles event channel from SBI (for WAIS subscriptions)
    - SBI event thread - (SBI) Handle events from AdvantEDGE (scenario updates, mobility events, etc.) & updates Loc. Service database)

  - Supports hot-restart
    - User / app subscriptions with and without a duration parameter survive WAIS pod restarts
