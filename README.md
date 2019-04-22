# AdvantEDGE

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
- [x] **Technologists that are simply interestied learning how the Edge works**

## Concepts

The following AdvantEDGE concepts are described [here](docs/concepts.md)

- [x] Micro-service Architecture
- [x] Macro-network Model
- [x] Network characteristics
- [x] Network mobility
- [x] External UE support

## Getting started

- [Setup runtime environment (Ubuntu/Dockers/Kubernetes/Helm)](docs/setup_runtime.md)

- Clone the AdvantEDGE repo<br>
  `git clone https://github.com/<your-fork>/AdvantEDGE.git`<br>
  (*assuming local gitdir =* `~/AdvantEDGE`)

- Setup [*meepctl*](docs/meepctl/meepctl.md) tool
  1. Copy to an executable path<br>
    `sudo cp ~/AdvantEDGE/bin/meepctl/meepctl /usr/local/bin/`
  2. Configure<br>
    `meepctl config set --ip <your-node-ip> --gitdir /home/<user>/AdvantEDGE`

- [Deploy AdvantEDGE](docs/deploy.md)

- [Use AdvantEDGE](docs/use.md)

## Building

- [Setup development environment (Ubuntu/Go/Node.js/NPM)](docs/setup_dev.md)

- Clone the AdvantEDGE repo<br>
  `git clone https://github.com/<your-fork>/AdvantEDGE.git`<br>
  (*assuming local gitdir =* `~/AdvantEDGE`)

- [Build AdvantEDGE](docs/build.md)

## Testing

The AdvantEDGE platform test procedures are described [here](docs/testing.md)

## Concepts
An understanding of some AdvantEDGE concepts is helpful towards effectively using the platform and understanding how it works.  These core AdvantEDGE concepts are described [here](docs/concepts.md)
- [x] Micro-service Architecture
- [x] Macro-network Model
- [x] Network characteristics
- [x] Mobility
- [x] External UE support

## Upstream communication

We use GitHub issues.

So just open an issue in the repo to provide user feedback, report software bugs or request enhancements.

## Licensing

Currently licensed under the *AdvantEDGE Limited Evaluation and Use License Agreement*
