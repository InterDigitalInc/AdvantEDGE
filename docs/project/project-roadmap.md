---
layout: default
title: Project Roadmap
parent: Project
nav_order: 2
---

- [Project Progression](#project-progression)
- [Project Roadmap](#project-roadmap)

## Project Progression

| Timeline | Version | Description |
|---|---|---|
| 07-2018 |  | Project start
| 09-2018 |  | Demo: Single-zone Drones Edge-DAA @ Edge Congress (Berlin)<br>_DAA = Detect-And-Avoid_
| 02-2019 |  | Demo: Multi-zone Drones Edge-DAA @ MWC (Barcelona)<br>_DAA = Detect-And-Avoid_
| 05-2019 | [1.0.0](https://github.com/InterDigitalInc/AdvantEDGE/releases/tag/v1.0.0) | Limited release to select partners |
| 07-2019 | [1.1.0](https://github.com/InterDigitalInc/AdvantEDGE/releases/tag/v1.1.0) | Feature release |
| 09-2019 | [1.2.0](https://github.com/InterDigitalInc/AdvantEDGE/releases/tag/v1.2.0) | - Apache 2.0 Open Source Release<br>- Demo: OpenRTiST real-time video style transfer @ Edge Congress (London)
| 11-2019 | [1.3.0](https://github.com/InterDigitalInc/AdvantEDGE/releases/tag/v1.3.0) | Feature release
| 03-2020 | [1.4.0](https://github.com/InterDigitalInc/AdvantEDGE/releases/tag/v1.4.0) | - Feature release<br>- Demo: In-Home Edge Gaming @ MWC (Barcelona)<br>_Demo Canceled due to COVA19 Virus_ :(
| 06-2020 | [1.5.0](https://github.com/InterDigitalInc/AdvantEDGE/releases/tag/v1.5.0) | - Feature release<br>- GIS support<br>- Multi-user sandboxes<br>- ETSI MEC 012 Radio Network Interface Service (RNIS)<br>- Ansible playbooks install (beta)<br>- containerized virt-engine<br>- Net.Char.: Asymmetric max throughput + selectable jitter distribution<br>- ingress controller (port 80/443)
| 12-2020 |[1.6.0](https://github.com/InterDigitalInc/AdvantEDGE/releases/tag/v1.6.0)|- New MEC Service: ETSI-MEC WLAN Access Information Service (WAIS)<br>- ETSI MEC RNIS & Location Service graduate to v2.1.1<br>- Support different PoA radio technologies (5G/4G/WLAN), UE radio technologies & UE disconnected state<br>- Possibility to limit CPU and Memory on scenario pods<br>- Network characteristics automaion based on geolocation<br>- Zone color coding on map + other map improvements<br>- Platform API and GIS engine API improvements<br>- Ansible playbook for dev. environment<br>- Support for Ubuntu 20.04/K8s 1.19 & transition to Helm v3<br>- Support for Open API v3.0 specifications<br>- Various deployment knobs added (Let's Encrypt CA certs, max users, session management, user roles, OAuth) to support diverse deployments of the platform<br>- Various bug fixes
|01-2021 | [1.6.1](https://github.com/InterDigitalInc/AdvantEDGE/releases/tag/v1.6.1)| - Helm related hotfix
|04-2021 | [1.7.0](https://github.com/InterDigitalInc/AdvantEDGE/releases/tag/v1.7.0)| - Process lifecycle API to control containers @ scenario runtime<br>- Network reachability/connectivity to control reachability of the network @ scenario runtime<br>- Improved platform monitoring & KPIs (Prometheus)<br>- New dashboards<br>- RNIS improvements (L2Meas, measurement reports) 
|09-2021 | [1.8.0](https://github.com/InterDigitalInc/AdvantEDGE/releases/tag/v1.8.0)| - New MEC Service: ETSI-MEC Edge Application Enablement Service<br>- New MEC Service: ETSI-MEC Application Mobility Service<br>- Geopspatial measurements caching (distance & signal strength)<br>- Wireless Metrics dashboards in frontend<br>- Long-term metrics storage capabilities
|12-2021 | [1.8.1](https://github.com/InterDigitalInc/AdvantEDGE/releases/tag/v1.8.1)| - WebSocket subscription support in ETSI MEC WLAN Access Information Service<br>- MEC Service maintenance updates (Edge Application Enablement & Application Mobility Services)<br>- Sandbox metrics garbage collection<br>- Plaform maintenance fixes
|04-2022 | 1.9.0| TBD

## Project Roadmap

The following features are considered & prioritized in "as-needed" basis

| Feature | Description |
| --- | --- |
| Value added KPIs | Platform provides new KPIs to allow application prototyping & experimentation |
| Scenario Validation | Platform provides scenario validation capabilities |
| New ETSI MEC Services | Platform supports new MEC services; BWM (MEC015) as the primary target |
| Mobile Edge/Fog Nodes | Platform allows mobility of edge nodes with creation and destruction of containers |
| Model extension | Platform allows to model new network topologies and types |
| Improved VM support | Provide guidance on installing AdvantEDGE in VM environment |
| Multi-user isolation | Isolate user sandboxes |
| Metrics recording trigger | Start recording metrics on a user generated events |
| WebSockets | Platform supports subscriptions via WebSockets; WAIS (MEC028) as the primary target |
| PoA Selection Algorithm | Improve algorithm to consider network KPIs |
| Edge/Fog Selection Algorithm | Improve "minimum hop-count" algorithm to consider network KPIs |
