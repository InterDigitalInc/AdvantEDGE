![AdvantEDGE-logo](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/docs/images/AdvantEDGE-logo_Blue-01.png)

_Complete AdvantEDGE documentation is now available in the [AdvantEDGE Wiki](https://github.com/InterDigitalInc/AdvantEDGE/wiki)_

------

AdvantEDGE is a Mobile Edge Emulation Platform (MEEP) that runs on Docker & Kubernetes.

> AdvantEDGE provides an emulation environment, enabling experimentation with Edge Computing Technologies, Applications, and Services.  The platform facilitates users to explore edge / fog deployment models and their impact to applications and services in short and agile iterations.

## Motivation

- [x] **Accelerate Mobile Edge Computing adoption**
- [x] **Discover new edge application use cases & services**
- [x] **Help to answer questions such as:**
  - Where should my application components be located in the edge network?
  - How do network characteristics (such as latency, jitter, and packet loss) impact my application or service?
  - How will my application behave when the user moves within and across access networks?

## Intended Users

- [x] **Edge Application Developers**
- [x] **Edge Network and Service Designers**
- [x] **Edge Researchers**
- [x] **Technologists that are simply interested learning how the Edge works**

## Concepts

An understanding of some AdvantEDGE concepts is helpful towards effectively using the platform and understanding how it works.

Before getting started we recommend familiarity with key [AdvantEDGE concepts](https://github.com/InterDigitalInc/AdvantEDGE/wiki/platform-concepts)

## Getting started
To get started using AdvantEDGE, the following high-level steps are needed:

- Setup runtime environment (Ubuntu/Dockers/Kubernetes/Helm)
- Clone AdvantEDGE repo
- Install & Configure meepctl tool
- Deploy AdvantEDGE micro-services

Step-by-step details are available in the [Wiki](https://github.com/InterDigitalInc/AdvantEDGE/wiki#getting-started)


## Building
The backend portion of AdvantEDGE is implemented as a collection of micro-services in Golang.

The frontend portion of AdvantEDGE is implemented using Javascript, React and Redux technologies.

To re-build either of these components, you first need to setup the development environment and then use the meepctl tool.

Step-by-step details are available in the [Wiki](https://github.com/InterDigitalInc/AdvantEDGE/wiki#building)

## Testing
The AdvantEDGE platform comes with automated system tests using Cypress.

Step-by-step details are available in the [Wiki](https://github.com/InterDigitalInc/AdvantEDGE/wiki/test-advantedge)

## Upstream communication

We use GitHub issues.

So just open an issue in the repo to provide user feedback, report bugs or request enhancements.

## Licensing

Currently licensed under the [Apache 2.0 License](./LICENSE.md)
