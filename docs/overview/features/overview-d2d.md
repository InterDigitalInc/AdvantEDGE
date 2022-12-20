---
layout: default
title: Device-to-device (D2D) Communication
parent: Features
grand_parent: Overview
nav_order: 9
permalink: docs/overview/features/d2d/
---

## Feature Overview
AdvantEDGE supports geospatial and network emulation for direct communication between UEs, or device-to-device (D2D) communication.

This feature provides the following capabilities:
- _Configurable D2D emulation_
  - D2D radius: maximum distance for D2D proximity
  - D2D via network: D2D via Uu (inter-UE communication via network)
  - D2D wireless type: D2D via PC5 (direct UE communication)
- _D2D Geospatial measurements_
  - Distance calculations between UEs that support D2D
- _D2D Network characteristics emulation_
  - Provides a "one-hop" data path between two UEs in D2D range
  - Abstracts D2D features such as keep-alive, advertisement, etc.
  - Allows UEs out of PoA range to use D2D link for network access
	
## Micro-Services
- _Traffic Control Engine:_ Enforces the real-time connectivity rules & network characteristics
- _Map server:_ Open Map Tiles is used to serve map data required by the frontend
- _GIS Engine:_ Determines D2D proximity for UEs
- _Databases:_ 
  - Postgres/Postgis backend database to store geospatial assets & perform calculations
  - Redis backend database to cache geospatial measurements

## Scenario Configuration

Element | Description
------ | ------
D2D Radius | Maximum distance between UEs to be considered in D2D range
Disable D2D via network	| Flag to enable/disable D2D via network for all UEs in the scenario
Supported Wireless Types | List of supported wireless access types; setting D2D enables D2D via PC5 emulation

## Scenario Runtime

Runtime | Description
------ | ------
GIS Engine | For UEs with D2D communication enabled, calculates distances between UEs and determines if they are in D2D range
Traffic Control Engine | Uses information about selected POA and D2D proximity to determine and enforce real-time connectivity rules & network characteristics

## Using D2D Feature
 
### Scenario Configuration
To enable D2D capabilities, D2D information must be configured in the scenario, either by updating an existing scenario or creating a new scenario from scratch.

As shown below, the configuration pane enables toggling of the visualization between the network hierarchy and map views (top left corner). Both views may be used to select and configure or edit geospatial data for the UEs.

The hierarchy view must be used to edit the scenario-level configuration. By clicking on the _Scenario_ network element (_Internet_ cloud in the diagram) the D2D radius can be set and D2D communication via the network may be disabled for all UEs in the scenario.
 
![D2DScenarioConfig]({{site.baseurl}}/assets/images/D2DScenarioConfig.png)

### Scenario Execution
At deployment time, a scenario containing D2D data will automatically populate the GIS Engine databases and trigger inter-UE distance calculations whenever a UE position changes.

From the execution tab, it is possible to visualize assets on a map by selecting the map view from the dashboard menu.

![D2DExecuteScenario.png]({{site.baseurl}}/assets/images/D2DExecuteScenario.png)

Once the scenario is successfully deployed, the GIS Engine will refresh the location of the UEs according to scenario evolution. When a UE position changes, the GIS Engine recalculates inter-UE distance to determine if UEs are within D2D proximity. The TC Engine is informed of D2D proximity and determines the required network connectivity rules and network characteristics to apply.

When in D2D range, applications on UEs benefit from 

When D2D-capable UEs are in range of each other and satify the D2D proximity requirements, they establish a direct connection. Applications deployed on such UEs benefit from a single hop link between UEs reflected in the observed network characteristics. Additionally, if a UE is not in range of a POA but is in D2D proximity with another UE, it will obtain network connectivity via its D2D connection.

![D2DAutomationConnected.png]({{site.baseurl}}/assets/images/D2DAutomationConnected.png)

When a D2D-capable UE is not in D2D proximity with another D2D-capable UE, the device remains in disconnected state.

![D2DAutomationDisconnected.png]({{site.baseurl}}/assets/images/D2DAutomationDisconnected.png)
