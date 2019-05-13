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
- Install the [demo1](../examples/demo1/README.md) & [demo2](../examples/demo2/README.md) scenarios:
  - [Import demo1 scenario](use/base-ops.md#import-demo1-scenario-in-advantedge)
  - Use the same import procedure for demo2 scenario
- Make sure there is no deployed scenario
- [Install Cypress](#install-cypress)

## Install Cypress

To install Cypress run the following commands:

```
cd ~/AdvantEDGE/test
npm ci
```

## Run tests

Test specification files are located [here](../test/cypress/integration/tests)

### Cypress CLI

```
# Run Cypress tests using CLI
cd ~/AdvantEDGE/test
npm run cy:run
```

>**NOTE**<br>
>Default AdvantEDGE URL used by cypress is http://127.0.0.1:30000<br>
>To run tests using another deployment:<br>
>`npm run cy:run -- --env meep_url="http://<Node IP>:<MEEP FE port>"`


## Cypress GUI

```
# Run/Debug Cypress tests using GUI
cd ~/AdvantEDGE/test
npm run cy:open
```

>**NOTE**<br>
>Default AdvantEDGE URL used by cypress is http://127.0.0.1:30000<br>
>To run tests using another deployment:<br>
>`npm run cy:open -- --env meep_url="http://<Node IP>:<MEEP FE port>"`
