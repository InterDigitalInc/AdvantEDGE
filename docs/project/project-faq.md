---
layout: default
title: Project FAQ
parent: Project
nav_order: 3
---

- [Unexplained disconnections using an egress service](#faq-1)

-----

## FAQ 1
_**When using an egress service with an external node, unexplained session disconnections are observed**_

"Egress service" is an AdvantEDGE feature used to include [external edge nodes]({{site.baseurl}}{% link docs/overview/features/overview-external-nodes.md %}) into a scenario. On certain systems, configuring & using an egress service will present unexplained spurious disconnections.

This condition may be related to the segmentation offload that certain network cards may offer. Segmentation offload comes in different flavors [Large Send Offload (LSO)](https://en.wikipedia.org/wiki/Large_send_offload), TSO (the TCP flavor) or GSO (generic flavor)

This problem may be related to specific hardware and/or hardware / driver combination and has been observed on NICs using kernel driver & modules `e1000e` such as the example shown below.
```
 lspci | awk '/[Nn]et/ {print $1}' | xargs -i% lspci -ks %
 00:1f.6 Ethernet controller: Intel Corporation Ethernet Connection (5) I219-LM
        Subsystem: Dell Ethernet Connection (5) I219-LM
        Kernel driver in use: e1000e
        Kernel modules: e1000e
 ```

To temporarily disable segmentation offload
```
# xxx is the device name (eth0, eno1... etc)
ethtool -K xxx gso off gro off tso off
```

To permanently disable segmentation offload, modify the interface configuration file (`/etc/network/interfaces`) as shown below
```
iface eno1 inet static
address xx.xx.xx.xx
netmask xx.xx.xx.xx
gateway xx.xx.xx.xx
broadcast xx.xx.xx.xx
dns-nameservers xx.xx.xx.xx
post-up /sbin/ethtool -K eno1 gso off gro off tso off
```
This problem has been observed as far as 2010.

Related resources:

- https://serverfault.com/questions/193114/linux-e1000e-intel-networking-driver-problems-galore-where-do-i-start
- https://bugs.launchpad.net/ubuntu/+source/linux/+bug/1766377

-----
