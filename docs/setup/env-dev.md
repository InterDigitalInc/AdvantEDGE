---
layout: default
title: Development Setup
parent: Setup
nav_order: 3
---

Topic | Abstract
------|------
[Ansible](#ansible) | Install using **Ansible** (_beta-feature_)
[Ubuntu](#ubuntu) | Supported OS
[Go Toolchain](#go-toolchain) | Golang toolchain to build backend
[Node.js & npm](#nodejs---npm) | NodeJS/NPM to build frontend
[Linters](#linters) | Go/JS linters
[IDE](#ide) | IDE we like
NEXT STEP: [Platform Management Workflow](#next-step) |

----
## Ansible

AdvantEDGE development environment installation procedures can be performed manually or automatically.

- To install **manually** - Read through the following sections
- To install using **Ansible** (_beta-feature_) - follow this [link]({{site.baseurl}}{% link docs/setup/env-ansible.md %})

----
## Ubuntu

See [Ubuntu runtime setup]({{site.baseurl}}{% link docs/setup/env-runtime.md %}#ubuntu)

----
## Go Toolchain

_:exclamation: **BREAKING CHANGE** :exclamation:<br> With AdvantEDGE release v1.7+, **pre-1.13 Go toolchains are no longer supported**._

We use the official [Go Programming Language install procedure](https://golang.org/doc/install)

Versions we use:

- 1.16 - 1.19 <br> (versions 1.14, 1.15 used to work - not tested anymore)

How we do it:

##### STEP 1 - Download Go linux tarball [(here)](https://golang.org/dl/)

##### STEP 2 - Unzip tarball & install

```
# Example tarball: go1.19.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
```

##### STEP 3 - Create Go working directory

```
mkdir -p ~/gocode/bin
```

##### STEP 4 - Update PATH in bashrc

```
# Add the following lines at the end of your $HOME/.bashrc
export GOPATH=$HOME/gocode
export PATH=$PATH:$GOPATH/bin:/usr/local/go/bin
```

##### STEP 5 - Apply profile updates

```
source ~/.bashrc
```

##### STEP 6 - Verify install

```
which go
go version

# Example output:
#   /usr/local/go/bin/go
#   go version go1.19.4 linux/amd64
```

----
## Node.js - npm

We use the _How To Install Using NVM_ method from [here](https://www.digitalocean.com/community/tutorials/how-to-install-node-js-on-ubuntu-16-04)

NVM is the _Node Version Manager (NVM)_ tool used to install Node.js. It allows concurrent use of different Node.js installations and simplifies the version upgrade procedure.

Versions we use:

- NVM: 0.34.0
- Node.js: 8.11, 10.15, 10.16, 10.19, 12.22
- npm: 6.8, 6.9, 6.11, 6.12, 6.14

How we do it:

##### STEP 1 - Install dependencies

```
sudo apt-get update
sudo apt-get install build-essential libssl-dev
```

##### STEP 2 - Download & install NVM

```
curl -skL https://raw.githubusercontent.com/creationix/nvm/v0.34.0/install.sh -o install_nvm.sh
bash install_nvm.sh
```

##### STEP 3 - Apply profile updates

```
source ~/.profile
```

##### STEP 4 - Install Node.js (latest LTS version)

```
# Retrieve & install latest v12 Node.js versions
nvm ls-remote | grep "Latest LTS"
nvm install 12.22.12
```

##### STEP 5 - Update NPM version bundled with Node.js

```
npm install -g npm@6.14.16
```

##### STEP 6 - Verify install

```
node -v
npm -v

# Example output:
#   v12.22.12
#   6.14.16
```


----
## Linters

The linting tools used for development are:

- [ESLint](#eslint)
- [golangci-lint](#golangci-lint)

### ESLint

[ESLint](https://eslint.org/) is a pluggable linting utility for JavaScript.

Versions we use:

- 5.16.0

How we do it:

##### STEP 1 - Install ESLint globally

```
npm install -g eslint@5.16.0
```

##### STEP 2 - Install ESLint React plugin

```
npm install -g eslint-plugin-react
```

##### STEP 3 - Verify install

```
eslint -v

# Example output:
#   v5.16.0
```

### GolangCI-Lint

[GolangCI-Lint](https://golangci-lint.run/) is a linters aggregator for Go.

Versions we use:

- 1.46.0

How we do it:

##### STEP 1 - Install GolangCI-Lint

```
# binary will be $(go env GOPATH)/bin/golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin 
```

##### STEP 2 - Verify install

```
golangci-lint --version

# Example output:
# golangci-lint has version 1.50.1 built from 8926a95f on 2022-10-22T10:50:47Z
```

----
## IDE

There is no strict requirement on which IDE to use for development.

We use:
- [Visual Studio Code](#visual-studio-code)
- [Atom](#atom)

### Visual Studio Code

[Visual Studio Code](https://code.visualstudio.com/) is 'a lightweight but powerful source code editor which runs on your desktop.'

Versions we use:

- 1.55

Extensions we use:

- ESLint 2.1.20
- Go 0.24.2

Settings we use:

```
{
    "files.eol": "\n",
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint",
    "go.lintFlags": [
        "--fast"
    ],
    "go.useLanguageServer": false,
    "go.testOnSave": true,
    "[go]": {
        "editor.snippetSuggestions": "none",
        "editor.formatOnSave": true,
        "editor.codeActionsOnSave": {
            "source.organizeImports": true
        },
    },
    "[javascript]": {
        "editor.codeActionsOnSave": {
            "source.fixAll.eslint": true,
        },
    },
    "eslint.format.enable": true
}
```

### Atom

[Atom](https://ide.atom.io/) is 'a hackable text editor for the 21st Century.'

Versions we use:

- 1.36.1

Packages we use:

- atom-ide-ui 0.13.0
- go-plus 6.1.0

Settings we use:

```
"*":
  "atom-ide-ui":
    "atom-ide-code-format":
      formatOnSave: true
    "atom-ide-diagnostics-ui":
      showDirectoryColumn: true
  "go-plus":
    lint:
      args: "--fast"
      tool: "golangci-lint"
    test: {}
```

## Next Step
Learn about the [Platform Management Workflow]({{site.baseurl}}{% link docs/platform-mgmt/mgmt-workflow.md %})
