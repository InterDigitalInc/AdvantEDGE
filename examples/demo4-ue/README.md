# Demo4
Demo4 scenario showcases the _Edge Platform Application Enablement_ and _Application Mobility_ edge services.

For more details, check out the [Demo4 Documentation](https://interdigitalinc.github.io/AdvantEDGE/docs/usage/usage-demo4/)



Demo4 scenario showcases the _Edge Platform Application Enablement_ and _Application Mobility_ edge services.

Demo4 includes:
- A Terminal application, named [demo4-ue](#demo4-ue), that can be deployed either as a container using the provided AdvantEDGE scenario, or as an external application that interacts with private or public AdvantEDGE deployments such as the ETSI MEC Sandbox.
- An onboarded MEC application, named [onboarded-demo4](#onboarded-demo4) , that is deploy on the MEC platform. It can be instantiated using MEC-016 Service APIs

Demo4 Terminal application does not provide a dashboard GUI, tools such as cUrl or Postman can be used (see clause [demo4-ue](#demo4-ue))

## Demo4 Scenario Overview

The Demo4 scenario consists of one or more instances of a single Terminal application running one or more fixed or mobile terminal equipment. 
MEC-016 Service APIs "/dev_app/v1/context" and "/dev_app/v1/context/{instance}" provide the capability to instantiate one or more MEC applications.
By default, there is one onboarded MEC application named onboarded-demo4, described below (see file ~/AdvantDEGE/examples/demo4-ue/src/onboarded-demo4.yaml):

```yaml
{
    "appList":
        [
            {
                "appInfoList":
                    [
                        {
                            "appDId": "onboarded-demo4",
                            "appName": "onboarded-demo4",
                            "appProvider": "ETSI",
                            "appSoftVersion": "v0.1.0",
                            "appDVersion": "v0.1.0",
                            "appDescription": "Basic HTTP Ping Pong",
                            "appLocation":
                                [
                                    {
                                        "area": null,
                                        "civicAddressElement": null,
                                        "countryCode": "33"
                                    }
                                ],
                            "appCharcs":
                                [
                                    {
                                        "memory": 1024,
                                        "storage": 1024,
                                        "latency": 1024,
                                        "bandwidth": 1024,
                                        "serviceCont": 0
                                    }
                                ],
                            "cmd": "/onboardedapp/onboarded-demo/onboarded-demo4",
                            "args":null
                        }
                    ],
                    "vendorSpecificExt": {
                        "vendorId": "ETSI"
                    }
            }
        ]
}
```

### demo4-ue

The Terminal application demo4-ue provide a simply HTTP REST APIs for MEC-016 in order to validate quickly MEC-016 support.

#### Check that the demo4-ue is up

The request `GET /info/application` provides a description of the demo4-ue application.
```sh
$ curl "http://mec-platform.etsi.org:31111/info/application" -H "accept: application/json"
{"config":"app_instance.yaml","ip":"http://demo4-ue1:80","id":"ce454f4c-16f8-4542-83d8-c8afd45bcfea","mecReady":true,"subscriptions":{"AppTerminationSubscription":{"subId":"sub-LzNr62ED50GbZL_S"},"SerAvailabilitySubscription:":{"subId":"sub-Ec2tT-tKAe_qLxI7"}},"offeredService":{"serName":"demo4","id":"2ad0e566-2054-4550-ba97-1fe1f100fff0","state":"ACTIVE","scopeOfLocality":"MEC_SYSTEM","consumedLocalOnly":true},"discoveredServices":[{"serName":"meep-dai","serInstanceId":"70bb83a5-7ec7-4ab4-b889-ca7033a5be2b","consumedLocalOnly":true,"link":"http://mec-platform.etsi.org/usersb/dev_app/v1/","version":"2.0"},{"serName":"meep-rnis","serInstanceId":"3e72aa67-ff20-4e7e-bf21-aea928a1743b","consumedLocalOnly":true,"link":"http://mec-platform.etsi.org/usersb/rni/v2/","version":"2.0"},{"serName":"meep-ams","serInstanceId":"99e8a85d-5f87-48ed-be5e-1302c18e4664","consumedLocalOnly":true,"link":"http://mec-platform.etsi.org/usersb/amsi/v1/","version":"2.0"},{"serName":"meep-loc-serv","serInstanceId":"85ddf74d-d5e7-4d6f-beb5-25e97f8b57c9","consumedLocalOnly":true,"link":"http://mec-platform.etsi.org/usersb/location/v2/","version":"2.0"},{"serName":"meep-vis","serInstanceId":"3e23a210-8b1a-4a01-bc0c-1df7d1ad6be3","consumedLocalOnly":true,"link":"http://mec-platform.etsi.org/usersb/vis/v2/","version":"2.0"},{"serName":"meep-wais","serInstanceId":"a22f614e-a5c1-4902-84d3-4361ac9e3425","consumedLocalOnly":true,"link":"http://mec-platform.etsi.org/usersb/wai/v2/","version":"2.0"}]}
```

The request `GET /info/logs` provides some executon logs for the demo4-ue application.
```sh
$ curl "http://mec-platform.etsi.org:31111/info/logs" -H "accept: application/json"
["5. demo4DaiAppListGET: applicationList succeed, len= 1","4. daiClient instance created","3. Subscribe to service-availability notification [201]","2. Subscribe to app-termination notification [201]","1. === Register Demo4 MEC Application [200]","0. Send confirm ready [204]"]
```

#### Retrieve the list of existing onboarded MEC applications

The request `GET /dai/apps` provides the list of the onboarded MEC application.
```sh
$ curl "http://mec-platform.etsi.org:31111/dai/apps" -H "accept: application/json"
{"appList":[{"appInfo":{"appCharcs":{"bandwidth":1024,"latency":1024,"memory":1024,"storage":1024},"appDId":"onboarded-demo4","appDVersion":"v0.1.0","appDescription":"Basic HTTP Ping Pong","appLocation":[{"countryCode":"33"}],"appName":"onboarded-demo4","appProvider":"ETSI","appSoftVersion":"v0.1.0"}}]}
```

#### Instantiate an onboarded MEC application

The request `POST /dai/instantiate` is used to create a new AppContext and instantiate the specified onboarded MEC application.
```sh
$ curl -X POST "http://mec-platform.etsi.org:31111/dai/instantiate" -H "accept: application/json"
{"appAutoInstantiation":true,"appInfo":{"appDId":"onboarded-demo4","appDVersion":"v0.1.0","appDescription":"Basic HTTP Ping Pong","appName":"onboarded-demo4","appProvider":"ETSI","appSoftVersion":"v0.1.0","appPackageSource":"appPackageSource1","userAppInstanceInfo":[{"appInstanceId":"20","appLocation":{"countryCode":"33"},"referenceURI":"https://mec-platform.etsi.org/usersb/onboarded-demo4"}]},"appLocationUpdates":true,"associateDevAppId":"04e71585-87c7-4d2e-b913-538cf1ef","callbackReference":"http://demo4-ue1:80","contextId":"20"}
```

Note that, in the response:
- The field "contextId" indicates the identifier to be used to [terminate the instance](#terminate-an-existing-instance-of-onboarded-mec-application).
- The field "referenceURI" indicates the URL to the new instance of the onboarded MEC application, see clause [Ping request](#ping-request]).

#### Ping request

The request `GET /onboarded-demo4/ping/{appcontextid}` is the ping/pong request to send to the instance of the onboarded MEC application [onboarded-demo4](#onboarded-demo4).
```sh
$ curl "http://mec-platform.etsi.org/usersb/onboarded-demo4/ping/20" -H "accept: application/json"
"pong"
```

#### Terminate an existing instance of onboarded MEC application

The request `DELETE /dai/delete/{appcontextid}` terminates the instance of the specified onboarded MEC application.

```sh
$ curl -X DELETE "http://mec-platform.etsi.org:31111/dai/delete/20" -H "accept: application/json"
```

#### application location availability task

The request `POST /dai/availability/{appcontextid}` provides the location availability for the specified instance of the onboarded MEC application.

```sh
$ curl -X POST "http://mec-platform.etsi.org:31111/dai/availability/20" -H "accept: application/json"
```

#### Postman workspace

The code below is the the Postman workspace equivalent to the curl requests above.

```json
{
	"info": {
		"_postman_id": "e8a79bf9-aeba-4869-91a1-45fb6f7320a7",
		"name": "demo4-ue",
		"description": "This collection provides all Curl command supported by the UE appliction demo4-ue",
		"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		"_exporter_id": "3317208"
	},
	"item": [
		{
			"name": "/info/application",
			"protocolProfileBehavior": {
				"disabledSystemHeaders": {
					"accept": true
				}
			},
			"request": {
				"method": "GET",
				"header": [
					{
						"key": "Accept",
						"value": "application/json",
						"type": "text"
					}
				],
				"url": {
					"raw": "try-mec.etsi.org:31111/info/application",
					"host": [
						"try-mec",
						"etsi",
						"org"
					],
					"port": "31111",
					"path": [
						"info",
						"application"
					],
					"query": [
						{
							"key": "",
							"value": null,
							"disabled": true
						}
					]
				}
			},
			"response": []
		},
		{
			"name": "/info/logs",
			"protocolProfileBehavior": {
				"disabledSystemHeaders": {
					"accept": true
				}
			},
			"request": {
				"method": "GET",
				"header": [
					{
						"key": "Accept",
						"value": "application/json",
						"type": "text"
					}
				],
				"url": {
					"raw": "try-mec.etsi.org:31111/info/logs",
					"host": [
						"try-mec",
						"etsi",
						"org"
					],
					"port": "31111",
					"path": [
						"info",
						"logs"
					],
					"query": [
						{
							"key": "",
							"value": null,
							"disabled": true
						}
					]
				}
			},
			"response": []
		},
		{
			"name": "/dai/apps",
			"protocolProfileBehavior": {
				"disabledSystemHeaders": {
					"accept": true
				}
			},
			"request": {
				"method": "GET",
				"header": [
					{
						"key": "Accept",
						"value": "application/json",
						"type": "text"
					}
				],
				"url": {
					"raw": "try-mec.etsi.org:31111/dai/apps",
					"host": [
						"try-mec",
						"etsi",
						"org"
					],
					"port": "31111",
					"path": [
						"dai",
						"apps"
					],
					"query": [
						{
							"key": "",
							"value": null,
							"disabled": true
						}
					]
				}
			},
			"response": []
		},
		{
			"name": "/dai/instantiate",
			"protocolProfileBehavior": {
				"disabledSystemHeaders": {
					"accept": true
				}
			},
			"request": {
				"method": "POST",
				"header": [
					{
						"key": "Accept",
						"value": "accept: application/json",
						"type": "text"
					}
				],
				"url": {
					"raw": "try-mec.etsi.org:31111/dai/instantiate",
					"host": [
						"try-mec",
						"etsi",
						"org"
					],
					"port": "31111",
					"path": [
						"dai",
						"instantiate"
					]
				}
			},
			"response": []
		},
		{
			"name": "/usersb/onboarded-demo4/ping",
			"protocolProfileBehavior": {
				"disabledSystemHeaders": {
					"accept": true
				}
			},
			"request": {
				"method": "GET",
				"header": [
					{
						"key": "Accept",
						"value": "accept: application/json",
						"type": "text"
					}
				],
				"url": {
					"raw": "https://195.238.226.94/usersb/onboarded-demo4/ping",
					"protocol": "https",
					"host": [
						"195",
						"238",
						"226",
						"94"
					],
					"path": [
						"usersb",
						"onboarded-demo4",
						"ping"
					]
				}
			},
			"response": []
		},
		{
			"name": "/usersb/onboarded-demo4/terminate",
			"protocolProfileBehavior": {
				"disabledSystemHeaders": {
					"accept": true
				}
			},
			"request": {
				"method": "DELETE",
				"header": [
					{
						"key": "Accept",
						"value": "accept: application/json",
						"type": "text"
					}
				],
				"url": {
					"raw": "https://try-mec.etsi.org/usersb/onboarded-demo4/ping",
					"protocol": "https",
					"host": [
						"try-mec",
						"etsi",
						"org"
					],
					"path": [
						"usersb",
						"onboarded-demo4",
						"ping"
					]
				}
			},
			"response": []
		},
		{
			"name": "/dai/doping",
			"protocolProfileBehavior": {
				"disabledSystemHeaders": {
					"accept": true
				}
			},
			"request": {
				"method": "GET",
				"header": [
					{
						"key": "Accept",
						"value": "accept: application/json",
						"type": "text"
					}
				],
				"url": {
					"raw": "try-mec.etsi.org:31111/dai/doping/25",
					"host": [
						"try-mec",
						"etsi",
						"org"
					],
					"port": "31111",
					"path": [
						"dai",
						"doping",
						"25"
					]
				}
			},
			"response": []
		},
		{
			"name": "/dai/delete",
			"protocolProfileBehavior": {
				"disabledSystemHeaders": {
					"accept": true
				}
			},
			"request": {
				"method": "DELETE",
				"header": [
					{
						"key": "Accept",
						"value": "accept: application/json",
						"type": "text"
					}
				],
				"url": {
					"raw": "try-mec.etsi.org:31111/dai/delete/25",
					"host": [
						"try-mec",
						"etsi",
						"org"
					],
					"port": "31111",
					"path": [
						"dai",
						"delete",
						"25"
					]
				}
			},
			"response": []
		},
		{
			"name": "/dai/availability",
			"protocolProfileBehavior": {
				"disabledSystemHeaders": {
					"accept": true
				}
			},
			"request": {
				"method": "POST",
				"header": [
					{
						"key": "Accept",
						"value": "accept: application/json",
						"type": "text"
					}
				],
				"url": {
					"raw": "try-mec.etsi.org:31111/dai/availability/25",
					"host": [
						"try-mec",
						"etsi",
						"org"
					],
					"port": "31111",
					"path": [
						"dai",
						"availability",
						"25"
					]
				}
			},
			"response": []
		}
	]
}
```

### onboarded-demo4

The MEC application onboarded-demo4 is a basic HTTP REST API "Ping/Pong" application. When it recieves an HTTP GET /ping request, it replies with 200 OK "pong".
The goal of this application is just to validate the support of MEC-016 Service APIs.

## Using Demo4

### Preamble (reserved to MEC Sandbox administrator)

Before using Demo4, the configuration files '~/AdvantDEGE/examples/demo4-ue/src/onboarded-demo4.yaml' and '~/AdvantEDGE/examples/demo4-ue/src/onboarded-demo/onboarded-demo-test1.json' and the binariy folders shall be copied into the folder /var/lib/docker/volumes/meep-dai/_data/ (required to use sudo command). This is the descriptor of the onboarded MEC application [onboarded-demo4](#onboarded-demo4).

**_Note_** For the MEC Sandbox platform, the configuration files and the binariy folders shall be copied into the folder /var/lib/docker/volumes/<sandbox network name>/_data (e.g. /var/lib/docker/volumes/meep-4g-5g-wifi-macro-mec016-1/_data)

### Using Demo4 with AdvantEDGE

To use Demo4 as an AdvantEDGE scenario container:

- Build & dockerize Demo4 server & frontend
- Import the provided scenario demo4-scenario.yaml
- Create a sandbox & deploy Demo4 scenario
- Start Demo4 application frontend in browser

#### Build from source
To build demo4-ue & onboarded-demo4 from source code:
```sh
$ cd ~/AdvantEDGE/examples/demo4-ue 
$ ./build-demo4-ue.sh --rebuild_dai
```

NOTE: Binary files are created in ./bin/ folder.

#### Dockerize demo applications
Demo Application binaries must be dockerized (containerized) as container images in the Docker registry. This step is necessary every time the demo binaries are updated.

NOTE: Make sure you have deployed the AdvantEDGE dependencies (e.g. docker registry) before dockerizing the demo binaries.

To generate docker images from demo binary files:

```sh
$ cd ~/AdvantEDGE/examples/demo4-ue
$ ./dockerize.sh
```

### Using Demo4 with ETSI MEC Sandbox

To use Demo4 as an external application that interacts with the ETSI MEC Sandbox

- Build Demo4 server & frontend
- Log in to the [ETSI MEC Sandbox](https://try-mec.etsi.org)
- Deploy either of the dual-mep scenarios
- Configure Demo4 application instances
- Start Demo4 application instances
- Demo4 does not have prior knowledge or configuration information of the MEC services offered by the MEC platform.

Therefore, the following steps need to be done prior to running Demo4 application instances.

#### Obtain demo binaries
- Use the same procedure described above for Demo4 with AdvantEDGE.

#### Create work directory for demo4-ue instance
- Create a work directory of your choice on the system (e.g. ~/tmp/demo4) and copy the files ~/AdvantEDGE/examples/demo4-ue/src/demo-server/demo4-ue-config.yaml and ~/AdvantEDGE/examples/demo4-ue/bin/demo-server/demo-server.

	The structure should look like this:

	```
	~/tmp/demo4
			|
			\____ demo4-ue-config.yaml
			|
			\____ demo-server
	```

Update the configuration file demo4-ue-config.yaml accordingly and launch the demo application:

```sh
$ ./demo-server ./demo4-ue-config.yaml
```