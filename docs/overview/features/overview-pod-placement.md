---
layout: default
title: Pod Placement Support
parent: Features
grand_parent: Overview
nav_order: 6
permalink: docs/overview/features/pod-placement/
---

## Feature Overview
When deploying on a multi-node cluster, AdvantEDGE supports pod placement on specific nodes.

This feature provides the following capabilities:
- _Override Kubernetes placement_
  - This may be useful when the cluster nodes have specific hardware characteristics
  - Ex: Applications may require specific hardware (GPU, CPU, etc.) available on a single node only
  - Ex: Scenarios with several client/server applications may require specific node placement in order to minimize network traffic between k8s nodes  

### Micro-Services
  - _virt-engine:_ virt engine ensures that the pod will be scheduled on the specified nodes

### Scenario Configuration

Element | Description
------ | ------
Placement Identifier | _Process:_ [terminal,edge,cloud] Specifies the host name of the node where the application must be scheduled

### Scenario Runtime
The AdvantEDGE platform provides a configuration field named **_Placement Identifier_** that allows the user to specify the host name of the node where an application must be scheduled. If not specified, the default k8s scheduling algorithm is used. If specified, the k8s scheduler will only schedule the pod on a node with the requested host name. If no matching node is found in the cluster then the pod will remain in _Pending_ state.

Pod Placement compares the user-provided host name with the k8s node label `kubernetes.io/hostname=<host name>`. This label value can be obtained from the k8s cluster nodes by running the command `kubectl describe nodes | grep hostname`.

If deploying scenario using user Helm charts, the following should be added to the application Helm chart to achieve the same result:

```
affinity:
  nodeAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      nodeSelectorTerms:
      - matchExpressions:
        - key: kubernetes.io/hostname
          operator: In
          values:
          - <host name>
```
