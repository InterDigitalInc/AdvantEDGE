---
layout: default
title: Demo 3 Scenario
parent: Usage
nav_order: 7
---

## Demo3
Demo3 scenario showcases the _Edge Platform Application Enablement_ and _Application Mobility_ edge services.

Demo3 includes a user MEC application that can be deployed either as a container using the provided AdvantEDGE scenario,
or as an external application that interacts with private or public AdvantEDGE deployments such as the [ETSI MEC Sandbox](https://try-mec.etsi.org/).

Demo3 MEC Application provides a dashboard GUI for visualizing MEC Service interactions such as:
- MEC011: _Application Support & Service Enablement_
  - MEC Application registration/deregistration
  - MEC Service discovery & offering
  - Event notifications for service and application availability
- MEC021: _Application Mobility Service_
  - User context relocation request registration
  - MEC Assisted state transfer subscriptions & notifications

## Demo3 Scenario Overview
The Demo3 scenario consists of two instances of a single MEC application running on two different mec platforms
- MEC Application instance 1 on MEC Platform _mep1_ with coverage area: _zone1_
- MEC Application instance 2 on MEC Platform _mep2_ with coverage area: _zone2_

The scenario includes:
- 2 instance of the same edge application: _demo3-mep1_ and _demo3-mep2_
- 1 MNO that has 2 Zones
  - Zone1 has 1 Edge node
  - Zone2 has 2 Edge node
- 3 terminals with pre-defined routes to move interchangeably from _mep1_ to _mep2_

**NOTE:** To view terminal movement on a map, you must first provision maps in AdvantEDGE as described
[here]({{site.baseurl}}{% link docs/overview/features/overview-gis.md %}#map-provisioning)

### Demo server
Demo server is a web server that interacts with the _Application Enablement_ & _Application Mobility_ services
and maintains MEC Application instance state such as:
- Application readiness
- Edge service availability
- Mobility & termination subscriptions
- Terminal device contexts (counters that continue to increment even with MEC application mobility)

### Demo frontend
The frontend is an instance-specific dashboard that provides MEC application information.

It provides controls to:
- Register/Deregister the MEC application instance
- Add/Remove terminal devices to track for Edge application mobility

## Using Demo3 with AdvantEDGE
To use Demo3 as an AdvantEDGE scenario container:
- Build & dockerize Demo3 server & frontend
- Import the provided scenario _demo3-scenario.yaml_
- Create a sandbox & deploy Demo3 scenario
- Start Demo3 application frontend in browser

### Build from source
To build _demo-frontend_ & _demo-server_ from source code:

```
cd ~/AdvantEDGE/examples/demo3
./build-demo3.sh
``` 

_**NOTE:** Binary files are created in ./bin/ folder_

### Optionally use pre-built binaries (from GitHub release)
```
# Get bin folder tarball from desired release
cd ~/AdvantEDGE/examples/demo3
tar -zxvf demo3.<version>.linux-amd64.tar.gz
```

### Dockerize demo applications
Demo Application binaries must be dockerized (containerized) as container images in the Docker registry.
This step is necessary every time the demo binaries are updated.

_**NOTE:** Make sure you have deployed the AdvantEDGE dependencies (e.g. docker registry) before dockerizing the demo binaries._

To generate docker images from demo binary files:

```
cd ~/AdvantEDGE/examples/demo3
./dockerize.sh
```

### Deploy Demo3 scenario
After building & dockerizing the Demo3 application, the Demo3 scenario must be imported in AdvantEDGE and deployed within a sandbox.

_**IMPORTANT NOTE:** For the API calls to work correctly, the platform IP address or domain name must be properly configured in the [AdvantEDGE deployment configuration]({{site.baseurl}}{% link docs/platform-mgmt/mgmt-cheat-sheet.md %}#deployment-configuration) file._

### Start Demo3 application frontend
After deploying the Demo3 scenario, the Demo3 application frontend can be accessed as follows:
- _demo3-mep1_: `http://<AdvantEDGE IP address>:31111`
- _demo3-mep2_: `http://<AdvantEDGE IP address>:31112`

Each Demo3 application instance provides a pre-configured frontend with controls to register/deregister the application instance to MEC011 and to track terminal devices using MEC021.

## Using Demo3 with ETSI MEC Sandbox
To use Demo3 as an external application that interacts with the ETSI MEC Sandbox
- Build Demo3 server & frontend
- Log in to the [ETSI MEC Sandbox](https://try-mec.etsi.org/)
- Deploy either of the _dual-mep_ scenarios
- Configure Demo3 application instances
- Start Demo3 application instances

Demo3 does not have prior knowledge or configuration information of the MEC services offered by the MEC platform.

Therefore, the following steps need to be done prior to running Demo3 application instances.

### Obtain demo binaries
Use the same procedure described above for Demo3 with AdvantEDGE.

### Create work directories for each Demo3 MEC application instance
Create work directories of your choice on the system of your choosing; we'll use `~/demo3-mep1` and `~/demo3-mep2` for
this example and create a folder named `static` inside each one of the folders.

The structure should look like this:
```
     ├── demo3-mep1
          ├── static
     ├── demo3-mep2
          ├── static
```

**IMPORTANT: _For this demo to work, the system running demo applications must be at a public IP address so that
notifications sent by the ETSI MEC Sandbox can be received by demo applications. If the system is behind a firewall,
ports will need to be opened._**

### Copy demo-server to working directories
For each application instance, copy the following files to the working directories:
- Server executable (`/AdvantEDGE/examples/demo3/bin/demo-server/demo-server`)
- Template configuration file (`/AdvantEDGE/examples/demo3/demo3-config.yaml`)
- Frontend bundle `/AdvantEDGE/example/demo3/bin/demo-frontend/*` to the static folder

The resulting structure should look like this:
```
     ├── demo3-mep1
          ├── demo-server
          ├── demo3-config.yaml
          ├── static
               ├── bundle.css
               ├── bundle.js
               ├── img
               ├── index.html
     ├── demo3-mep2
          ├── demo-server
          ├── demo3-config.yaml
          ├── static
               ├── bundle.css
               ├── bundle.js
               ├── img
               ├── index.html
```

### Create application instances in the MEC Sandbox
Login via the [ETSI MEC Sandbox](https://try-mec.etsi.org/) frontend.

Deploy either of teh dual-MEP scenarios. Note that the _dual-mep-short-path_ network scenario will trigger
AMS mobility procedure notifications much quicker that the _dual-mep-4g-5g-wifi-macro_ network scenario.

Create two Application Instance IDs called `demo3`, one on _mep1_ and one on _mep2_ respectively.

### Edit application instance configuration files
For each application instance, edit the copied template configuration file with instance-specific configuration values.

The resulting configuration should look like this:
```
mode: 'sandbox'             # demo3 runs against ETSI MEC Sandbox
https: true                 # ETSI MEC sandbox uses https
sandbox: 'https://try-mec.testfqdn.dev/mep1'  # URL to your sandbox, this info is available in the ETSI MEC Sandbox frontend
mecplatform: 'mep1'         # the MEP where the instance is running, one of your application should be mep1 and the other mep2
appid: ''                   # these are created in the ETSI MEC Sandbox frontend
localurl: 'http://'         # the public IP address where demo3 instance is running
port: ''                    # the port number that demo3 is listening on for incoming traffic
```

### Start Demo3 instances
Start the demo3 instances `./demo3-server demo3-config.yaml`

After starting the servers, the frontend can be accessed at `<your-ip-address>:<your-port>`

From the frontend, Demo3 can register to MEC011 and then devices present in the scenario can be added.
