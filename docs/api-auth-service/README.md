# Documentation for AdvantEDGE Auth Service REST API

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *http://localhost/auth/v1*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*AuthApi* | [**authenticate**](Apis/AuthApi.md#authenticate) | **GET** /authenticate | Authenticate service request
*AuthApi* | [**authorize**](Apis/AuthApi.md#authorize) | **GET** /authorize | OAuth authorization response endpoint
*AuthApi* | [**login**](Apis/AuthApi.md#login) | **GET** /login | Initiate OAuth login procedure
*AuthApi* | [**loginSupported**](Apis/AuthApi.md#loginsupported) | **GET** /loginSupported | Check if login is supported
*AuthApi* | [**loginUser**](Apis/AuthApi.md#loginuser) | **POST** /login | Start a session
*AuthApi* | [**logout**](Apis/AuthApi.md#logout) | **GET** /logout | Terminate a session
*AuthApi* | [**triggerWatchdog**](Apis/AuthApi.md#triggerwatchdog) | **POST** /watchdog | Send heartbeat to watchdog


<a name="documentation-for-models"></a>
## Documentation for Models

 - [Sandbox](./Models/Sandbox.md)


<a name="documentation-for-authorization"></a>
## Documentation for Authorization

All endpoints do not require authorization.
