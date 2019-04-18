# Build Procedure
  
## Prerequisites

- Set up [AdvantEDGE Development Environment](setup_dev.md)

## Build & Install meepctl

The [_meepctl CLI tool_](meepctl/meepctl.md) tool is built & installed using a bash script.

```
cd ~/AdvantEDGE/go-apps/meepctl
./install.sh
```

A first time install of meepctl must also be configured.

```
meepctl config set --ip <your-node-ip> --gitdir /home/<user>/AdvantEDGE
```

## Build AdvantEDGE micro-services

The [_meepctl CLI tool_](meepctl/meepctl.md) is used to build the AdvantEDGE binaries using the [_meepctl build_](meepctl/meepctl_build.md) command.

```
meepctl build all
```

This command generates the _core_ micro-service binaries, as well as the frontend web application.

To deploy the new binaries, follow the [Deploy AdvantEDGE](deploy.md) procedure.
