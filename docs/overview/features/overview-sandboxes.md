---
layout: default
title: Sandbox Subsystem
parent: Features
grand_parent: Overview
nav_order: 4
permalink: docs/overview/features/sandbox/
---

## Feature Overview
AdvantEDGE provides a built-in Sandbox Subsystem that allows to share the platform with multiple friendly users.

This feature provides the following capabilities:
- _Sandbox management_
  - Create/Delete sandboxes
- _Scenario isolation_
  - Execute/monitor/terminate scenarios in an isolated manner
- _Collaboration_
  - Allows multiple users to observe the same sandbox

### Micro-Services
See [Micro-service Architecture]({{site.baseurl}}{% link docs/overview/overview-architecture.md %}#micro-service-architecture) to learn about the micro-services running in a sandbox.

### Scenario Configuration
No scenario configuration

### Scenario Runtime
|Note|
|:-------------|
| AdvantEDGE default deployment configuration provided in the repository assumes a **friendly & secure environment** where co-workers/collaborators can share a single platform (e.g. lab, home, etc.) |

In order to manage expectations, the following is to be expected when using AdvantEDGE in a multi-users environment.
- platform does not provide user authentication by default
- platform does not require REST API keys - hence all endpoints are exposed by default

Because of the above, default configuration is not meant to be deployed in a public manner or an insecure environment.
The following should therefore be expected:
- using the AdvantEDGE frontend - users can impact other users
- users can view/modify/delete all scenarios in the scenario store
- users can view/deploy/terminate execution of scenario in other sandboxes
- users can view/add/modify/delete dashboards present in Grafana
- sandbox lifecycle monitoring & management is not implemented in version 1.5 and has to be manually performed using the frontend

#### Using Sandboxes

Using the platform frontend, in the execution tab, a user can create a sandbox by using the `NEW` sandbox button in the upper left corner.

![sandbox-create.md]({{site.baseurl}}/assets/images/sandbox-create.png)

This opens a dialog where you must enter the name of the sandbox to be created.

![sandbox-create.md]({{site.baseurl}}/assets/images/sandbox-create-dialog.png)

Once the sandbox is created, the scenario buttons are enabled and a scenario can be deployed as per the usual procedure.<br> Once a scenario is deployed, everything that is shown in the dashboards is related to the scenario running in the sandbox.

Selecting another sandbox using the drop-down menu switches the dashboards and displayed information to the new sandbox.

A single user or many users can create multiple sandboxes until they exhaust the system resources.

![sandbox-create2.md]({{site.baseurl}}/assets/images/sandbox-select.png)
