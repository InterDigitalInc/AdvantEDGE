---
layout: default
title: Monitoring View
parent: GUI Overview
grand_parent: Usage
nav_order: 3
permalink: docs/usage/gui/mon-view/
---

Topic | Abstract
------|------
[Monitoring View](#monitoring-view) | Monitoring View overview
[Headline bar](#headline-bar) | Monitoring View controls
[Dashboard configuration pane](#dashboard-configuration-pane) | Editable dashboard list
[Visualization iframe](#visualization-iframe) | Dashboard visualization

---
# Monitoring View
The Monitoring view provides an interface for the AdvantEDGE platform user to view and edit external monitoring dashboards.

Other views can be found in the [GUI Overview wiki page]({{site.baseurl}}{% link docs/usage/usage-gui-overview.md %})

# Headline Bar
This bar is always visible within the Monitoring view. It provides the following buttons:

## DASHBOARD
Loads the selected dashboard URL in the visualization iframe.

## EDIT
Opens the dashboard configuration pane to edit the list of dashboard URLs.

## OPEN GRAFANA
Loads the Grafana frontend in a new browser tab with a connection to the _meep-grafana_ backend microservice. The Grafana web application may be used to view the default dashboards or to create and edit new user-defined dashboards.

---
# Dashboard Configuration Pane
This section provides an editable list of dashboards. The default dashboards can't be removed or modified. Additional dashboards can be added, modified or removed as follows:

_**NOTE:** Any user-configured dashboards will be deleted if the web application local storage is cleared via the browser or the Settings View button to clear the UI cache._

**Create Dashboard:**
- Click _NEW_
- Enter dashboard Name & URL
- Click _APPLY_

**Edit Dashboard:**
- Modify dashboard Name & URL
- Click _APPLY_

**Delete Dashboard:**
- Click checkbox next to dashboards to be deleted
- Click _DELETE_
- Click _APPLY_

---
# Visualization Iframe
This section uses an iframe to load the selected dashboard URL into the frontend.

_**NOTE:** Certain websites don't allow iframe embedding. Dashboards will only load if the **X-Frame-Options** are not set in the http response._
