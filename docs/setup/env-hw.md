---
layout: default
title: Hardware Setup
nav_order: 1
parent: Setup
---

Topic | Abstract
------|------
[Hardware Requirements](#hardware-requirements) | Recomended hardware
NEXT STEP: [Runtime environment](#next-step) |

## Hardware Requirements
AdvantEDGE may be installed on a single K8s node or a cluster of K8s nodes. These may run on bare-metal Linux hosts or within VMs.

Recommended System Requirements:

- Ubuntu 20.04 LTS or 22.04 LTS
- Intel I7-8750H ~4GHz
  - 6 Cores - 12 threads
- 16GB RAM (32GB better)
- 500 GB SSD (see Note)
- Ethernet/WLAN

### Hardware Considerations
- Lesser hardware is possible, but may limit scenario size
- Certain platform features (such as geospatial calculations, network traffic modeling) may consume significant amount of processing depending on the scenario size
- Sharing a platform between multiple users (sandboxes) will share hardware resources and a user may workload my impact other users.
- We observed slow 3rd party database operation using HDD and therefore recommend SSD for better performance storage

## Next Step
Learn about configuring the [Runtime Environment]({{site.baseurl}}{% link docs/setup/env-runtime.md %})
