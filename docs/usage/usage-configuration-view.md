---
layout: default
title: Configuration View
parent: GUI Overview
grand_parent: Usage
nav_order: 1
permalink: docs/usage/gui/config-view/
---

Topic | Abstract
------|------
[Configuration View](#configuration-view) | Configuration View overview
[Headline bar](#headline-bar) | Configuration View controls
[Scenario visualization area](#scenario-visualization-area) | Network topology & map visualization
[Network element configuration pane](#network-element-configuration-pane) | Element configuration controls
[Network element table](#network-element-table) | Configured element list

---
## Configuration View
The Configuration view provides a graphical interface to perform scenario management operations. You can use this view to create, update & delete scenarios in the platform backend persistent store.

To learn more about scenarios check out the Network Model section of the [Platform Concepts wiki page]({{site.baseurl}}{% link docs/overview/overview-architecture.md %}#network-model)

Other views can be found in the [GUI Overview wiki page]({{site.baseurl}}{% link docs/usage/usage-gui-overview.md %})

---
## Headline Bar
The headline bar is always visible within the Configuration view. It provides the following buttons:

### View
Drop-down menu for selecting visualization area mode to display:
- Network: Network topology view
- Map: Map view

### NEW
Create a new scenario:
- Opens a dialog prompt for a new scenario name
  - Validates the scenario name format & length
  - Verifies that scenario does not already exist in backend store
- Creates a new empty scenario
  - Displays the new scenario in the visualization area
  - Opens the Network Element configuration pane
  - Shows the Network Element table with a single scenario element

_**NOTE:** Operation fails silently if a scenario with the same name already exists. In this case you must click on **NEW** again and enter a unique scenario name._

_**NOTE:** At this point the scenario exists only in the frontend. Any operations performed on the scenario will be lost if the browser is refreshed. To persist the new scenario you must send it the the backend store using the **SAVE** button._

### OPEN
Open an existing scenario:
- Opens a dialog with a drop-down list of scenarios stored in the backend store
- Retrieves the selected scenario and loads it into the frontend
  - Displays the selected scenario in the visualization area
  - Opens the Network Element configuration pane
  - Shows the Network Element table with the list of elements in the selected scenario

_**NOTE:** Any unsaved changes to a scenario being configured in the frontend will be lost when another scenario is opened._

_**NOTE:** Operations performed on the scenario will only be persisted when sent to the backend store using the **SAVE** button._

### SAVE
Save scenario modifications:
- Opens a dialog prompt for the scenario name to be saved
  - Defaults to the name of the currently configured scenario
  - Provides _Save As_ functionality by allowing you to change the scenario name
  - Validates the scenario name format & length
- Sends the latest scenario to the backend where it is added to the persistent store

_**NOTE:** Once saved, there is no means of retrieving the previously saved version._

_**NOTE:** This operation overwrites any existing scenario with the same name. Be careful to use a unique name when saving to another scenario name._

### IMPORT
Import a scenario from a _yaml_ file:
- Opens a file explorer to find the scenario file to import
- Loads the chosen scenario file into the frontend
  - Displays the imported scenario in the visualization area
  - Opens the Network Element configuration pane
  - Shows the Network Element table with the list of elements in the imported scenario

_**NOTE:** Any unsaved changes to a scenario being configured in the frontend will be lost when a scenario is imported._

_**NOTE:** At this point the scenario exists only in the frontend. To persist the imported scenario you must send it the the backend store using the **SAVE** button._

_**NOTE:** Scenario validation during import operation is limited. We recommend a visual inspection of the network elements after import._

### EXPORT
Export a scenario to a _yaml_ file:
- Opens a dialog prompt for the name of the _yaml_ file to be created
  - Defaults to the name of the currently configured scenario
- Downloads the exported file in the browser
  - Creates a unique download file name if file already exists

---
## Scenario Visualization Area
This area is only visible when a scenario is being configured. It provides a visual representation of either the network element hierarchy or geographic position on a map for the scenario being configured. The View drop-down menu is used to toggle the visualization mode.
The visualization area is updated whenever a scenario is created, updated, opened, imported or deleted.

### Network View
We use a 3rd party graphing tool called [vis.js](https://visjs.org/) to create a network view of the scenario elements. You can drag and zoom the network view using the provided controls or using your mouse click & scroll buttons.

The scenario visualization shows basic network element information such as type, name & network characteristics directly in the graph. For more details on a specific element you can hover over the element to view an information tooltip, or you can click on it to open it in the network element configuration pane. You can switch between elements by clicking on a new one or on the graph background. You can set the network element positions using click & drag.

_**NOTE:** Hovering over a link will also open a tooltip with additional network characteristics information._

### Map View
We use a 3rd party mapping client called [Leaflet](https://leafletjs.com/) to render map tiles. Physical location assets such as UEs, PoAs & Compute nodes are then overlaid on the maps using markers, circles and lines. You can drag and zoom the map view using the provided controls or using your mouse click & scroll buttons.

We use [Leaflet-Geoman](https://geoman.io/leaflet-geoman) to enable location & path editing directly in the Map View for element configuration. You can click on any marker to select it for configuration and then use the buttens on the top-left of the Map View to edit the element. Once changes are finished, the element being configured is updated and the new values can then be applied to the configured scenario.

---
## Network Element Configuration Pane
This pane is only visible when a scenario is being configured. It provides a means of creating, editing, deleting & cloning network elements in the scenario. It provides the following buttons:

_**NOTE:** There is no **EDIT** button. To edit an element you must click on it in the scenario visualization area._

_**NOTE:** Any updates made to the scenario using this pane will only be persisted when sent to the backend using the headline bar **SAVE** button._

### NEW
Configure a new scenario element:
- Enabled when no element is currently selected for editing
- Provides a set of forms to configure a new network element
- When applied, adds the new element to the scenario

Element configuration fields:

##### Element Type
Drop-down list of scenario element types. Divided into layers according to the [Network Model]({{site.baseurl}}{% link docs/overview/overview-architecture.md %}#network-model).

##### Parent Node
Drop-down list of available parent nodes. Only shows valid & available parent elements.

##### Unique Element Name
Text field to input network element name. Accepts lowercase alphanumeric or '-' or '.' with a limit of 30 characters.

##### Latency / Jitter / Packet Loss / Throughput (DL/UL)
Number fields to input element network characteristics. Network characteristics are described in more detail [here]({{site.baseurl}}{% link docs/overview/overview-architecture.md %}#network-characteristics).

##### Connectivity Model
Drop-down list of supported connectivity models:
- _OPEN:_ Allows full connectivity between terminals and all other scenario elements
- _PDU:_ Limits terminal connectivity to other scenario elements according to active PDU sessions

Cellular Connectivity Control feature & usage is described in more detail [here]({{site.baseurl}}{% link docs/overview/features/overview-cell-connectivity-control.md %}).

##### Data Network Name
Text field to input data network name. Accepts alphanumeric or '-' or '.' with a limit of 50 characters.

This value may be used as a target data network for terminal PDU sessions.

Cellular Connectivity Control feature & usage is described in more detail [here]({{site.baseurl}}{% link docs/overview/features/overview-cell-connectivity-control.md %}).

##### Local (LADN)
Checkbox to indicate if the data network is a _Local Area Data Network_.

When checked, terminal connectivity to the data network is blocked when the terminal is outside the data network service area.

##### Service Provider / Edge Compute Service Provider
Text field to input service provider name. Accepts alphanumeric or ' ' with a limit of 50 characters.

_**NOTE:** Currently ignored by the platform._

##### MCC / MNC
Text fields to input 3-digit hexadecimal Mobile Country Code (MCC) and Mobile Network Code (MNC) for simulated 3GPP cellular network.

##### Default Cell ID / Cell ID
Text field to enter 28-character bit string cell identifier for a cellular PoA in the simulated 3GPP cellular network.

##### Location & Path Coordinates, Radius, End-of-path mode, Velocity & Supported Wireless Types
Text & Number input fields for provisioning geospatial data for physical & network location assets. GIS feature & usage is described in more detail [here]({{site.baseurl}}{% link docs/overview/features/overview-gis.md %}).

##### Zone Color
Text field to enter a 6-character hexadecimal zone color. Value may be set manually or using the provided color picker. This color is used in the map view for PoA icons and ranges that are in the configured zone.

##### Container Image Name
Text field to input docker container image name to run. By default, pulls images from [Docker Hub](https://hub.docker.com/) to your local registry. You may use another registry by adding the URL in the container image name (e.g. _meep-docker-registry:30001/demo-server_).

##### Port
Number field to input the service port to expose in the k8s pod. This port is pod-specific and can be accessed by all other pods.

_**NOTE:** Although the backend supports a set of service ports, the frontend currently supports a single service port. If you need to expose more than 1 service port, you must do so using user-defined charts._

##### External Port
Number field to input the service port to expose in the k8s cluster. This port maps to a k8s NodePort and must be unique within the cluster. Valid values range from 30000 to 32767.

External ports provide direct access from clients running outside the k8s cluster to a service running in a k8s pod. This connection **is not impacted by network characteristics** defined in the scenario.

_**NOTE:** To access an internal service and apply network characteristics to the connection you must configure an external node._

##### Protocol
Drop-down list of currently supported service protocol types.

##### Group Service Name
Text field to input the Mobility Group service name. This service identifies the service group to which a service instance belongs.

Several instances may join the same mobility group in order to offer a multi-edge service. Each service instance must register with the Mobility Group Manager as specified in the [Application State Transfer service]({{site.baseurl}}{% link docs/overview/edge-services/overview-state-transfer.md %}).

_**NOTE:** Multi-edge services that do not require state transfer support are also enabled using this field. If multiple service instances specify the same group service name but do not register to the Mobility Group Manager, application traffic is automatically routed to the nearest service instance based on hop count. Any changes to this behavior must be coordinated through the Mobility Group Manager._

##### GPU Count
Number field to input the number of GPUs to reserve. The pod will not be scheduled unless the requested number of GPUs are available.

##### GPU Type
Drop-down list of currently supported GPU types.

##### Min CPU Count
Number field to input the minimum number of CPUs to reserve. The pod will not be scheduled unless the requested number of CPUs are available.

CPU resources are measured in CPU units as described in the [k8s documentation](https://kubernetes.io/docs/tasks/configure-pod-container/assign-cpu-resource/#cpu-units). AdvantEDGE currently supports whole or fractional CPU units using a float.

##### Max CPU Count
Number field to input the maximum number of CPUs to reserve. The pod will not be provided more than _Max CPU_ resources.

CPU resources are measured in CPU units as described in the [k8s documentation](https://kubernetes.io/docs/tasks/configure-pod-container/assign-cpu-resource/#cpu-units). AdvantEDGE currently supports whole or fractional CPU units using a float.

##### Min Memory (MB)
Number field to input the minimum amount of memory (in MB) to reserve. The pod will not be scheduled unless the requested amount of memory is available.

##### Max Memory (MB)
Number field to input the maximum amount of memory (in MB) to reserve. The pod will not be provided more than _Max Memory_ resources.

##### Environment Variables
Text field to input a list of environment variables to set in the k8s container. Accepts a comma-separated list of variables (e.g. VAR=value\[,VAR=value\]) with alphanumeric or '_' or '-' or '.' characters.

##### Command
Text field to input the command to execute on container startup. This command replaces the default container entry point.

##### Arguments
Text field to input the arguments to pass to the container entry point command.

##### Placement Identifier
Text field to input the k8s node identifier where the pod must be scheduled. More details are provided in the [Pod Placement]({{site.baseurl}}{% link docs/overview/features/overview-pod-placement.md %}) wiki page.

##### External App / Ingress Service Mapping / Egress Service Mapping
Checkbox to indicate if the application is running externally to the k8s cluster.

When configuring an external application, text fields are enabled to input ingress & egress service mappings as detailed in the [External Nodes]({{site.baseurl}}{% link docs/overview/features/overview-external-nodes.md %}) wiki page.

##### User-Defined Chart
Checkbox to indicate if the application must be installed using a user-provided helm chart. When unchecked, the backend will automatically generate a helm chart using the configured values. When checked, the backend will install the user-provided helm chart.

When installing the user-defined chart, the _--name_ argument and the _fullnameOverride_ variable are both passed and set to the value **_\<scenario name\>-\<unique element name\>_**. This permits the installation of several instances using the same helm chart.

_**NOTE:** For more information on writing helm charts see the [Helm](https://helm.sh/) website or the [Helm Charts repository](https://github.com/helm/charts) for examples._

##### User Chart Location
Text field to input the user-defined chart location used by the Virtualization Engine. User-defined charts must be placed on the k8s master node at the following location: `~/.meep/virt-engine/user-charts/`. The Virtualization Engine uses the **relative path** provided in the _user chart location_ field to retrieve the requested user-defined charts from this folder.

_**NOTE:** For AdvantEDGE versions older than v1.5, user charts may be placed anywhere on the backend host where the Virtualization Engine is running. This user chart location must to be specified as a full path, with the exception of the '~' character that is replaced by the **$HOME** path._

##### User Chart Group
Text field to input the user chart service information. This field provides the necessary details about the user chart service to enable network characteristics and mobility group management. If not specified, the user chart services may be accessed but will bypass traffic management.

The field is formatted as follows: _`Svc instance:svc group name:port:protocol`_
- **_Svc instance:_** Name of exposed user chart service
- **_Svc group name:_** Mobility group name to which service belongs (optional)
- **_Port:_** Service port
- **_Protocol:_** Service protocol

_**NOTE:** The platform currently supports provisioning for a single user chart service._

##### User Chart Alternate Values
Text field to input the path to an alternate _values.yaml_ file to use during chart installation. Alternate _values.yaml_ files must be placed on the k8s master node at the following location: `~/.meep/virt-engine/user-charts/`. The Virtualization Engine uses the **relative path** provided in the _user chart alternate values_ field to retrieve the requested alternate _values.yaml_ files from this folder.

_**NOTE:** For AdvantEDGE versions older than v1.5, alternate **values.yaml** files may be placed anywhere on the backend host where the Virtualization Engine is running. A full path is required, with the exception of the '~' character that is replaced by the **$HOME** path._

### DELETE
Remove an existing scenario element:
- Enabled when an element is selected in the visualization area
- Removes the existing element from the scenario
  - Updates the network visualization
  - Updates the network element table

### CLONE
Copy an existing scenario element with its child elements:
- Enabled when an element is selected in the visualization area
- Creates a new element and copies all configuration field values
  - A new unique name is automatically generated
  - The parent node field is cleared and must be manually selected
- When applied, adds the cloned node to the scenario and automatically clones the child elements
  - To avoid creating an invalid scenario, the following configuration fields are cleared for all cloned child elements:
    - External Port
    - Ingress service mapping
    - Egress service mapping

### CANCEL
Cancels network element configuration changes:
- Enabled when an element is being configured
- Makes no changes to the scenario
- Returns element configuration pane to its original state

### APPLY
Applies the network element configuration changes:
- Enabled when an element is being configured
- Verifies element field values for scenario validity
  - Provides visual error feedback
- Creates a new element or updates an existing element in the scenario
  - Updates the network visualization
  - Updates the network element table
  - Returns element configuration pane to its original state

---
## Network Element Table
This table is only visible when a scenario is being configured. It provides a list of network elements for the scenario being configured. It is updated whenever a scenario is created, updated, opened, imported or deleted.

The table provides basic element information such as name, type, parent name & total element count. The table is sortable, which makes it useful for quickly identifying network elements in larger scenarios.
