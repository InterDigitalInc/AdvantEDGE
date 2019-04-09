![logo](./docs/images/advantedge.png)

AdvantEDGE is a Mobile Edge Emulation Platform (MEEP) that runs on Docker & Kubernetes

MEEP provides an environment to experiment with Mobile Edge Computing (MEC) technologies and edge / fog deployment models in short and agile iterations.

## Motivations
- [x] **Accelerate Mobile Edge Computing adoption**
- [x] **Discover new edge application use cases & services**
- [x] **Answer these questions:**
  * Where should my application components be located in the network?
  * What are network characteristics limitations of my application?
  * How will my application behave when the user moves in the network?

## Getting started
* [Setup runtime environment (Ubuntu/Dockers/Kubernetes/Helm)](docs/setup_runtime.md)

* Clone the AdvanteDGE repo<br>
  `git clone https://github.com/<your-fork>/AdvantEDGE.git`<br>
  (*assuming local gitdir =* `~/AdvantEDGE`)

* Setup [*meepctl*](docs/meepctl/meepctl.md) tool
  1. Copy to an executable path<br>
    `sudo cp ~/AdvantEDGE/bin/meepctl/meepctl /usr/local/bin/`
  2. Configure<br>
  `meepctl config set --ip <your-node-ip> --gitdir /home/<user>/AdvantEDGE`

* [Deploy AdvantEDGE](docs/deploy.md)

* [Use AdvantEDGE](docs/use.md)

## Concepts
The following AdvantEDGE concepts are described [here](docs/concepts.md)
- [x] Micro-service Architecture
- [x] Macro-network Model
- [x] Network characteristics
- [x] Network mobility
- [x] External UE support

## Upstream communication
We use GitHub issues.

So just open an issue in the repo to provide user feedback, report software bugs or request enhancements.

## Licensing
Currently licensed under the *AdvantEDGE Limited Evaluation and Use License Agreement*
