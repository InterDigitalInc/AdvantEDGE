# AdvantEDGE Usage Tutorials
## Goal
- [Introduction to AdvantEDGE GUI](#gui-overview)
- [Basic operations (link)](./use/base-ops.md)
- [Create a first user scenario (link)](./use/create-scenario.md)
- [Multi-edge services (link)](./use/me-services.md)
- [Application State Transfer (link)](./use/app-state-transfer.md)
- [External UEs (link](./use/ext-ue.md)

## Pre-requisites
- Familiarize with [AdvantEDGE Concepts](./concepts.md)

## GUI Overview
The AdvantEDGE GUI can be accessed using a standard browser on port 30000 of the node where AdvantEDGE is deployed

Either `<your-node-ipaddress>:30000` or `<your-node-URL>:30000` should do the trick
> We currently do not perform cross-browser compatibility tests <br>Therefore your best chance of success is using Chrome

The GUI is sub-divided in 3 main areas:
- Top Bar
- Drawer (to the left)
- Main Area

The Top Bar is a fixed element that is always visible. It identifies the AdvantEDGE Contoller application and contains a health indicator of the platform
> Health indicator verifies that the AdvantEDGE Core components are present and running

The Drawer allows to select what is visible in the Main Area. It can be hidden by clicking on the InterDigital bullet in the top bar.

Drawer Item | Description
------ | --------
_Configure_ | Scenarios management operations (create/delete/open/import/export)
_Execute_ | Scenarios runtime operations (deploy/terminate/events)
_Monitor_ | Scenario monitoring (visualization dashboards)
_Settings_ | AdvantEDGE platform settings
