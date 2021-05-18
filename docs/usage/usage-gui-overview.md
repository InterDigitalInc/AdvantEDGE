---
layout: default
title: GUI Overview
parent: Usage
nav_order: 2
has_children: true
permalink: docs/usage/gui/
---

Topic | Abstract
------|------
[GUI Overview](#gui-overview) | AdvantEDGE platform GUI overview
[Home](#home) | Platform description
[Configuration](#configuration) | Scenario management operations (create/delete/open/import/export)
[Execution](#execution) | Scenario runtime operations (deploy/terminate/events)
[Monitoring](#monitoring) | Scenario monitoring (visualization dashboards)
[Settings](#settings) | AdvantEDGE platform settings

----
## GUI Overview
The AdvantEDGE GUI can be accessed using a standard browser on standard HTTP/HTTPS ports 80/443 of the node where AdvantEDGE is deployed.

Either `<your-node-ipaddress>` or `<your-node-URL>` should do the trick.
> We currently do not perform cross-browser compatibility tests <br>Therefore your best chance of success is using Chrome

The GUI is sub-divided in 3 areas:
- Top Bar
- Main Area
- Footer

The Top Bar is a fixed element that is always visible. It contains navigation tabs and a platform health indicator.
> Health indicator verifies that the AdvantEDGE Core components are present and running

The Top Bar navigation tabs allow to select which view to display in the Main Area.

The footer is at the bottom of each page and contains a copyright notice and links to the Wiki, GitHub repository, GitHub Discussions, platform License & documentation on how to contribute to the project.

----
## Home
The Home view provides basic platform information and useful links for further learning & help.

> For public platform deployments, this is the only visible view when a user is not signed in.

----
## Configuration
The Configuration view provides a graphical interface to perform scenario management operations. You can use this page to create, update & delete scenarios in the platform backend persistent store.

Detailed usage: [Configuration View]({{site.baseurl}}{% link docs/usage/usage-configuration-view.md %})

----
## Execution
The Execution view provides a graphical interface to perform scenario run-time operations. You can use this view to deploy and terminate scenarios in the backend. You can also inject events to change the network topology or link characteristics.

Detailed usage: [Execution View]({{site.baseurl}}{% link docs/usage/usage-execution-view.md %})

----
## Monitoring
The Monitoring view provides an interface for the AdvantEDGE platform user to view and edit external monitoring dashboards.

Detailed usage: [Monitoring View]({{site.baseurl}}{% link docs/usage/usage-monitoring-view.md %})

----
## Settings
The Settings view provides an interface for the AdvantEDGE platform user to configure frontend settings and view platform information.

Detailed usage: [Settings View]({{site.baseurl}}{% link docs/usage/usage-settings-view.md %})
