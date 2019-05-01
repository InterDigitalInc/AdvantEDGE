# Development Environment Setup

- Guidance on development pre-requisites installation
- List versions known to work

## Overview

AdvantEDGE development environment prerequisites:

- [Ubuntu](#ubuntu)
- [Go Toolchain](#go-toolchain)
- [Node.js & npm](#nodejs-npm)
- [Linters](#linters)
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
export GOPATH=$HOME/gocode
export PATH=$PATH:$GOPATH/bin:/usr/local/go/bin
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


## Linters

The linting tools used for development are:

- [ESLint](#eslint)
- [golangci-lint](#golangci-lint)

### ESLint

[ESLint](https://eslint.org/) is a pluggable linting utility for JavaScript.

Versions we use:

- 5.16.0

How we do it:

###### STEP 1 - Install ESLint globally

```
sudo npm install -g eslint
```

###### STEP 2 - Install ESLint React plugin

```
sudo npm install -g eslint-plugin-react
```

###### STEP 3 - Verify install

```
eslint -v

# Example output:
#   v5.16.0
```

### GolangCI-Lint

[GolangCI-Lint](https://github.com/golangci/golangci-lint) is a linters aggregator for Go.

Versions we use:

- 1.16.0

How we do it:

###### STEP 1 - Install GolangCI-Lint

```
cd ~
GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.16.0
```

###### STEP 2 - Verify install

```
golangci-lint --version

# Example output:
#   golangci-lint has version v1.16.0 built from (unknown, mod sum: "h1:PcWAN9JHflZzJQaZVY1JXZE0Tgjq+jO2v4QLqJ/Azvw=") on (unknown)
```


## IDE

There is no strict requirement on which IDE to use for development.

We use:
- [Visual Studio Code](#visual-studio-code)
- [Atom](#atom)

### Visual Studio Code

[Visual Studio Code](https://code.visualstudio.com/) is 'a lightweight but powerful source code editor which runs on your desktop.'

Versions we use:

- 1.33

Extensions we use:

- ESLint 1.8.2
- Go 0.10.1

Settings we use:

```
{
    "go.formatTool": "goimports",
    "go.lintTool":"golangci-lint",
    "go.lintFlags": [
        "--fast"
    ],
    "[go]": {
        "editor.snippetSuggestions": "none",
        "editor.formatOnSave": true,
        "editor.codeActionsOnSave": {
            "source.organizeImports": true
        },
    }
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