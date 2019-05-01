# Monitoring
AdvantEDGE uses Elastic Stack to provide monitoring capabilities.

Prior to using AdvantEDGE's monitoring feature, it is necessary to configure it.
The following steps explain how to configure monitoring.

## Pre-requisites
- Familiarize with [AdvantEDGE Concepts](../concepts.md)
- [Deploy AdvantEDGE](../deploy.md)

## Elastic Stack
AdvantEDGE uses [Elastic Stack](https://www.elastic.co/products/) as a monitoring pipeline. It provides centralized logging for AdvantEDGE scenarios and core components. Elastic Stack components run in pods on the platform.

[Kibana](https://www.elastic.co/products/kibana) is the visualization component of Elastic Stack, it runs in a pod on the K8s cluster and provides a frontend of its own.

To access Kibana frontend, open AdvantEDGE frontend `<your-node-ipaddress>:30000`, then select Monitor from the drawer and click the Kibana button in the monitoring tab.

This will open a new browser tab with the Kibana frontend in it.

## Configure Monitoring
Prior to using the monitoring pipeline, it is necessary to configure it.<br>
Configuration is done by issuing the meepctl command as follow:<br>
'meepctl config kibana'<br>

As a result, all the saved objects needed to populate Kibana's dashboards and index pattern are set.
