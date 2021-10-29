---
layout: default
title: Project FAQ
parent: Project
nav_order: 3
---

- [FAQ 1 - Unexplained disconnections using an egress service](#faq-1-unexplained-disconnections-using-an-egress-service)
- [FAQ 2 - K8s Docker container runtime deprecation](#faq-2-k8s-docker-container-runtime-deprecation)

-----

## FAQ 1: Unexplained disconnections using an egress service
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

To temporarily disable segmentation offload:
```
# xxx is the device name (eth0, eno1... etc)
ethtool -K xxx gso off gro off tso off
```

To permanently disable segmentation offload, modify the interface configuration file (`/etc/network/interfaces`) as shown below:
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

- [Serverfault question thread](https://serverfault.com/questions/193114/linux-e1000e-intel-networking-driver-problems-galore-where-do-i-start)
- [Ubuntu bug report](https://bugs.launchpad.net/ubuntu/+source/linux/+bug/1766377)

-----
## FAQ 2: K8s Docker container runtime deprecation
_**Starting with k8s version 1.22, Docker container runtime is no longer supported**_

With release 1.20, Kubernetes announced the deprecation of Docker as a container runtime, promoting instead other runtimes that support the
Container Runtime Interface (CRI) such as containerd and CRI-O. As of release 1.22, Docker container runtime was officially unsupported.

**Why deprecate Docker container runtime?**

To support interoperability with different container runtimes, Kubernetes requires a runtime that implements the Container Runtime Interface (CRI).
Docker container runtime is not CRI-compliant and requires an adaptation layer called _dockershim_, maintained by k8s. As of release 1.22, k8s decided
to stop maintaining _dockershim_ for Docker, in favor of other CRI-compatible runtimes.

More details about the deprecation can be found here:
- [K8s dockershim deprecation v1.20 release notes](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.20.md#dockershim-deprecation)
- [K8s dockershim deprecation FAQ](https://kubernetes.io/blog/2020/12/02/dockershim-faq/)
- [K8s & Docker](https://kubernetes.io/blog/2020/12/02/dont-panic-kubernetes-and-docker/)
- [K8s Docker deprecation explanation](https://medium.com/better-programming/kubernetes-is-deprecating-docker-8a9f7566fbca)

**Impact on AdvantEDGE deployment**

AdvantEDGE [runtime installation procedure]({{site.baseurl}}{% link docs/setup/env-runtime.md %}) deploys Docker as the
k8s container runtime. For this reason, AdvantEDGE installation as currently documented only supports k8s versions up to 1.21. We have validated
that the Docker container runtime can be replaced seamlessly by _containerd_ for k8s version 1.20, however we have not yet tested with the latest
k8s release 1.22.

We will update the AdvantEDGE documentation with the new container runtime installation procedure as soon as we officially support 
the latest k8s version. Users wishing to migrate to another CRI-compatible container runtime may do so at any time.
