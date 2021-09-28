---
layout: default
title: Monitoring Subsystem
parent: Features
grand_parent: Overview
nav_order: 3
permalink: docs/overview/features/monitoring/
---

## Feature Overview
AdvantEDGE provides a built-in Monitoring Subsystem that integrates with scenarios.

This feature provides the following capabilities:

- _Scenario local measurements_
  - Automated Network Characteristics: Latency, UL/DL throughput, UL/DL packet loss are automatically recorded
  - Automated Events: Scenario events generated towards the Events API are recorded; recorded events can originate from the frontend, from an external source, from a replay file or from one of the automation.
- _Custom measurements_
  - Custom metrics: InfluxDB API is available for logging your own time-series metrics; justneed to include an InfluxDB client in your application and start logging.
- _Dashboard visualization and management interface_
  - Built-in dashboards: visualize network characteristics point-to-point (source to dest.) or aggregated (source to all)
  - Custom dashboards: create your own dashboards; allows access to display automated measurements (net.char/events) with your own measurements.
- _Metrics API_
  - Expose metrics to applications: Metrics can be exposed to external applications for condicting network adaptative experiments.
- _Platform metrics local monitoring_
  - Automated Platform Micro-Services monitoring: Prometheus collects metrics locally about the platform micro-services; this allows AdvantEDGE platform usage metrics in your deployments.

### Micro-Services
- _InfluxDB:_ Time-Series database - used to monitor scenario network characteristics, events & custom user metrics.
- _Grafana:_ Dashboard visualization and management solution
- _metrics-engine:_ Collects automated measurements and implements the metrics API
- _Prometheus:_ Collects platform micro-services metrics

### Scenario Configuration
No scenario configuration

### Scenario Runtime
#### InfluxDB
Influx DB is a time series database; it provides a central aggregation point to store AdvantEDGE metrics.

Out-of-the-box collected metrics are:
- Latency
- UL/DL Throughput
- UL/DL Packet loss
- Events

InfluxDB runs as a dependency pod in the platform.
When deploying a scenario, AdvantEDGE creates a database for the scenario and stores aforementioned metrics in it.
After scenario termination, the stored metrics remain available until the scenario is re-deployed.
A user willing to preserve metrics must export these in between scenario runs.

InfluxDB is provided as a platform facility; if desired, users can use the InfluxDB database instance to store demo specific metrics & re-use them for graphing.

Externally from the platform, access to InfluxDB are proxied through Grafana.

#### Grafana
Grafana is a flexible graphing service that can pull metrics directly from known data sources such as InfluxDB or Prometheus.

Grafana integrates with AdvantEDGE by providing dashboards that are embedded in AdvantEDGE frontend.
On platform bring-up, default AdvantEDGE dashboards are imported in the platform.

Grafana is provided as a platform facility; if desired, users can use Grafana to create and store demo specific dashboards.
Grafana provides a frontend that can be accessed from the Montitoring page; using Grafana frontend.
Demo-specific dashboards can be added to the Monitoring page or the execution page.

#### Metrics engine
AdvantEDGE provides a `/metrics` endpoint in its REST API to allow user to collect/use metrics from their scenario control software or to experiment from their edge applications.

The service currently allows to query/subscribe to metrics related to:
- network KPIs (latency, UL/DL throughput, jitter, packet-loss)
- events received on the `/events` endpoint (mobility, net.char. update, etc.)
- http requests received by the various REST APIs of the platform

Example usage of this API: in a past demo, we subscribed to this API to feed scenario data (throughput usage) into a ML algorithm of ours.

#### Prometheus
Prometheus is a monitoring & alerting toolkit that collects and stores metrics; its 2 main components are:
- _Prometheus Server:_ Scrapes metrics from services and stores them in a time-series database; monitors alert conditions
- _Alert Manager:_ Manages and publishes alert notifications

Prometheus metric are stored in a database as time-series uniquely identified by metric name and applied labels. Supported metric types are:
- _Counter:_ Single increasing numerical value
- _Gauge:_ Single numerical value that can go up or down
- _Histogram:_ Observations (durations, sizes, etc.) grouped into configurable buckets; includes counters for sample number & sum
- _Summary:_ Observations (durations, sizes, etc.) with calculated quantiles; includes counters for sample number & sum

Prometheus is best used for metrics collection; by grouping data into metric types, Prometheus efficiently supports data storage, queries & alerting. It is an excellent tool for monitoring platform or system usage trends over time.

_**NOTE:** InfluxDB is better suited for event logging and long-term data storage._

##### Prometheus Server
Prometheus server pulls metrics from configured services by periodically _scraping_ the well-known `/metrics` endpoint. Each _scrape interval_, it collets samples from each configured service and stores them in the appropriate time-series.

Services wishing to provide metrics to the Prometheus server must expose the `/metrics` endpoint and create a custom `ServiceMonitor` resource. There are several readily available Prometheus exporters and libraries to easily instrument microservices for metrics exposure.

Prometheus exposes its data with the PromQL query language; allows retrieving and aggregating time series data in real time. Queries can be made using the HTTP API; Grafana uses this API and supports Prometheus as a data source for graphing data in its Dashboards.

Prometheus server also monitors its configured alert thresholds, informing the Alert Manager of any alert conditions.

##### Alert Manager
Alert Manager processes alerts received from Prometheus server. When an alert is received, the Alert Manager sends an alert notification to its configured listeners via e-mail, chat or notification systems.

Alert Manager also supports alert silencing and aggregation.

