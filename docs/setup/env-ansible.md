---
layout: default
title: Ansible (beta)
nav_order: 4
parent: Setup
---

AdvantEDGE comes with a set of Ansible playbooks to streamline platform environment setup.

_**Playbooks are currently a beta-feature**<br>
They are available for configuring the runtime & development environments<br><br>
We've used them to setup multiple systems/VMs and are confident they will work, so you are welcome to try them<br>
To help us sanitize Playbooks, let us know in [GitHub Issues](https://github.com/InterDigitalInc/AdvantEDGE/issues) if you encounter problems._

## Ansible Installation
To run playbooks, you must first install ansible on your management system. This machine can be one of your AdvantEDGE platform cluster nodes or any other system with SSH access to the cluster.

Official install procedure and Ansible documentation can be found [here](https://docs.ansible.com/)

```
# Ubuntu install procedure
sudo apt-get update
sudo apt-get install software-properties-common
sudo apt-get install ansible
sudo apt-get install ssh sshpass
```

## Playbook Inventory Configuration
Playbooks require an inventory of hosts where playbook tasks will be run. The _hosts.ini_ inventory file is used to configure the list of hosts. It can be found in the AdvantEDGE [playbooks folder](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/playbooks).

AdvantEDGE hosts are grouped into _master_, _worker_, _cluster_ and _dev_ hosts. These groups are used by the playbooks to run group-specific tasks for each of the managed hosts.

You must update the _hosts.ini_ inventory file to match your setup before running the playbooks.

## Usage
To run a playbook you must:
- Enable SSH access to your managed hosts (e.g. you must be able to ssh into the system)
- Update the _hosts.ini_ inventory file (the one in the repo's playbooks folder) with host IP address(es) and user name(s)
- Run the playbook using the _ansible-playbook_ command:

```
# Install Runtime Environment
ANSIBLE_CONFIG=./ansible.cfg ansible-playbook install-runtime-env.yml

# Uninstall Runtime Environment
ANSIBLE_CONFIG=./ansible.cfg ansible-playbook uninstall-runtime-env.yml

# Install Development Environment
ANSIBLE_CONFIG=./ansible.cfg ansible-playbook install-development-env.yml
```

_**NOTE:** You will be prompted for SSH & sudo passwords when you run a playbook.<br>
This allows ansible to connect to the managed hosts and to run tasks with elevated privileges when necessary._
