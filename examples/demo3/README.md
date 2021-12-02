# Demo3

Demo 3 demonstrates how to use AdvantEDGE mec services showcasing Application Mobility Service (AMS) and Edge Application Enablement Service API.

Demo 3 compose a user MEC application that can be integrated with AdvantEDGE scenarios or MEC Sandbox. 

###  Use-cases

A user can deploy demo 3 as a scenario onto AdvantEDGE or as an external application to use MEC Services APIs from a MEC application


Demo3 provides a GUI to interact with MEC Service APIs such as:

1. Edge Platform Application Enablement allow MEC Applications to interact with the MEC System allowing Demo 3 to perform the following:

- MEC Application registration/deregistration 
- MEC Service discovery & offering 
- Monitor event notifications about service and application availability 

2. Application Mobility Service provides support for relocation of user context and/or application instance between MEC hosts enabling:

- AMS to trigger and assists the context transfer 
- Demo3 user-context is transferred from one instance to another user-based target mec application



## Getting Started:
Demo 3 can be started on MEC Sandbox or AdvantEDGE

###  Sandbox Procedure

- Demo 3 does not have any prior knowledge or configuration information of the MEC services offered by the MEC platform

The following steps need to be done prior to running demo 3

| Operation: | Notes: |
| ---------  | ------ |
| 1. Login via the MEC Sandbox frontend |  |
| 2. Select a network to deploy in the user sandbox |  |
| 3a. Pre-configure MEC Application named `demo3-config.yaml` under Demo 3 backend directory refer to - [File Structure](#file-structure)  | Fill configuration fields <br> `mode:  'sandbox' `| 
| 3b. Pre-configure MEC Application `sandbox` with Application Enablement service endpoints | Example: <br> `sandbox: 'https://try-mec.etsi.org/<my-sandbox-key>/<mep-host>'` <br> MEC Application will learn MEC services availability via mep host interface
| 3c. Pre-configure MEC Application `https` if sandbox url is using https and `mepplatform` with mec platform name demo-3 will run on |  Example: <br> `https: 'false'` <br> `mecplatform: 'mep1' `| 
| 3d. Pre-configure MEC Application `appid` with an Application Instance ID (e.g. appInstanceId) | 
| 3e. Pre-configure MEC Application `localurl` with your I.P address and `port` to indicate port number that demo3 server will run at | `localurl: 'http://<my-ip-address>'` <br> `port: '<my-port-number>'` 
| 4. Optional: If running a dual mep scenario on MEC sandbox. The above steps needs to be repeated to run a seperate instance of demo 3 application by applying configurations into `demo3-config-instance-two.yaml` |  | 

How configuration is expected :
```sh .env
# This file defines the configuration of Demo3 edge application. All fields are required to run demo-3 on MEC Sandbox 

# Set where mec application is running either on MEC Sandbox or AdvantEDGE. Expected fields: sandbox | advantedge
mode: 'sandbox'
# Set url of mec platform. Example field format: http://{MEC_IP_ADDRESS}/{SANDBOX_NAME}/{MEP_NAME}/ 
sandbox: 'http://{mec-host}/{sandbox-key}/{platform-name}'
# Set if sandbox url uses https. Expected fields: true | false 
https: false 
# Set the mec platform name demo-3 will run on. Example field: mep1
mecplatform: ''
# Set user-application ID that is generated on MEC Sandbox frontend. Example field format: 7930ba6d-4581-444c-b966-3312517f3a51
appid: ''
# Set host address of demo-3. 
localurl: 'http://{local-url}/'
# Set host port number of demo-3. Example field: '8093'
port: '8093'

```

### Build Demo 3 Server

```shell
# Build demo 3 backend binary 
cd AdvantEdge/examples/demo3/src/backend
go build -o demo-server .
go run demo-server mep1.yaml
```

### Build Demo 3 Frontend 
1. Change directories to demo 3 frontend
```
cd ~/AdvantEDGE/examples/demo3/src/frontend
```
2. Install dependencies
```
npm i 
```
3. Repeat above step by installing dependencies for frontend client package
```
cd ~/AdvantEDGE/examples/demo3/src/client
npm i 
```
4. Set up your .env file with values of where demo 3 backend is served  
```
cd AdvantEdge/examples/demo3/src
# Modify .env
# Apply configurations 
ENVIRONMENT=SANDBOX
URL=http://<my-ip-address>:<my-demo3-port-number>
```
5. Run demo 3 in development mode or build into binaries
```
npm run build:dev
npm run build
```

## Demo 3 AdvantEdge Procedure

### Demo server & Demo frontend 

Demo server is a web server that will run internally in AdvantEdge as the backend of demo 3 application. The frontend is a dashboard to provide information on the MEC application instance in the scenario. One capability When the external UE moves in the network and transitions from one edge instance to another, the "UE state" (e.g. the counter value) is transferred using the application state transfer. On the demo server, it maintains a state of external UE as a counter that will continue incrementing (e.g. not reset to zero) when the UE moves in the network.



The following steps need to be done prior to using this scenario

#### Configure demo frontend and obtain binaries

##### Build from source

To build _demo-server_ binaries from source code:

```
cd ~/AdvantEDGE/examples/demo3/
./build-demo3.sh
```

#### Dockerize demo applications

Demo Application binaries must be dockerized (containerized) as container images in the Docker registry. This step is necessary every time the demo binaries are updated.

> **NOTE:** Make sure you have deployed the AdvantEDGE dependencies (e.g. docker registry) before dockerizing the demo binaries.

To generate docker images from demo binary files:

```
cd ~/AdvantEDGE/examples/demo3/
./dockerize.sh
```

#### Using the scenario & deploy

Import `demo3-scenario.yaml` under demo3 directory into AdvantEDGE as a scenario process then deploy

#### Start demo 3 frontend

To start the application, load the following page in the browser `<AdvantEDGE-node-ip-address>:31111`for first instance of demo 3 | `<AdvantEDGE-node-ip-address>:31112`for second instance of demo 3

## File Structure

A quick look at the top-level relevant files and directories in demo 3 project.

    .
    ├── backend
         ├── server
         ├── api
         ├── util
               ├── mep1.yaml
               ├── mep2.yaml
         ├── go.mod
         ├── main.go
    ├── frontend
         ├── node_modules
         ├── src
              ├── js
              ├── css
              ├── .env
    ├── client
          ├── node_modules 
          
1.  **`/backend/server`**: This directory contains all code related to demo 3 backend

2.  **`/backend/api`**: This directory will contain all open-api specification for REST API of demo 3 

3.  **`./backend/util`**: This directory will contain configurations for an external demo 3 application running on MEC sandbox

4.  **`./backend/util.demo3-config.yaml`**: This file will configurations for an external demo 3 application running on mep1

5.  **`./backend/main.go`**: This file is the entry for launching demo 3 backend
   
6.  **`/frontend/node_modules`**: This directory contains all of the modules of code that your frontend depends on (npm packages) are automatically installed.

7.  **`/frontend/src`**: This directory will contain all of the code related to what you will see on the front-end of your site (what you see in the browser) 

8.  **`/frontend/.env.local`**: This file is used to store Environmental Variables to tell frontend where to poll resource from 

9.  **`/client/node_modules`**: This directory contains all of the modules of code that your frontend client package depends on (npm packages) are automatically installed.

## Scenario composition

The scenario is composed of two instance of same MEC applications running on two different mec platform

MEC Application instance 1 on MEC Platform mep1 coverage area: zone01 (PoA:4g-macro-cell-1)
MEC Application instance 2 on MEC Platform mep2 coverage area: zone2, zone3 (PoA: 4g1 & wifi1) 

The scenario is composed of the following components:

- 2 instance of the same edge application: demo-app and demo-app2
- 1 MNO that has 3 Zones
  - Zone1 has 1 Edge node
  - Zone2&3 has 1 Edge node
- 3 UEs with pre-defined route will move interchangeably from mep 1 to mep2

