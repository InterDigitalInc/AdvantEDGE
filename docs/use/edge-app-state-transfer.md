# Edge Application State Transfer

## Overview
AdvantEDGE provides a state transfer service that facilitates UE state transfer between instances of a multi-edge group.

To use the state transfer service, multi-edge application instances must:
1. Create & Configure a Mobility Group (MG) using the MG Manager service REST API
> *Note: The MG is automatically created by the AdvantEDGE platform at scenario deployment time, based on the multi-edge group defined in the scenario; therefore, there is no need to create the MG via the MG Manager API.*
2. Register to the MG using the MG Manager service REST API
3. Implement the MG Application State Transfer REST API to handle MG application state events

When multi-edge application instances are registered to the MG, the MG Manager informs them when a UE state transfer is needed. They can then transfer the requested UE state to the target application instance(s) via the MG Manager service API.

> *Note: The MG Manager uses the MG configuration to determine when to send events and which events to send to the application instances.*

The Edge Application State Transfer & MEEP Controller REST API definitions can be found here:
- [MEEP Mobility Group Manager REST API](../../go-apps/meep-mg-manager/api/swagger.yaml)
- [MEEP Mobility Group Application REST API](../../go-apps/meep-mg-manager/client-app-api/meep-mg-app-api-client.yaml)
- [MEEP Controller REST API](../../go-apps/meep-ctrl-engine/api/swagger.yaml)

## Mobility Group Creation
Mobility Groups are automatically created at scenario deployment time based on the multi-edge groups defined in the scenario.

The MG Manager service REST API endpoints used to create and configure MGs are the following:
- /mg
- /mg/{mgName}
> *Note: The MG Manager currently supports only the default configuration. DO NOT use this API to set a different configuration.*

The AdvantEDGE platform uses the following default MG configuration:

MG Config Param        | Default Value    | Description
-----------------------|------------------|------------
stateTransferMode      | STATE-MANAGED    | MG Manager is used to transfer UE state
stateTransferTrigger   | NET-LOC-IN-RANGE | The following state transfer events are sent to trigger UE state transfer:<ul><li>STATE-TRANSFER-START:<ul><li>Sent when UE is in range of a POA with a different edge service instance</li><li>Triggered by a POAS-IN-RANGE event sent to MEEP Controller</li></ul> <li>STATE-TRANSFER-COMPLETE: <ul><li>Sent when UE moves to a POA with a different edge service instance</li><li>Triggered by a UE-MOBILITY event sent to MEEP Controller</li></ul> <li>STATE-TRANSFER-CANCEL:<ul><li>Sent when UE is no longer in range of a POA with a different edge service</li><li>Triggered by a POAS-IN-RANGE event sent to MEEP Controller</li></ul> </li></ul>
sessionTransferMode    | FORCED           | Session between UE & Edge App instance is forcefully dropped after STATE-TRANSFER-COMPLETE event
loadBalancingAlgorithm | HOP-COUNT        | Target multi-edge app instance is determined using minimum hop count to reach service

## Multi-Edge Application Instance Registration
The MG Manager service REST API endpoints used for multi-edge application instance registration are the following:
- /mg/{mgName}/app
- /mg/{mgName}/app/{appId}

The following figure presents the MG creation and application instance registration procedure.

![me-app-deploy.](../images/edge-app-state-transfer-deployment.png)

1. On scenario deployment, MG Manager service creates the Mobility Group(s)
2. Ctrl Engine (via the Virt Engine) deploys the multi-edge application instances
3. Upon start-up, multi-edge application instances register to the MG Manager
    - POST /mg/multi-edge-svc/app/multi-edge-app1
    - POST /mg/multi-edge-svc/app/multi-edge-app2

## UE State Transfer on UE Mobility Event
The MG Manager service REST API endpoint used for transferring edge application UE state is the following:
- /mg/{mgName}/app/{appId}/state

The MG Application REST API endpoint used for handling UE state transfer events is the following:
- /mg/event

The MEEP Controller REST API endpoint used for injecting UE Mobility events is the following:
- /events

The following figure presents the UE State Transfer procedure.

![me-app-runtime.](../images/edge-app-state-transfer-runtime.png)

1. At connection time, UE traffic is routed to the closest multi-edge app instance (multi-edge-app1)
2. A UE Mobility event is sent to the MEEP Ctrl Engine to simulate a change of POA
    - POST /events {UE-MOBILITY}
3. MG Manager processes change of POA and determines that UE is now closest to a new multi-edge app instance (multi-edge-app2)
4. MG Manager sends a state transfer event on the initial multi-edge app instance (multi-edge-app1)
    - POST /mg/event {STATE-TRANSFER-COMPLETE}
5. Initial multi-edge app instance (multi-edge-app1) sends its latest UE state to the MG Manager
    - POST /mg/multi-edge-svc/app/multi-edge-app1/state {UE, state}
6. MG Manager sends the UE state to the new multi-edge app instance (multi-edge-app2)
    - POST /mg/event {STATE-UPDATE}
7. MG Manager forcefully terminates the UE connection to the initial multi-edge app instance (multi-edge-app1)
    - The UE should create a new connection to the multi-edge service, which will be routed to the new multi-edge app instance (multi-edge-app2)

## UE State Transfer on POAs in Range Event

> *Note: In cases where UE state transfer is required before a UE Mobility event occurs, the MG Manager service may be used to trigger a state transfer on a MEEP Controller POAS-IN-RANGE event. This is the default MG Manager behavior, however the Demo Applications provided with the AdvantEDGE platform do not use this functionality. More details will be provided at a later time or upon request.*
