# Testing Procedure

## Summary

The AdvantEDGE platform currently supports end-to-end testing using [Cypress](https://www.cypress.io/). This Node-based Javascript testing tool simulates user interactions with the frontend and validates expected UI updates. It captures DOM snapshots for each event and allows for quick debugging of issues within the browser.

## Prerequisites

- Set up [AdvantEDGE Development Environment](setup_dev.md)
- Set up [AdvantEDGE Runtime Environment](setup_runtime.md)

## Set up testing environment

Before running the tests, do the following:

- [Build AdvantEDGE](build.md)
- [Deploy AdvantEDGE](deploy.md)
- [Install the demo1 scenario](../examples/demo1/README.md)
- Make sure there is no deployed scenario

## Run tests using Cypress CLI

```
cd js-apps
npm run cy:run
```

## Run tests using Cypress GUI

```
npm run cy:open
```


- To run tests from web-ui folder:
          - CLI: npm run cy:run [-- --env meep_url="http://<Node IP>:<MEEP FE port>"]
          - GUI: npm run cy:open [-- --env meep_url="http://<Node IP>:<MEEP FE port>"]
          - NOTES:
            - Default IP & port is 127.0.0.1:30000
           - MEEP must be running
            - No scenario must be deployed
            - demo-svc scenario must exist in scenario DB
        - Updated package versions


The [_meepctl CLI tool_](meepctl/meepctl.md) tool is built & installed using a bash script.

```
cd ~/AdvantEDGE/go-apps/meepctl
./install.sh
```

A first time install of meepctl must also be configured.

```
meepctl config set --ip <your-node-ip> --gitdir /home/<user>/AdvantEDGE
```

