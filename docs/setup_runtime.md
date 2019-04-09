# Runtime Environment Setup
- Guidance on runtime pre-requisites installation
- List versions known to work

## Overview
AdvantEDGE requires the following pre-requisites
- [Ubuntu](#ubuntu)
- [Dockers](#dockers)
- [Kubernetes](#kubernetes)
- [Helm](#helm)


## Ubuntu
There are many installation guides out there; we use [this one](https://tutorials.ubuntu.com/tutorial/tutorial-install-ubuntu-desktop#0)

Versions we use:
- 16.04 LTS and 18.04 LTS
- Kernel: 4.4, 4.15 and 4.18

## Dockers
We typically use the convenience script procedure for the community edition from [here](https://docs.docker.com/install/linux/docker-ce/ubuntu/)

Versions we use:
- 17.03, 18.03

How we do it:
```
curl -fsSL https://get.docker.com -o get-docker.sh

sudo sh get-docker.sh

# Add user to docker group
sudo usermod -aG docker <your-user>
```
## Kubernetes
We use the kubeadm method from [here](https://kubernetes.io/docs/setup/independent/install-kubeadm/)

Versions we use:
- 1.09, 1.10, 1.12, 1.13

>**IMPORTANT NOTE #1**<br>
K8s deployment has a dependency on the node's IP address.<br>
From our experience, it is **strongly recommended** to ensure that your platfrom always gets the same IP address for the main interface when it reboots. It also makes usage of the platform easier since it will reside at a well-known IPon your network.<br>
Depending on your network setup, this can be achieved either by setting a static IP address on the host or configuring the DHCP server to always give the same IP address to your platform.<br>

>**IMPORTANT NOTE #2**<br>
Latest version of K8s (1.14) has some changes incompatible with AdvantEDGE.<br>
We are currently working at resolving these.

How we do it:
###### STEP 1 - Verify pre-requisites [(here)](https://kubernetes.io/docs/setup/independent/install-kubeadm/#before-you-begin)

```
# Then disable swap
sudo swapoff -a
sudo sed -i '/ swap / s/^/#/' /etc/fstab
```
###### STEP 2 - Setup Docker daemon [(details)](https://kubernetes.io/docs/setup/cri/#docker)

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

# 
sudo mv ~/daemon.json /etc/docker

#
mkdir -p /etc/systemd/system/docker.service.d

# Restart docker.
systemctl daemon-reload
systemctl restart docker
```
###### STEP 3 - Install kubeadm, kubelet & kubectl [(details)](https://kubernetes.io/docs/setup/independent/install-kubeadm/#installing-kubeadm-kubelet-and-kubectl)
```
sudo apt-get update && sudo apt-get install -y apt-transport-https curl

curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -

sudo sh -c 'cat <<EOF >/etc/apt/sources.list.d/kubernetes.list
deb https://apt.kubernetes.io/ kubernetes-xenial main
EOF'

sudo apt-get update

sudo apt-get install -y kubelet=1.13.0-00 kubeadm=1.13.0-00 kubectl=1.13.0-00 kubernetes-cni=0.6.0-00

# Lock current version
sudo apt-mark hold kubelet kubeadm kubectl
```
###### STEP 4 - Initialize master [(details)](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#initializing-your-master)
```
kubeadm init

# Once completed, follow onscreen instructions
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

```
###### STEP 5 - Install the network add-on [(details)](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#pod-network)
We use [WeaveNet](https://www.weave.works/docs/net/latest/kubernetes/kube-addon/)
```
sudo sysctl net.bridge.bridge-nf-call-iptables=1
kubectl apply -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')"
```
###### STEP 6 - Allow scheduling pods on master node [(details)](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#control-plane-node-isolation)
We currently support only single node K8s
```
kubectl taint nodes --all node-role.kubernetes.io/master-
```
###### STEP 7 - Enable kubectl auto-completion
```
echo "source <(kubectl completion bash)" >> ~/.bashrc
```
###### STEP 8 - Enable initializers
```
# Edit kube-apiserver configuration
sudo vi /etc/kubernetes/manifests/kube-apiserver.yaml

# Add this flag to the kube-apiserver command
--runtime-config=admissionregistration.k8s.io/v1alpha1

# Reboot the platform
sudo reboot
```

## Helm
We use [this](https://docs.helm.sh/using_helm/#installing-helm) procedure

Versions we use:
- 2.8.2, 2.12.3

How we do it:
###### STEP 1 - Install Helm [(details)](https://docs.helm.sh/using_helm/#installing-helm)
```
sudo snap install helm --classic
```
###### STEP 2 - Install Tiller [(details)](https://docs.helm.sh/using_helm/#installing-tiller)
```
helm init
```
###### STEP 3 - Configure Tiller
```
# Create Tiller service account
kubectl create serviceaccount tiller --namespace kube-system

# Create Tiller cluster role binding
cat > tiller-crb.yaml <<EOF
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: tiller
subjects:
- kind: ServiceAccount
  name: tiller
  namespace: kube-system
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: ""
EOF

kubectl create -f tiller-crb.yaml

# Re-initialize Tiller with crb
helm init --service-account tiller --upgrade
```
###### STEP 4 - Configure repo
```
# Enable incubator charts
helm repo add incubator https://kubernetes-charts-incubator.storage.googleapis.com/

helm repo update
```
