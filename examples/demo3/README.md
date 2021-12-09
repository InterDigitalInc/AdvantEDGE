# Demo3

Demo 3 demonstrates how to use AdvantEDGE mec services showcasing Application Mobility Service (AMS) and Edge Application Enablement Service API.

Demo 3 is a user MEC application that can be integrated within an AdvantEDGE scenario or used as an external application with [ETSI MEC Sandbox](https://try-mec.etsi.org/). 




Demo3 MEC Application provides a dashboard GUI that allow a user to visualize interactions with MEC Service such as:

1. MEC011 - Application Enablement and Service Enablement Services:

- MEC Application registration/deregistration 
- MEC Service discovery & offering 
- Monitor event notifications about service and application availability 

2. Application Mobility Service provides support for relocation of user context and/or application instance between MEC hosts enabling:

- MEC Assisted state transfer operationn



## Getting Started:
Demo 3 can be used with an AdvantEDGE deployment or with the ETSI MEC Sandbox public deployment.
#### AdvantEDGE procedure overview
- Demo3 is built and dockerized using provided bash scripts
- Demo3 scenario (provided) is imported in AdvantEDGE
- Demo3 scenario is deployed on AdvantEDGE
#### MEC Sandbox procedure overview
- Demo3 is built using provided bash scripts
- Demo3 executable and frontend are copied to a common folder
- A MEC Sandbox is created and configured on ETSI MEC Sandbox [site](https://try-mec.etsi.org/)
- Demo3 configuration files (provided) are updated with your environment values
- Demo3 application is executed and accesses the ETSI MEC Sandbox

###  Using Demo3 with ETSI MEC Sandbox

- Demo 3 does not have any prior knowledge or configuration information of the MEC services offered by the MEC platform

The following steps need to be done prior to running demo 3

| Operation: | Notes: |
| ---------  | ------ |
| 0a. Build demo3 server and frontend by invoking `Advantedge/example/demo3/build-demo3.sh` |  |
| 0b. Create work directories of your choice on the system of your choosing; we'll use `~/demo3-mep1` and `~/demo3-mep2` for this example and create a folder named `static` inside each one of the folders. The structure should look like this <br>
     ├── demo3-mep1
          ├── static
     ├── demo3-mep2
          ├── static
**IMPORTANT** ``` For this demo to work, the system running demo applications must be at a public IP address so that notifications sent by the ETSI MEC Sandbox can be received by demo applications. If the system is behind a firewall, ports will need to be opened.  ``` 

|  |  |
| ---------  | ------ |
| 0c. For each application instance, copy Demo3 server (`/AdvantEDGE/example/demo3/bin/demo-server/demo-server`) in the work directories and copy Demo3 frontend bundle `/AdvantEDGE/example/demo3/bin/demo-frontend/*` in the static folder <br> The resulting should look like this <br>
     ├── demo3-mep1
          ├── demo-server
          ├── static
               ├── bundle.css
               ├── bundle.js
               ├── img
                    ├── AdvantEDGE-logo-NoTagline_White_RGB.png
                    ├── ID-Icon-01-idcc.svg
                    ├── network.png
               ├── index.html
     ├── demo3-mep2
          ├── demo-server
          ├── static
               ├── bundle.css
               ├── bundle.js
               ├── img
                    ├── AdvantEDGE-logo-NoTagline_White_RGB.png
                    ├── ID-Icon-01-idcc.svg
                    ├── network.png
               ├── index.html
|  |  |                  
| ---------  | ------ |
| 1. Login via the MEC Sandbox frontend |  |
| 2. Deploy a Dual MEP scenario and create two Application Instance IDs calles `demo3`, one on MEP1 and one on MEP2 respectively |  |
| 3. Using provided `demo3-config.yaml` template, create a configuration file for each of the MEP1 and MEP2 applications and save each configuration files in their respective work folder created in step 0 <br> The resulting structure should look like this <br>` mode: 'sandbox'` - demo3 runs against ETSI MEC Sandbox <br> ` https: true ` - ETSI MEC sandbox uses https<br>`sandbox: 'https://try-mec.testfqdn.dev/mep1' `- URL to your sandbox, this info is available in the ETSI MEC Sandbox frontend <br> `mecplatform: 'mep1'`- the MEP where the instance is running, one of your application should be mep1 and the other mep2<br> `appid: '' ` - these are created in the ETSI MEC Sandbox frontend<br> ` localurl: 'http://'` - the public IP address where demo3 instance is running<br> `port: ''` - the port number that demo3 is listening on for incoming traffic<br>
| 4. Start the demo3 instances `./demo3-server demo3-config.yaml` -- after starting the servers, the frontend can be accessed at `<your-ip-address>:<your-port>` <br> From the frontend, demo3 can register to MEC011 and then devices present in the scenario can be added.

Example of a configuration file
```sh .env
# This file defines the configuration of Demo3 edge application
# Required to define if application will run on MEC Sandbox or AdvantEDGE. Expected Values: sandbox or advantedge
mode: 'sandbox'
# Field to define MEC platform URL that app will run on. Example: http://{Mec_IP_Address}/{Sandbox_name}/{Mep_name}/
https: true
sandbox: 'https://try-mec.testfqdn.dev/sbx-1234abcd/mep1'
# Field to define MEC platform name. Example: mep1
mecplatform: 'mep1'
# Field to define user-application ID that can be generated using MEC Sandbox frontend
appid: 'cd9e4234-d7b7-4d49-be64-850ca436e24c'
# Local I.P address application will run on
localurl: 'http://1.2.3.4'
# Port number of user-application example: port: '8093'
port: '31111'

```

### Build Demo 3 Server and Frontend

``` shell
cd AdvantEdge/examples/demo3
./build-demo3.sh
``` 
###  Using Demo3 with AdvantEDGE

### Demo server & Demo frontend 

Demo server is a web server that will run internally in AdvantEdge as the backend of demo 3 application. The frontend is a dashboard to provide information on the MEC application instance in the scenario. One capability When the external UE moves in the network and transitions from one edge instance to another, the "UE state" (e.g. the counter value) is transferred using the application state transfer. On the demo server, it maintains a state of external UE as a counter that will continue incrementing (e.g. not reset to zero) when the UE moves in the network.



The following steps need to be done prior to using this scenario

#### Configure demo frontend and obtain binaries
| Operation | Notes |   
| ---------  | ------ |
| 0. Build and dockerize demo3 by invoking `AdvantEDGE/example/demo3/build-demo3.sh` and `AdvantEDGE/example/demo3/dockerize.sh`|  |
| 1. Import provided `demo3-scenario.yaml` in AdvantEDGE and save it | |
| 2. Deploy `demo3-scenario` from the frontend | this scenario uses geo-localization, therefore it is necessary to provision a map as described [here](https://interdigitalinc.github.io/AdvantEDGE/docs/overview/features/gis/#map-provisioning) |


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

