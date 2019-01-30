# AdvantEDGE Environment Setup
## Goals
- Guidance on pre-requisite installation
- List of versions known to work

## Overview
AdvantEDGE requires the following pre-requisites
- [Ubuntu](#ubuntu)
- [Dockers](#dockers)
- [Kubernetes](#kubernetes)
- [Helm](#helm)


## Ubuntu
There are many installation guides out there; we use this one:

https://tutorials.ubuntu.com/tutorial/tutorial-install-ubuntu-desktop#0

Tested versions
- 16.04 LTS and 18.04 LTS
- Kernel: 4.4, 4.15 and 4.18

## Dockers
We typically use the convenience script procedure for the community edition:

https://docs.docker.com/install/linux/docker-ce/ubuntu/

In a nutshell
```
# Install using convenience script
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
# Add user to docker group
sudo usermod -aG docker <your-user>
```
Tested versions
- ???

## Kubernetes
We use the kubeadm method from here:

https://kubernetes.io/docs/setup/independent/install-kubeadm/

Steps summary:
```
# STEP 1
# Verify pre-requisites here: https://kubernetes.io/docs/setup/independent/install-kubeadm/#before-you-begin
# Then disable swap
sudo swapoff -a
sudo sed -i '/ swap / s/^/#/' /etc/fstab
```
```
# STEP 2
# Docker was previously installed
# Now, setup Docker daemon: https://kubernetes.io/docs/setup/cri/#docker
cat > /etc/docker/daemon.json <<EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2"
}
EOF

mkdir -p /etc/systemd/system/docker.service.d

# Restart docker.
systemctl daemon-reload
systemctl restart docker
```
```
# STEP 3
# Install kubeadm, kubelet & kubectl: https://kubernetes.io/docs/setup/independent/install-kubeadm/#installing-kubeadm-kubelet-and-kubectl
sudo apt-get update && sudo apt-get install -y apt-transport-https curl

curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -

cat <<EOF >/etc/apt/sources.list.d/kubernetes.list
deb https://apt.kubernetes.io/ kubernetes-xenial main
EOF

sudo apt-get update

sudo apt-get install -y kubelet kubeadm kubectl

# Lock current version
sudo apt-mark hold kubelet kubeadm kubectl
```
```
# STEP 4
# Initialize the master: https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#initializing-your-master
kubeadm init
```
```
# STEP 5
# Install the network add-on: https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#pod-network
# We use WeaveNet (https://www.weave.works/docs/net/latest/kubernetes/kube-addon/)
sysctl net.bridge.bridge-nf-call-iptables=1
kubectl apply -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')"
```
```
# STEP 6
# We currently support only single node K8s: https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#control-plane-node-isolation
# Allow scheduling pods on the master node
kubectl taint nodes --all node-role.kubernetes.io/master-
```
```
# STEP 7
# For convenience, allow kubectl auto-completion
echo "source <(kubectl completion bash)" >> ~/.bashrc
```
Tested versions
- 1.09, 1.10, 1.12, 1.13

## Helm
We use this procedure: https://docs.helm.sh/using_helm/#installing-helm

In a nutshell
```
# Install Helm: https://docs.helm.sh/using_helm/#installing-helm
sudo snap install helm --classic

# Install Tiller: https://docs.helm.sh/using_helm/#installing-tiller
helm init

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

# Enable incubator charts
helm repo add incubator https://kubernetes-charts-incubator.storage.googleapis.com/

helm repo update
```

Tested versions
- 2.8.2, 2.12.3
