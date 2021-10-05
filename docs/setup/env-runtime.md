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
[Dockers](#dockers) | Dockers installation
[Kubernetes](#kubernetes) | Kubernetes installation
[Helm](#helm) | Helm installation
[GPU Support](#gpu-support) | [Optional] To run sceanrios using GPUs
NEXT STEP: [Development environment](#next-step) |

----
## Ansible
AdvantEDGE runtime environment installation procedures can be performed manually or automatically.

- To install **manually** - Read through the following sections
- To install using **Ansible** (_beta-feature_) - follow this [link]({{site.baseurl}}{% link docs/setup/env-ansible.md %})

----
## Ubuntu

There are many installation guides out there; we use [this one](https://tutorials.ubuntu.com/tutorial/tutorial-install-ubuntu-desktop#0)

Versions we use:

- 18.04 LTS and 20.04 LTS <br> _(version 16.04 LTS used to work - not tested anymore)_
- Kernel: 4.4, 4.15, 4.18, 5.3 and 5.4

----
## Dockers

We typically use the convenience script procedure for the community edition from [here](https://docs.docker.com/install/linux/docker-ce/ubuntu/)

Versions we use:

- 19.03 and 20.10 <br> _(versions 17.03, 18.03, 18.09 used to work - not tested anymore)_

How we do it:

```
curl -fsSL https://get.docker.com -o get-docker.sh

sudo sh get-docker.sh

# Add user to docker group
sudo usermod -aG docker <your-user>

# Restart shell to apply changes
```

----
## Kubernetes

_:exclamation: **BREAKING CHANGE** :exclamation:<br> With AdvantEDGE release v1.7+, **pre-1.16 k8s releases are no longer supported**._

We use the kubeadm method from [here](https://kubernetes.io/docs/setup/independent/install-kubeadm/)

Versions we use:

- 1.19, 1.20 <br> _(versions 1.16 used to work - not tested anymore)_

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

##### STEP 2 - Setup Docker daemon [(details)](https://kubernetes.io/docs/setup/cri/#docker)

```
# Docker was previously installed
# Now, setup Docker daemon
cat > ~/daemon.json <<EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2"
}
EOF

# Copy daemon config
sudo mv ~/daemon.json /etc/docker

# Create systemd entry for docker daemon
sudo mkdir -p /etc/systemd/system/docker.service.d

# Reboot
sudo reboot
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
sudo apt-get install -y kubelet=1.19.1-00 kubeadm=1.19.1-00 kubectl=1.19.1-00 kubernetes-cni=0.8.7-00

# Lock current version
sudo apt-mark hold kubelet kubeadm kubectl
```

##### STEP 4 - Initialize master [(details)](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#initializing-your-master)

```
sudo kubeadm init

# Once completed, follow onscreen instructions
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

Allow scheduling pods on master node [(details)](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#control-plane-node-isolation)

```
kubectl taint nodes --all node-role.kubernetes.io/master-
```

Install the network add-on [(details)](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#pod-network)

We use [WeaveNet](https://www.weave.works/docs/net/latest/kubernetes/kube-addon/)

```
sudo sysctl net.bridge.bridge-nf-call-iptables=1
kubectl apply -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')&env.WEAVE_MTU=1500"
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
kubeadm join --token <token> <master-ip>:<master-port> --discovery-token-ca-cert-hash sha256:<hash>

# Configure the worker node
mkdir ~/.kube
scp <user>@<master-ip>:~/.kube/config ~/.kube/
```

##### STEP 6 - Enable kubectl auto-completion

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
```

----
## Helm

We use [this](https://helm.sh/docs/intro/install/) procedure

Versions we use:

- 3.3 <br> _(Helm v2 deprecated)_

_**NOTE:** Procedure is slightly different when upgrading Helm v2 to v3 versus installing Helm v3 from scratch_

How we do it:

### Install Helm from scratch[(details)](https://helm.sh/docs/intro/install/)

```
sudo snap install helm --channel=3.3/stable --classic

# If you have already installed helm v3, use the refresh command below to configure it to 3.3 instead
sudo snap refresh helm --channel=3.3/stable --classic
```

### Upgrade Helm v2 to v3

##### STEP 1 - Delete all your deployment running in k8s.

##### STEP 2 - Install Helm v3

```
sudo snap refresh helm --channel=3.3/stable --classic
```

##### STEP 3 - Check helm installation

```
helm version
# Output should show version as 3.3.1
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

In order for Kubernetes to be aware of available GPU resources on its nodes, each host with a GPU must install the necessary drivers and vendor-specific device plugin. Also, the docker container runtime must be changed to make the GPUs visible within the containers. More information can be found in this [blog post](https://devblogs.nvidia.com/gpu-containers-runtime/).

How we do it:

##### STEP 1 - Install NVIDIA drivers

Determine which NVIDIA GPU hardware is installed on your setup using the command `lspci | grep NVIDIA` and find the recommended driver version for your GPU by searching the [NVIDIA driver download page](https://www.nvidia.com/Download/index.aspx).

Install the NVIDIA drivers:

```
# Update the NVIDIA driver repo
sudo add-apt-repository ppa:graphics-drivers/ppa
sudo apt update

# Install the NVIDIA drivers
# sudo apt-get install nvidia-<version>
sudo apt-get install nvidia-415
```

Verify driver installation:

```
# Get driver information
nvidia-smi

# Sample output:
+-----------------------------------------------------------------------------+
| NVIDIA-SMI 415.27       Driver Version: 415.27       CUDA Version: 10.0     |
|-------------------------------+----------------------+----------------------+
| GPU  Name        Persistence-M| Bus-Id        Disp.A | Volatile Uncorr. ECC |
| Fan  Temp  Perf  Pwr:Usage/Cap|         Memory-Usage | GPU-Util  Compute M. |
|===============================+======================+======================|
|   0  GeForce GTX 1050    Off  | 00000000:01:00.0  On |                  N/A |
| 40%   34C    P8    N/A /  75W |    296MiB /  1999MiB |      0%      Default |
+-------------------------------+----------------------+----------------------+

+-----------------------------------------------------------------------------+
| Processes:                                                       GPU Memory |
|  GPU       PID   Type   Process name                             Usage      |
|=============================================================================|
|    0       983      G   /usr/lib/xorg/Xorg                           161MiB |
|    0      7315      G   compiz                                       131MiB |
+-----------------------------------------------------------------------------+
```

##### STEP 2 - Install NVIDIA Container Runtime

Starting with Docker 19.03, NVIDIA GPU support is included in the default _runc_ container runtime. However, the NVIDIA device plugin requires the NVIDIA container runtime. We describe how to install it here.

We use the [NVIDIA Container Runtime for Docker](https://github.com/NVIDIA/nvidia-docker) procedure.

_**IMPORTANT NOTE:** For older versions of docker you must install the nvidia-docker2 runtime as described [here](https://github.com/NVIDIA/nvidia-docker#upgrading-with-nvidia-docker2-deprecated)_

Install container-toolkit & nvidia-runtime, and verify _runc_ runtime GPU support:

```
# Add the package repositories
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | sudo apt-key add -
curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | sudo tee /etc/apt/sources.list.d/nvidia-docker.list

sudo apt-get update && sudo apt-get install -y nvidia-container-toolkit nvidia-container-runtime
sudo systemctl restart docker

# Test nvidia-smi with the latest supported official CUDA image
docker run --gpus all nvidia/cuda:9.0-base nvidia-smi
```

Update the default docker runtime by setting the following in `/etc/docker/daemon.json`:

```
{
  "default-runtime": "nvidia",
  "runtimes": {
    "nvidia": {
      "path": "/usr/bin/nvidia-container-runtime",
      "runtimeArgs": [
      ]
    }
  },
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2"
}
```

Restart the docker daemon & verify _nvidia_ runtime GPU support:

```
sudo systemctl restart docker

# Test nvidia-smi with the latest supported official CUDA image
docker run --rm nvidia/cuda:9.0-base nvidia-smi
```

##### STEP 3 - Install NVIDIA device plugin for Kubernetes

We use the [NVIDIA Device Plugin for Kubernetes](https://github.com/nvidia/k8s-device-plugin) procedure.

Install the device plugin DaemonSet using the following command:

```
kubectl create -f https://raw.githubusercontent.com/NVIDIA/k8s-device-plugin/1.0.0-beta/nvidia-device-plugin.yml
```

##### STEP 4 - Deploy a scenario requiring GPU resources

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
