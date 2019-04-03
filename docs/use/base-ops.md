# Basic Operations
AdvantEDGE comes pre-bundled with a demo scenario that allows for rapid experimentation.

Going through the deployment/termination steps of that scenario is a good introduction to basic AdvantEDGE operations.

Prior to using the demo scenario for the first time, we need to containerize the scenario applications and import the scenario in AdvantEDGE

## Pre-requisites
- Familiarize with [AdvantEDGE Concepts](../concepts.md)
- [Deploy AdvantEDGE](../deploy.md)

## Containerize Demo Applications
In a command line shell
- Go to the `AdvantEDGE/examples/demo1/bin` directory
- Execute `./dockerize.sh`

> After completing this steps, Demo Application binaries are dockerized (containerized) and the container images are stored in the local Docker registry

## Import Demo Scenario in AdvantEDGE
From AdvantEDGE GUI
- Select _Configure_ from Drawer
- Click on _Import_
- Browse to `AdvantEDGE/examples/demo1/` and select `demo1.yaml`
- Once the scenario topology appears, click on _Save_

> You successfully imported the scenario in AdvantEDGE internal storage <br>Next time you need to use it, simply click on _Open_ and select it from the drop-down menu
> Familiarize with [various demo scenarios](../../examples/README.md)

## [Back to usage top level](../use.md)
