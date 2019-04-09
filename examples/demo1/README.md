# Demo1
A simple scenario used to showcase platform capabilities.

This has two applications _iperf server_ and _demo server_ deployed across multiple tiers of the network (fog/edge/cloud).

It has two clients (internal & external) that communicate with the servers.
Internal client traffic is iperf only and has no GUI.
External client accesses both iperf and demo servers and has a GUI.

The platform capabilities demonstrated with this scenario are:
- Scenario deployment (dynamic charts)
- Network tiering
- Internal & external clients
- Network characteristics
- Edge-application deployment model: one-to many relationship
- UE mobility
- Applciation state transfer (_demo server_)
- Monitoring (demo specific dashboards)

## Scenario composition
The scenario is composed of the following components:

- 2 distant cloud application: _iperf_ server and _demo_ server
- 1 MNO that has 2 Zones
  - Zone1 has 2 PoAs & 1 Edge node
  - Zone2 has 1 PoA & 1 Edge node
  - PoA1 in Zone1 is equipped with a Fog node
  - Each Fog/Edge node runs 2 Edge servers (_iperf_ and _demo_)
- 2 UEs
  - 1 internal UE that runs an iperf client
  - 1 external UE that runs a Demo client

#### Internal UE application
Upon scenario startup, internal UE application (an iperf client) connects automatically to the closest iperf server and starts transferring traffic.

As the UE moves around the network, edge node instance will change.

#### External UE
External UE application is a javascript application running in an external browser.

To start the aooplcation, load the following page in the browser `<AdvantEDGE-node-ip-address>:31111`

The application shows details about the connection, allows to start a state counter and iperf traffic and presents an image. See Iperf & Demo server sub-sections for more details.

#### Iperf server
This is a standard iperf server that will terminate iperf client connections.

There is an iperf client running in the internal UE and another one in the external UE.

External UE needs the iperf proxy running to be able to control the iperf client from the javascript GUI.

#### Demo server
Demo server is a web server that maintains a UE state and also stores unique data.
Only the external UE accesses the demo server.

Unique data is an image that is different for each location; therefore depending on the UE location in the network, the external UE GUI will show a different image. This image provides a visual indication that the UE is connected to a different network location.

UE state is in the form of a counter incremented on the server every second. Counter lives on the server side and UE can only start/reset the counter and display its value.
On the UE GUI, the counter is started by pressing the button.

When the external UE moves in the network and transitions from one edge instance to another, the "UE state" (e.g. the counter value) is transferred using the application state transfer. On the UE GUI, the counter continue incrementing (e.g. not reset to zero) when the UE moves in the network.

## Using the scenario
The following steps need to be done prior to using this scenario

#### Dockerize demo applications
Only need to do it once, or when the demo application changes

In a command line shell
```
cd ~/AdvantEDGE/examples/demo1/
./dockerize.sh
```
> Demo Application binaries are dockerized (containerized) and the container images are stored in the local Docker registry.<br> Next time you want to use the demo scenario; demo application containers will be available

#### Configure demo specific dashboards
Only need to do it once, or when the demo dashboard changes

Follow the procedure described in [Configuring Monitoring](../../docs/use/monitoring.md#configure-dashboards) and import the demo specific dashboard from `AdvantEDGE/examples/demo1/demo1-dashboards.json`

> Demo specific dashboards are stored in Kibana.<br> Next time you want to use the demo scenario; demo specific dashboard will be available.

#### Start iperf proxy
Do it everytime you start using the demo when the iperf-proxy is not running

This demo scenario requires iperf installed on the AdvantEDGE host and the iperf proxy running.

If `which iperf` returns nothing, install iperf
```
sudo apt-get install iperf
# the following should now return /usr/bin/iperf
which iperf
```

Start iperf proxy, in a command line shell
```
cd ~/AdvantEDGE/examples/demo1/bin/iperf-proxy
./iperf-proxy
```
