# Demo3

Demo 3 is used to showcase platform capability of 
Application Mobility Service (AMS) API and Edge Application Enablement Service APIs

## Introduction

Edge Application Enablement Service APIs (MEC011) allow MEC applications to:
1. Discover & offer MEC platform services
2. Subscribe for application lifecycle and MEC platform service availability notifications

MEC Applications should use the Application Mobility Service (AMS) API to perform MEC-assisted application context transfers.
Using the AMS API, MEC Applications can register for simulated terminal device Mobility Procedure notifications to trigger user application context transfers across MEC Applications

## Getting Started

There are two ways to run Demo3:
1. Running Demo3 as an external MEC Application with requirement to configure details of appInstanceId that has to be generated from the MEC Sandbox frontend along with MEC Service base path. A scenario needs to be deployed on MEC Sandbox prior to running Demo3
2. Use a pre-defined demo 3 scenario to deploy Demo3 application that runs on AdvantEdge as a container within a kubernetes environment 

## Running Demo 3 as an external application with MEC Platform

The following steps need to be done prior to running demo 3

| Operation: | Notes: |
| ---------  | ------ |
| 1. Login via the MEC Sandbox frontend |  |
| 2. Select a network to deploy in the user sandbox |  |
| 3a. Pre-configure `mep1.yaml` or `mep2.yaml` located in Demo 3 backend directory refer to - [File Structure](#file-structure)  | Add configuration fields <br> `mode:`  "sandbox" |
| 3b. Configure `mep1.yaml` sandbox field with MEC Service base path | Base path format `https://<my-mec-url>/<my-sandbox-key>/mep1` mep1 is an indicator to refer to an edge application running on mep1. MEC Application will learn MEC services availability via mep1 interface
| 3c. Configure `mep1.yaml` mecplatform field with which platform the demo 3 application will be running on | If demo3 will be running on mep1, it would make sense to use `mep1.yaml` as your demo3 configuration file <br> `mecplatform: 'mep1' `| 
| 3d. Obtain an application instance ID using the MEC application startup procedure and configure `mep1.yaml` appid field with application instance ID running on mep1 | `appid: '<app-instance-id>'` | 
| 3e. Pre-configure `localurl` with your I.P address and `port` to indicate port number that demo3 server will run at | `localurl: 'http://<my-ip-address>'` <br> `port: ':<my-port-number>'` 
| 4. Optional: If running a dual mep scenario on MEC sandbox. The above steps needs to be repeated to run a seperate instance of demo 3 application by applying configurations into `mep2.yaml` |  | 

configuration will look like this at the end:
```sh .env
mode: 'sandbox'
sandbox: `https://<my-mec-url>/<my-sandbox-key>/mep1`
mecplatform: 'mep1'
appid: '44a0a575-916d-4cac-874c-514833dc3035'
localurl: 'http://<my-ip-address>'
port: ':8093'
```

### To Build Demo 3 Server

```shell
# Build demo 3 backend binary 
cd AdvantEdge/examples/demo3/src/backend
go build -o demo-server .
go run demo-server mep1.yaml
```

### To Build Demo 3 Frontend 
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
4. Set up your .env.local file with values of where demo 3 backend is served  
```
cd AdvantEdge/examples/demo3/src
# Create a file called .env.local
# Apply configurations 
URL=http://<my-ip-address>:<my-demo3-port-number>
```
5. Run demo 3 in development mode or build into binaries
```
npm run build:dev
npm run build
```

## Running Demo 3 scenario on AdvantEdge

#### Demo server & Demo frontend 

Demo server is a web server that will run internally in AdvantEdge as the backend of demo 3 application. The frontend is a dashboard to provide information on the MEC application instance in the scenario. One capability When the external UE moves in the network and transitions from one edge instance to another, the "UE state" (e.g. the counter value) is transferred using the application state transfer. On the demo server, it maintains a state of external UE as a counter that will continue incrementing (e.g. not reset to zero) when the UE moves in the network.

## Using the scenario

The following steps need to be done prior to using this scenario

#### Obtain demo binaries

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

#### Open demo 3 scenario in the AdvantEdge configure tab & save

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
              ├── .env.local
    ├── client
          ├── node_modules 
          
1.  **`/backend/server`**: This directory contains all code related to demo 3 backend

2.  **`/backend/api`**: This directory will contain all open-api specification for REST API of demo 3 

3.  **`./backend/util`**: This directory will contain configurations for an external demo 3 application running on MEC sandbox

4.  **`./backend/util.mep1.yaml`**: This file will configurations for an external demo 3 application running on mep1

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

