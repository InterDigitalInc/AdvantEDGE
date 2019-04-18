# Development Environment Setup

- Guidance on development pre-requisites installation
- List versions known to work

## Overview

AdvantEDGE development environment prerequisites:

- [Ubuntu](#ubuntu)
- [Go Toolchain](#go-toolchain)
- [Node.js & npm](#nodejs-npm)
- [IDE](#IDE)

## Ubuntu

See [Ubuntu runtime setup](setup_runtime.md#ubuntu)

## Go Toolchain

We use the official [Go Programming Language install procedure](https://golang.org/doc/install)

Versions we use:

- 1.12.1, 1.12.4

>**IMPORTANT NOTE**<br>
Minumum required version 1.12 to support Go modules<br>

How we do it:

###### STEP 1 - Download Go linux tarball [(here)](https://golang.org/dl/)

###### STEP 2 - Unzip tarball & install

```
# Example tarball: go1.12.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
```

###### STEP 3 - Create Go working directory

```
mkdir ~/gocode
mkdir ~/gocode/bin
```

###### STEP 4 - Update PATH in profile

```
# Add the following lines at the end of your $HOME/.profile
GOPATH=$HOME/gocode
PATH=$PATH:$GOPATH/bin:/usr/local/go/bin
```

###### STEP 5 - Apply profile updates

```
source ~/.profile
```

###### STEP 6 - Verify install

```
which go
go version

# Example output:
#   /usr/local/go/bin/go
#   go version go1.12.4 linux/amd64
```

## Node.js & npm

We use the _How To Install Using NVM_ method from [here](https://www.digitalocean.com/community/tutorials/how-to-install-node-js-on-ubuntu-16-04)

NVM is the _Node Version Manager (NVM)_ tool used to install Node.js. It allows concurrent use of different Node.js installations and simplifies the version upgrade procedure.

Versions we use:

- NVM: 0.34.0
- Node.js: 8.11.1, 10.15.3
- npm: 6.8.0, 6.9.0

How we do it:

###### STEP 1 - Install dependencies

```
sudo apt-get update
sudo apt-get install build-essential libssl-dev
```

###### STEP 2 - Download & install NVM

```
curl -sL https://raw.githubusercontent.com/creationix/nvm/v0.34.0/install.sh -o install_nvm.sh
bash install_nvm.sh
```

###### STEP 3 - Apply profile updates

```
source ~/.profile
```

###### STEP 4 - Install Node.js (latest LTS version)

```
# Retrieve & install latest LTS Node.js versions
nvm ls-remote | grep "Latest LTS"
nvm install 10.15.3
```

###### STEP 5 - Update NPM version bundled with Node.js

```
npm install -g npm
```

###### STEP 6 - Verify install

```
node -v
npm -v

# Example output:
#   v10.15.3
#   6.9.0
```

## IDE

There is no strict requirement on which IDE to use for development.

We use:
- [Visual Studio Code](https://code.visualstudio.com/)
- [Atom](https://ide.atom.io/)