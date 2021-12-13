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
- WLAN Access Information Service is compliant with the ETSI MEC028 Specification, v2.2.1:
  - [ETSI GS MEC 028 V2.2.1](https://www.etsi.org/deliver/etsi_gs/MEC/001_099/028/02.02.01_60/gs_MEC028v020201p.pdf)
  - [ETSI Forge - WLAN Access Information API repository](https://forge.etsi.org/rep/mec/gs028-wai-api)
- API
  - [API Definition](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/docs/api-wais)
  - Based on OpenAPI Specification (OAS) 3.0

## AdvantEDGE Integration
WAIS is implemented as a single sandbox pod within AdvantEDGE, providing service for all applications running as part of that sandbox

WAIS components:
- Northbound Interface (NBI) & Core - (tightly coupled) implements the WAIS API & internal service logic
- Southbound Interface (SBI) - (decoupled from NBI/Core) implements glue logic between AdvantEDGE & the NBI/Core
- Threads:
  - _Main thread:_ (NBI/Core) Handles requests to the WAIS API (server) from users (i.e. scenario pods)
  - _NBI event thread:_ (NBI/Core) Handles event channel from SBI (for WAIS subscriptions)
  - _SBI event thread:_ (SBI) Handle events from AdvantEDGE (scenario updates, mobility events, etc.) & updates Loc. Service database)
- Supports hot-restart
  - User / app subscriptions with and without a duration parameter survive WAIS pod restarts

### WebSocket Support
WAIS implements the [ETSI MEC009](https://www.etsi.org/deliver/etsi_gs/MEC/001_099/009/03.01.01_60/gs_MEC009v030101p.pdf) _REST-based subscribe/notify with Websocket
fallback_ pattern.

This pattern requires sending a test notification to the notification endpoint URI in the subscription request to determine
if it is directly reachable or not. On success, standard REST-based notification methods continue; on failure (and upon client request),
WAIS creates a WebSocket subscription endpoint where the client may connect to receive notifications.

_**NOTE:**_ WebSockets are required when client applications subscribing for notifications are behind a NAT/firewall.

