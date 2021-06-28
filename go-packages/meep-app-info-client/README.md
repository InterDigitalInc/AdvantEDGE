# Go API client for client

AdvantEDGE implementation to create an Application Instance information using OpenAPI. Developed as an extension to Application Enablement API.

## Overview
This API client was generated by the [swagger-codegen](https://github.com/swagger-api/swagger-codegen) project.  By using the [swagger-spec](https://github.com/swagger-api/swagger-spec) from a remote server, you can easily generate an API client.

- API version: 1.0.0
- Package version: 1.0.0
- Build package: io.swagger.codegen.v3.generators.go.GoClientCodegen

## Installation
Put the package under your project folder and add the following in import:
```golang
import "./client"
```

## Documentation for API Endpoints

All URIs are relative to *https://localhost/sandboxname/app_info/v1*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*AppsApi* | [**ApplicationsAppInstanceIdDELETE**](docs/AppsApi.md#applicationsappinstanceiddelete) | **Delete** /applications/{appInstanceId} | 
*AppsApi* | [**ApplicationsAppInstanceIdGET**](docs/AppsApi.md#applicationsappinstanceidget) | **Get** /applications/{appInstanceId} | 
*AppsApi* | [**ApplicationsAppInstanceIdPUT**](docs/AppsApi.md#applicationsappinstanceidput) | **Put** /applications/{appInstanceId} | 
*AppsApi* | [**ApplicationsGET**](docs/AppsApi.md#applicationsget) | **Get** /applications | 
*AppsApi* | [**ApplicationsPOST**](docs/AppsApi.md#applicationspost) | **Post** /applications | 


## Documentation For Models

 - [ApplicationInfo](docs/ApplicationInfo.md)
 - [ApplicationState](docs/ApplicationState.md)
 - [LocalityType](docs/LocalityType.md)
 - [ProblemDetails](docs/ProblemDetails.md)


## Documentation For Authorization
 Endpoints do not require authorization.


## Author

AdvantEDGE@InterDigital.com
