---
layout: default
title: Execution View
parent: GUI Overview
grand_parent: Usage
nav_order: 2
permalink: docs/usage/gui/exec-view/
---

Topic | Abstract
------|------
[Execution View](#execution-view) | Execution View overview
[Headline bar](#headline-bar) | Execution View controls
[Event configuration bar](#event-configuration-bar) | Event pane selection & replay file creation controls
[Dashboard configuration bar](#dashboard-configuration-bar) | Visualization dashboard controls
[Visualization dashboard](#visualization-dashboard) | Network topology, map & metrics dashboard visualization
[Manual event pane](#manual-event-pane) | Manual event controls
[Automation pane](#automation-pane) | Automated event controls
[Replay event pane](#replay-event-pane) | Event replay execution controls
[Process status table](#process-status-table) | Application execution status

---
## Execution View
The Execution view provides a graphical interface to perform scenario run-time operations. You can use this view to deploy and terminate scenarios in the backend. You can also inject events to change the network topology or link characteristics.

Other views can be found in the [GUI Overview wiki page]({{site.baseurl}}{% link docs/usage/usage-gui-overview.md %})

---
## Headline Bar
This bar is always visible within the Execution view. It provides the following buttons:

### Sandbox
Drop-down menu for selecting a sandbox to use:
- Menu list includes all existing sandboxes
- Selecting a sandbox updates the Execution View with the current sandbox state

### NEW
Create a new sandbox:
- Opens a dialog prompt for a new sandbox name
  - Validates the sandbox name format & length
- Creates a new sandbox
  - Creates & starts new sandbox pods
  - Enables scenario deployment buttons

_**NOTE:** Operation fails silently if a sandbox with the same name already exists. In this case you must click on **NEW** again and enter a unique sandbox name._

### DELETE
Delete the selected sandbox:
- Opens a dialog to confirm delete action
- Deletes the sandbox
  - Removes the sandbox pods
  - Disables scenario deployment buttons until a new sandbox is selected or created

_**NOTE:** Scenario pods running in the sandbox will also be removed when deleting a sandbox._

### DEPLOY
Activate a scenario:
- Opens a dialog prompt to select a scenario to deploy
  - Scenario list is retrieve from the backend persistent store
- Sends the chosen scenario name to the backend for deployment
  - Opens the Visualization dashboard with the selected views
  - Shows the Process status table with the scenario pod status

### SAVE
Save deployed scenario:
- Opens a dialog prompt for the scenario name to be saved
  - Validates the scenario name format & length
- Sends the current state of the deployed scenario to the backend where it is added to the persistent store

_**NOTE:** This operation overwrites any existing scenario with the same name. Be careful to use a unique name when saving the deployed scenario._

### TERMINATE
Deactivate a deployed scenario:
- Opens a dialog prompt to confirm deployed scenario termination
- Sends the termination request to the backend
  - Closes the Visualization dashboard
  - Closes the Process status table only after all scenario pods have terminated

### EVENT
Opens the event configuration bar.

### DASHBOARD
Opens the dashboard configuration bar.

---
## Event Configuration Bar
This section provides controls for event pane selection and replay file creation. It also displays the status of a running replay file. The following controls are available:

### MANUAL
Opens the Manual event pane for event configuration.

### AUTOMATION
Opens the Automation pane for automation configuration.

### AUTO-REPLAY
Opens the Auto-replay event pane for event configuration.

### SAVE EVENTS
Create a new replay file:
- Opens a dialog prompt for the replay file name & description to be created
- Sends a replay file generation request to the backend
  - Backend queries the event database and creates a new replay file with the provided name & description

_**NOTE:** You must click on the **AUTO-REPLAY** button for the new replay file to be added to the list of available replay files in the frontend._

### CLOSE
Closes the event configuration bar.

---
## Dashboard Configuration Bar
This section provides the following configuration controls for the Visualization dashboard:

### VIEW 1 & 2
Drop-down menus to select the visualizations to display in the Visualization dashboard. You may configure a single view or two views to display side-by-side. You can choose from the following views:
- Network Topology
- Map View
- Sequence Diagram
- Data Flow Diagram
- Network Metrics Point-to-Point
- Network Metrics Aggregation
- Wireless Metrics Point-to-Point
- Wireless Metrics Aggregation
- None

---
## Visualization Dashboard
This section is only visible when a scenario is deployed. It displays the views selected in the Dashboard configuration bar. Each view is described below:

### Network Topology
This view provides a visual representation of the network element hierarchy for the deployed scenario. It is updated whenever a scenario is deployed, refreshed or terminated. It shows basic network element information such as type, name & network characteristics directly in the graph.

We use a 3rd party graphing tool called [vis.js](https://visjs.org/) to create a network view of the scenario elements. You can drag and zoom the network view using the provided controls or using your mouse click & scroll buttons. Hovering over a node or link in the graph will open an informative tooltip with details about the currently configured element state.

#### Dashboard Options:
#### Show Apps
Checkbox to show or hide applications in the Visualization dashboard views. If checked, the Network Topology view includes edge, UE & cloud application network elements.

#### Close
Closes the dashboard configuration bar.

### Map View
This view provides a visual representation of the network element geographic positions on a map for the deployed scenario. It is updated whenever a scenario is deployed, refreshed or terminated.

We use a 3rd party mapping client called [Leaflet](https://leafletjs.com/) to render map tiles. Physical location assets such as UEs, PoAs & Compute nodes are then overlaid on the maps using markers, circles and lines. You can drag and zoom the map view using the provided controls or using your mouse click & scroll buttons. Hovering over a marker displays its name, and clicking on a marker opens a tooltip with the asset name, position, range, end-of-path mode & velocity, when applicable.

#### Dashboard Options:
#### Close
Closes the dashboard configuration bar.

### Sequence Diagram
This view provides a visual representation of the messages exchanged between components for the deployed scenario. Applications may use the _meep-http-logger_ package to log http messages with a source and destination identifier. The http metrics must be logged to the active scenario database by first invoking the _ReInit()_ function of the _meep-http-logger_ package and providing the current sandbox name and the active scenario name. Subsequent logging calls then create and store the message logs in the scenario database. The sequence diagram visualization retrieves and displays the logged messages.

To display the messages, we use a 3rd party sequence diagram generation tool called [Mermaid](https://mermaid.js.org/). Messages can also be exported to a text file formatted either for [Mermaid](https://mermaid.live/edit) or [SequenceDiagram.org](https://sequencediagram.org/).

#### Dashboard Options:
#### Participants
Text field to input the ordered set of participants to be rendered in the sequence digaram.

#### Count
Text field to input the maximum number of logs to be shown in the sequence diagram.

#### Clear
Button to clear the rendered diagram.

#### Fetch All
Button to fetch all the logs/events for the given scenario and render the diagram on screen.

#### Export Mermaid
Button to export the sequence diagram logs to a text file in mermaid format.

#### Export SDorg
Button to export the sequence diagram logs to a text file in sdorg (SequenceDiagram.org) format.

#### Pause
Checkbox to pause or unpause the automatic rendering of the updated sequence diagram.

#### Close
Closes the dashboard configuration bar.

### Data Flow Diagram
This view provides a visual representation of the number of component interactions for the deployed scenario. Applications may use the _meep-http-logger_ package to log http messages with a source and destination identifier. The http metrics must be logged to the active scenario database by first invoking the _ReInit()_ function of the _meep-http-logger_ package and providing the current sandbox name and the active scenario name. Subsequent logging calls then create and store the message logs in the scenario database. The data flow diagram visualization retrieves and displays the number of component interactions.

To display the number of component interactions, we use a 3rd party flow diagram generation tool called [Mermaid](https://mermaid.js.org/). Component interactions can also be exported to a text file formatted for [Mermaid](https://mermaid.live/edit).

#### Dashboard Options:
#### Export Mermaid
Button to export the component interaction logs in mermaid format.

#### Pause
Checkbox to pause or unpause the automatic rendering of the updated data flow diagram.

#### Close
Closes the dashboard configuration bar.

### Network Metrics Point-to-Point
This view loads the Network Metrics Point-to-Point Grafana Dashboard in an iframe. This dashboard monitors the deployed scenario and provides some latency and traffic metrics between the configured source and destination nodes. It displays instantaneous measurements for round-trip ping time and UL & DL throughput, as well as a graph of these measurements over the last minute.

_**NOTE:** You must select both a source and destination node in the Dashboard configuration bar to visualize data in this view._

#### Dashboard Options:
#### Source Node
Drop-down menu to select the source application node that must be provided to the Visualization dashboard.

#### Destination Node
Drop-down menu to select the destination application node that must be provided to the Visualization dashboard.

##### Close
Closes the dashboard configuration bar.

### Network Metrics Aggregation
This view loads the Network Metrics Aggregation Grafana Dashboard in an iframe. This dashboard monitors the deployed scenario and provides some latency and traffic metrics between the configured source and all other scenario application nodes. It displays graphs for round-trip ping times, UL & DL throughput measurements over the last minute. It also shows a table of events received by the backend.

_**NOTE:** You must select a source node in the Dashboard configuration bar to visualize data in this view._

#### Dashboard Options:
#### Source Node
Drop-down menu to select the source application node that must be provided to the Visualization dashboard.

#### Close
Closes the dashboard configuration bar.

### Wireless Metrics Point-to-Point
This view loads the Wireless Metrics Point-to-Point Grafana Dashboard in an iframe. This dashboard monitors the deployed scenario and provides some signal strength metrics and distance measurements between the configured UE and POA. It displays instantaneous measurements for RSRP & RSRQ when attached to a 4G or 5G POA, and the RSSI when attached to a WiFi POA. The dashboard also provides the current distance between the UE and POA, as well as a graph of all of these measurements over the last minute (overlayed with mobility events to easily track POA transitions).

#### Dashboard Options:
#### UE
Drop-down menu to select the UE node that must be provided to the Visualization dashboard.

#### POA
Drop-down menu to select the POA node that must be provided to the Visualization dashboard.

#### Close
Closes the dashboard configuration bar.

_**NOTE:** You must select both a UE and a POA in the Dashboard configuration bar to visualize data in this view._

### Wireless Metrics Aggregation
This view loads the Wireless Metrics Aggregation Grafana Dashboard in an iframe. This dashboard monitors the deployed scenario and provides some signal strength metrics and distance measurements between the configured UE and all POAs in range. It displays graphs for RSRP & RSRP (4G/5G POAs), RSSI (WiFi POAs), and distance measurements over the last minute (overlayed with mobility events to easily track POA transitions). It also shows a table of events received by the backend.

_**NOTE:** You must select a UE in the Dashboard configuration bar to visualize data in this view._

#### Dashboard Options:
#### UE
Drop-down menu to select the UE node that must be provided to the Visualization dashboard.

#### Close
Closes the dashboard configuration bar.

### None
Disables the view, giving its frontend real-estate to the other view.

---
## Manual Event Pane
This pane is used to manually configure and send events to the backend in order to dynamically alter the currently deployed scenario. Once fully configured you may send the event using the _SUBMIT_ button or cancel the event using the _CLOSE_ button. The event response status is displayed below the _SUBMIT_ button after sending an event.

You may configure the following events:

### MOBILITY
This event moves a chosen network element to a new location in the network topology. The following options must be configured:

##### Target
Drop-down menu to select the network element to be moved.

##### Destination
Drop-down menu to select the parent network element where the target element should be moved. This list is filtered to display valid destinations according to the selected target element.

When destination is set to _DISCONNECTED_, 100% packet loss is set on the selected target element to simulate network connectivity loss.

### NETWORK-CHARACTERISTICS-UPDATE
This event provides network characteristics that must be applied to a specific network element. The following options must be configured:

##### Network Element Type
Drop-down menu to select the type of network element to be updated.

##### Network Element
Drop-down menu to select the network element to be updated.

##### Latency / Latency Variation / Packet Loss / Throughput
Number fields to input element network characteristics to be sent in the event. Network characteristics are described in more detail [here]({{site.baseurl}}{% link docs/overview/overview-architecture.md %}#network-characteristics).

### SCENARIO-UPDATE
This event dynamically adds or removes terminal, edge or cloud applications in the network topology. This event enables the [Process Lifecycle feature]({{site.baseurl}}{% link docs/overview/features/overview-process-lifecycle.md %}).

The following options must be configured:

##### Action Type
Drop-down menu to select the action to perform.

_**Adding an element**_
When adding an application, element configuration fields must be filled before submitting the event. Frontend & backend validation is performed on the event to prevent errors during process instantiation. Once successfully submitted, the platform performs all necessary actions to start the new application.

_**Removing an element**_
When removing an application, the target _Process Type_ & _Process Name_ must be selected from the drop-down menus. Once successfully submitted, the platform performs all necessary actions to stop the selected application.

### PDU-SESSION
This event adds or removes terminal PDU Sessions. This event enables the [Cellular Connectivity Control feature]({{site.baseurl}}{% link docs/overview/features/overview-cell-connectivity-control.md %}).

The following options must be configured:

##### Action Type
Drop-down menu to select the action to perform.

##### Terminal
Drop-down menu to select the terminal.

##### PDU Session ID
Text Field to input PDU Session ID.

##### Data Network Name
Drop-down menu to select the data network name (only visible when creating a new PDU Session).

---
## Automation Pane
This pane is used to enable/disable GIS automations. More details on GIS automations can be found [here]({{site.baseurl}}{% link docs/overview/features/overview-gis.md %}).

### Movement
Toggles UE movement according to their respective paths.

### Mobility
Toggles automatic generation of mobility events.

### PoA-in-Range
Generates the PoA in range events necessary for the state transfer service.

### Network Characteristics
Toggles automatic generation of network characteristic update events.

---
## Replay Event Pane
This pane provides the following controls for replay file execution in the backend:

### REPLAY FILE
Drop-down menu to select the replay file to execute.

_**NOTE:** New replay files created while this pane is open will only be added the the list after the pane is closed and re-opened._

### LOOP
Checkbox to enable replay file execution looping. If checked, the backend replay manager will run the selected replay file to completion, wait for 5 seconds and then restart the same replay file execution. It continues this loop until stopped.

### PLAY
Sends a request to the backend replay manager to start the selected replay file execution. While running, the following replay file status information appears in the Event configuration bar:
- Running replay file name
- Current event count / Total event count
- Time to next event / Remaining replay file execution time

### STOP
Sends a request to the backend replay manager to stop the selected replay file execution.

### CLOSE
Closes the auto-replay event pane.

---
## Process Status Table
This table is only visible when a scenario is deployed. It displays key information about applications in the deployed scenario, such as:
- Process name
- Application instance ID
- Execution status
- Ingress & egress service mappings for external applications

_**NOTE:** Ingress & egress service mappings are prepended with 'I:' & 'E:' respectively. More details on external node service mappings are provided in the [External Nodes]({{site.baseurl}}{% link docs/overview/features/overview-external-nodes.md %}) wiki page._
