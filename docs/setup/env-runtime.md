---
layout: default
title: Runtime Setup
parent: Setup
nav_order: 2
---

Topic | Abstract
------|------
[Ansible](#ansible) | Install using **Ansible** (_beta-feature_)
[Ubuntu](#ubuntu) | Supported OS
[Docker](#docker) | Docker installation
[Kubernetes](#kubernetes) | Kubernetes installation
[Helm](#helm) | Helm installation
[GPU Support](#gpu-support) | [Optional] To run sceanrios using GPUs
NEXT STEP: [Development environment](#next-step) |

----
## Ansible
_:exclamation: **IMPORTANT NOTE** :exclamation:<br>
With AdvantEDGE release v1.9+, Ansible playbooks are no longer maintained; they are left here for reference only.<br>_

AdvantEDGE runtime environment installation procedures can be performed manually or automatically.

- To install **manually** - Read through the following sections
- To install using **Ansible** (_beta-feature_) - follow this [link]({{site.baseurl}}{% link docs/setup/env-ansible.md %})

----
## Ubuntu

There are many installation guides out there; we use [this one](https://tutorials.ubuntu.com/tutorial/tutorial-install-ubuntu-desktop#0)

Versions we use:

- 18.04 LTS, 20.04 LTS and 22.04 LTS<br> _(version 16.04 LTS used to work - not tested anymore)_
- Kernel: 4.4, 4.15, 4.18, 5.3, 5.4 and 5.15

----
## Docker

We use the procedure for the community edition from [here](https://docs.docker.com/install/linux/docker-ce/ubuntu/)

Versions we use:

- 19.03 and 20.10 <br> _(versions 17.03, 18.03, 18.09 used to work - not tested anymore)_
- Containerd: 1.5.11 _(v1.6+ not supported)_

How we do it:

If upgrading from an older version, start by uninstalling it:

```
# Uninstall the Docker Engine, CLI, Containerd, and Docker Compose packages
sudo apt-get purge docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Delete all images, containers, and volumes
sudo rm -rf /var/lib/docker
sudo rm -rf /var/lib/containerd

sudo reboot
```

To install the latest supported version:

```
# Install dependencies
sudo apt-get update
sudo apt-get install ca-certificates curl gnupg lsb-release

# Add Docker’s official GPG key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor --yes -o /usr/share/keyrings/docker-archive-keyring.gpg

# Set up the stable repository
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker engine
sudo apt-get update
sudo apt-get install -y docker-ce=5:20.10.22~3-0~ubuntu-$(lsb_release -cs) docker-ce-cli=5:20.10.22~3-0~ubuntu-$(lsb_release -cs) containerd.io=1.6.14-1 docker-compose-plugin=2.14.1~ubuntu-$(lsb_release -cs)

# Lock current version
sudo apt-mark hold docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Add user to docker group
sudo usermod -aG docker <your-user>

# Allow user to access containerd socket 
sudo setfacl --modify user:<your-user>:rw /run/containerd/containerd.sock

# Restart shell to apply changes
```

----
## Kubernetes

_:exclamation: **BREAKING CHANGES** :exclamation:<br>
With AdvantEDGE release v1.7+, **pre-1.16 k8s releases are no longer supported**.<br>
With AdvantEDGE release v1.9+, **pre-1.19 k8s releases are no longer supported**._

_:exclamation: **IMPORTANT NOTE** :exclamation:<br>
With AdvantEDGE release v1.9+, Docker container runtime has been replaced by containerd to support k8s versions 1.22+.<br>
For more information, see the [Docker container runtime deprecation FAQ]({{site.baseurl}}{% link docs/project/project-faq.md %}#faq-2-k8s-docker-container-runtime-deprecation)._

We use the kubeadm method from [here](https://kubernetes.io/docs/setup/independent/install-kubeadm/)

Versions we use:

- 1.19 to 1.26<br> _(versions 1.16 to 1.18 used to work - not tested anymore)_

_**NOTE:** K8s deployment has a dependency on the node's IP address.<br>
From our experience, it is **strongly recommended** to ensure that your platform always gets the same IP address for the main interface when it reboots. It also makes usage of the platform easier since it will reside at a well-known IP on your network.<br>
Depending on your network setup, this can be achieved either by setting a static IP address on the host or configuring the DHCP server to always give the same IP address to your platform._

How we do it:

##### STEP 1 - Verify pre-requisites [(here)](https://kubernetes.io/docs/setup/independent/install-kubeadm/#before-you-begin)

```
# Disable swap
sudo swapoff -a
sudo sed -i '/ swap / s/^/#/' /etc/fstab
```

##### STEP 2 - Setup container runtime [(details)](https://kubernetes.io/docs/setup/production-environment/container-runtimes/)

Containerd is used as the k8s container runtime.

_**NOTE:** Containerd was installed during Docker installation._

To install the container runtime prerequisites:

```
cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
overlay
br_netfilter
EOF

sudo modprobe overlay
sudo modprobe br_netfilter

# sysctl params required by setup, params persist across reboots
cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
EOF

# Apply sysctl params without reboot
sudo sysctl --system
```

To configure containerd:

```
# configure containerd
sudo mkdir -p /etc/containerd
containerd config default | sudo tee /etc/containerd/config.toml
sudo sed -i 's/SystemdCgroup \= false/SystemdCgroup \= true/g' /etc/containerd/config.toml

# restart containerd
sudo systemctl restart containerd
```

##### STEP 3 - Install kubeadm, kubelet & kubectl [(details)](https://kubernetes.io/docs/setup/independent/install-kubeadm/#installing-kubeadm-kubelet-and-kubectl)

If upgrading from an older version, start by uninstalling it:

```
sudo kubeadm reset

sudo apt-get purge kubeadm kubectl kubelet kubernetes-cni kube*
sudo apt-get autoremove  
sudo rm -rf ~/.kube

sudo reboot
```

To install the latest supported version:

```
sudo apt-get update && sudo apt-get install -y apt-transport-https curl

curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -

sudo sh -c 'cat <<EOF >/etc/apt/sources.list.d/kubernetes.list
deb https://apt.kubernetes.io/ kubernetes-xenial main
EOF'

# Install latest supported k8s version
sudo apt-get update
sudo apt-get install -y kubelet=1.26.0-00 kubeadm=1.26.0-00 kubectl=1.26.0-00 kubernetes-cni=1.1.1-00

# Lock current version
sudo apt-mark hold kubelet kubeadm kubectl
```

##### STEP 4 - Initialize master [(details)](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#initializing-your-master)

```
sudo kubeadm init --cri-socket unix:///run/containerd/containerd.sock

# Once completed, follow onscreen instructions
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

Allow scheduling pods on master node [(details)](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#control-plane-node-isolation)

```
kubectl taint nodes --all node-role.kubernetes.io/control-plane-
# For older k8s deployments:
kubectl taint nodes --all node-role.kubernetes.io/master-
```

Install the network add-on [(details)](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#pod-network)

We use [WeaveNet](https://www.weave.works/docs/net/latest/kubernetes/kube-addon/)

```
sudo sysctl net.bridge.bridge-nf-call-iptables=1

# Based on https://github.com/weaveworks/weave/releases/download/v2.8.1/weave-daemonset-k8s.yaml
# WEAVE_MTU set to 1500
kubectl apply -f https://raw.githubusercontent.com/dilallkx/AdvantEDGE/gh-pages-v1.9.0-kd/setup/weave-daemonset-k8s.yaml
```

##### STEP 5 - Optionally add worker nodes to K8s cluster [(details)](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-join/)

_**NOTE: This step is necessary only if using Worker Nodes; if you are only using 1 node, skip this step and go to STEP #6**_

On the master node:

```
# Create token for worker nodes to join
kubeadm token create

# Get CA certificate token hash
openssl x509 -pubkey -in /etc/kubernetes/pki/ca.crt | openssl rsa -pubin -outform der 2>/dev/null | openssl dgst -sha256 -hex | sed 's/^.* //'
```

On each worker node:

```
# Enable Netfilter bridging
sudo sysctl net.bridge.bridge-nf-call-iptables=1

# Join worker node
sudo kubeadm join --token <token> <master-ip>:<master-port> --discovery-token-ca-cert-hash sha256:<hash>

# Configure the worker node
mkdir ~/.kube
scp <user>@<master-ip>:~/.kube/config ~/.kube/
```

##### STEP 6 - Enable kubectl auto-completion

_**NOTE:** This step should only be run once._

```
echo "source <(kubectl completion bash)" >> ~/.bashrc
```

##### STEP 7 - Configure Docker Registry

Each node (master & worker) must be able to access the docker registry where container images are stored. By default, we install and use a private cluster registry. To enable access to the registry, run the following commands on each node:

```
# Add the internal docker registry to the host file
# Add the following line to /etc/hosts
# <Master Node IP>   meep-docker-registry
#   example: 192.168.1.1 meep-docker-registry
sudo vi /etc/hosts

# Add K8s CA to list of trusted CAs
sudo cp /etc/kubernetes/pki/ca.crt /usr/local/share/ca-certificates/kubernetes-ca.crt
sudo chmod 644 /usr/local/share/ca-certificates/kubernetes-ca.crt
sudo update-ca-certificates

# Restart docker daemon
sudo systemctl restart docker

# Restart containerd daemon
sudo systemctl restart containerd
```

----
## Helm

We use [this](https://helm.sh/docs/intro/install/) procedure

Versions we use:

- 3.3, 3.7 <br> _(Helm v2 deprecated)_

_**NOTE:** Procedure is slightly different when upgrading Helm v2 to v3 versus installing Helm v3 from scratch_

How we do it:

### Install Helm from scratch[(details)](https://helm.sh/docs/intro/install/)

```
sudo snap install helm --channel=3.7/stable --classic

# If you have already installed helm v3, use the refresh command below
sudo snap refresh helm --channel=3.7/stable --classic
```

### Upgrade Helm v2 to v3

##### STEP 1 - Delete all your deployment running in k8s.

##### STEP 2 - Install Helm v3

```
sudo snap refresh helm --channel=3.7/stable --classic
```

##### STEP 3 - Check helm installation

```
helm version
# Output should show version as 3.7.0
```

##### STEP 4 - Download helm v2 to v3 plugin to get the helm v2 configuration and data

```
helm plugin install https://github.com/helm/helm-2to3  

helm plugin list
# This should show that 2to3 plugin is downloaded
# Note: Please check that all Helm v2 plugins that you have installed previously, work fine with the Helm v3, and remove the plugins that do not work with v3.
```

##### STEP 5 - Migrate Helm v2 configurations

```
helm 2to3 move config

helm repo list
# This will show all the repositories you had added for Helm v2
```

##### Optional Step - Clean up of Helm v2 data and releases

```
helm 2to3 cleanup
# It will clean configurations (helm v2 home directory), remove tiller and delete v2 release data. It will not be possible to restore them if you haven't made a backup of the releases. Helm v2 will not be usable afterwards.
```

----
## GPU Support

### NVIDIA

In order for Kubernetes to be aware of available GPU resources on its nodes, each host with a GPU must install the necessary drivers. The NVIDIA GPU Operator must also be installed in order to configure, install & validate all other components required to enable GPUs on k8s, such as the NVIDIA container runtime, device plugin & CUDA toolkit. More information can be found in this [blog post](https://developer.nvidia.com/blog/announcing-containerd-support-for-the-nvidia-gpu-operator/).

How we do it:

##### STEP 1 - Install NVIDIA drivers

Determine which NVIDIA GPU hardware is installed on your setup using the command `lspci | grep NVIDIA` and find the recommended driver version for your GPU by searching the [NVIDIA driver download page](https://www.nvidia.com/Download/index.aspx).

Install the NVIDIA drivers:

```
# Update the NVIDIA driver repo
sudo add-apt-repository ppa:graphics-drivers/ppa
sudo apt update

# Install the NVIDIA drivers
# sudo apt-get install nvidia-driver-<version>
sudo apt-get install nvidia-driver-510
```

Verify driver installation:

```
# Get driver information
nvidia-smi

# Sample output:
+-----------------------------------------------------------------------------+
| NVIDIA-SMI 510.73.05    Driver Version: 510.73.05    CUDA Version: 11.6     |
|-------------------------------+----------------------+----------------------+
| GPU  Name        Persistence-M| Bus-Id        Disp.A | Volatile Uncorr. ECC |
| Fan  Temp  Perf  Pwr:Usage/Cap|         Memory-Usage | GPU-Util  Compute M. |
|                               |                      |               MIG M. |
|===============================+======================+======================|
|   0  NVIDIA GeForce ...  Off  | 00000000:17:00.0 Off |                  N/A |
|  0%   36C    P8     2W / 190W |     99MiB /  6144MiB |      0%      Default |
|                               |                      |                  N/A |
+-------------------------------+----------------------+----------------------+

+-----------------------------------------------------------------------------+
| Processes:                                                                  |
|  GPU   GI   CI        PID   Type   Process name                  GPU Memory |
|        ID   ID                                                   Usage      |
|=============================================================================|
|    0   N/A  N/A      1447      G   /usr/lib/xorg/Xorg                 39MiB |
|    0   N/A  N/A      1690      G   /usr/bin/gnome-shell               57MiB |
+-----------------------------------------------------------------------------+
```

##### STEP 2 - Install NVIDIA GPU Operator

The NVIDIA GPU Operator configures, installs and validates the NVIDIA container runtime, device plugin & CUDA toolkit required to support GPUs within k8s containers. We use the NVIDIA method documented [here](https://docs.nvidia.com/datacenter/cloud-native/gpu-operator/getting-started.html#install-nvidia-gpu-operator)

_**NOTE:** This procedure will take some time during first installation_

```
# Add the NVIDIA helm repository
helm repo add nvidia https://helm.ngc.nvidia.com/nvidia && helm repo update

# Install NVIDIA GPU Operator in Bare-metal/Passthrough with pre-installed NVIDIA drivers
helm install gpu-operator --create-namespace nvidia/gpu-operator --set driver.enabled=false
```

##### STEP 3 - Deploy a scenario requiring GPU resources

This can be done via AdvantEDGE frontend scenario configuration by selecting the number of requested GPUs for a specific application (GPU type must be set to _NVIDIA_). The application image must include or be based on an official NVIDIA image containing the matching NVIDIA drivers. DockerHub images can be found [here](https://hub.docker.com/r/nvidia/cuda/).

GPU resources may also be requested via user charts in the configured AdvantEDGE scenario by adding the following lines to the container specification:

```
spec:
  containers:
    - name: <container name>
      image: nvidia/cuda:<version>
      resources:
        limits:
          nvidia.com/gpu: <# of requested GPUs. E.g. 1>
```

## Next Step

Learn about configuring the [Development Environment]({{site.baseurl}}{% link docs/setup/env-dev.md %})
