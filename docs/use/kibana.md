# Monitoring Tool
AdvantEDGE comes pre-bundled with monitoring tools.

Prior to using the monitoring features of AdvantEDGE, it is recommended to configure the monitoring visualizations/dashboard through Kibana.

The following steps shows how to import such visualizations/dashboards

## Pre-requisites
- Familiarize with [AdvantEDGE Concepts](../concepts.md)
- [Deploy AdvantEDGE](../deploy.md)

## Logging to Kibana
As per the GUI Overview where it was mentionned that the AdvantEDGE GUI can be accessed using a standard browser on port 30000 of the node where AdvantEDGE is deployeds, the same can be said for Kibana but through port 32003 
Either <your-node-ipaddress>:32003 or <your-node-URL>:32003 should do the trick

> You now successfully logged in Kibana

## Import Visualizations/Dashboards
From Kibana GUI
- Select _Management_ from the left menu
- Click on _Saved Objects_
- Click on _Import_ on the top right corner
- Click on _Import_ in the _Import saved object_ box
- Browse to `AdvantEDGE/dashboard/` and select `basic-dashboards.json`
- Click on _Open_
- Click on _Import_ button on the bottom right corner
- Click on _Done_

## Setting a Default Index Pattern
From Kibana GUI
- Select _Management_ from the left menu
- Click on _Index Patterns_
- In the _Create Index Pattern_ box, select _filebeat*_
- Click on the button with a _STAR icon_ on the top right corner

> You now selected a default index pattern that was imported from the step above. This index is used throughout all the visualizations/dashboards. The monitoring feature is now enabled.
