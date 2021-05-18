---
layout: default
title: Settings View
parent: GUI Overview
grand_parent: Usage
nav_order: 4
permalink: docs/usage/gui/settings-view/
---

Topic | Abstract
------|------
[Settings View](#settings-view) | Settings View overview
[Development Settings](#development-settings) | Development settings details
[Local Storage Settings](#local-storage-settings) | Local storage settings details
[About Information](#about-information) | Platform information details

---
## Settings View

The Settings view provides an interface for the AdvantEDGE platform user to configure frontend settings and view platform information.

Other views can be found in the [GUI Overview wiki page]({{site.baseurl}}{% link docs/usage/usage-gui-overview.md %})

---
## Development Settings
This section provides the following developer options:

### VIS Configuration Mode
Checkbox to enable development settings for the Configuration & Execution View visualizations. If checked, you will have access to visualization configuration options to edit the [vis.js](https://visjs.org/) graph look & feel.

---
## Local Storage Settings
This section provides the following browser options:

### Clear UI Cache
Button to trigger deletion of the web application local storage. Certain frontend settings are persisted in the browser local storage to recover from a browser refresh. This button clears the local storage and resets the frontend state to its default values.

> **NOTE:** Clearing the UI cache should only be necessary in certain version upgrade cases where state values may have changed. This should not happen often.

---
## About Information
This section provides details about the running frontend version.
