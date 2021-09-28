---
layout: default
title: Demo 2 Scenario
parent: Usage
nav_order: 6
---

## Demo2
This scenario is the same as [demo1](../demo1/README.md) except that it uses _user charts_ to deploy its components instead of using dynamic chart generation.

## Prerequisites
- Running AdvantEDGE platform
- Demo1 docker images
  - [Build demo1](../demo1/README.md)
  - [Dockerize demo1](../demo1/README.md)

## Using the scenario
The following steps need to be done prior to using this scenario

### Build
To build demo2:

```
cd ~/AdvantEDGE/examples/demo2/
./build-demo2.sh
```

### Import
To import demo2 follow steps given for demo1 just import demo2-scenario.yaml from AdvantEDGE/examples/demo2/ instead of demo1-scenario.yaml:

[Import scenario in AdvantEDGE]({{site.baseurl}}{% link docs/usage/usage-basic.md %}#import-demo1-scenario-in-advantedge)

### Deploy
To deploy demo2 follow steps given for demo1 just select demo2 from dropdown instead of demo:

[Deploy scenario in AdvantEDGE]({{site.baseurl}}{% link docs/usage/usage-basic.md %}#deploy-demo1-scenario)
