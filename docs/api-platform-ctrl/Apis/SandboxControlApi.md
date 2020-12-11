# SandboxControlApi

All URIs are relative to *http://localhost/platform-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createSandbox**](SandboxControlApi.md#createSandbox) | **POST** /sandboxes | Create a new sandbox
[**createSandboxWithName**](SandboxControlApi.md#createSandboxWithName) | **POST** /sandboxes/{name} | Create a new sandbox
[**deleteSandbox**](SandboxControlApi.md#deleteSandbox) | **DELETE** /sandboxes/{name} | Delete a specific sandbox
[**deleteSandboxList**](SandboxControlApi.md#deleteSandboxList) | **DELETE** /sandboxes | Delete all active sandboxes
[**getSandbox**](SandboxControlApi.md#getSandbox) | **GET** /sandboxes/{name} | Get a specific sandbox
[**getSandboxList**](SandboxControlApi.md#getSandboxList) | **GET** /sandboxes | Get all active sandboxes


<a name="createSandbox"></a>
# **createSandbox**
> Sandbox createSandbox(config)

Create a new sandbox

    Create a new sandbox with a server-generated name

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **config** | [**SandboxConfig**](../Models/SandboxConfig.md)| Sandbox configuration information |

### Return type

[**Sandbox**](../Models/Sandbox.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="createSandboxWithName"></a>
# **createSandboxWithName**
> createSandboxWithName(name, config)

Create a new sandbox

    Create a new sandbox using provided name

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**| Sandbox name | [default to null]
 **config** | [**SandboxConfig**](../Models/SandboxConfig.md)| Sandbox configuration information |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

<a name="deleteSandbox"></a>
# **deleteSandbox**
> deleteSandbox(name)

Delete a specific sandbox

    Delete the sandbox with the provided name

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**| Sandbox name | [default to null]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

<a name="deleteSandboxList"></a>
# **deleteSandboxList**
> deleteSandboxList()

Delete all active sandboxes

    Delete all active sandboxes

### Parameters
This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

<a name="getSandbox"></a>
# **getSandbox**
> Sandbox getSandbox(name)

Get a specific sandbox

    Get sandbox information for provided sandbox name

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**| Sandbox name | [default to null]

### Return type

[**Sandbox**](../Models/Sandbox.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getSandboxList"></a>
# **getSandboxList**
> SandboxList getSandboxList()

Get all active sandboxes

    Returns a list of all active sandboxes

### Parameters
This endpoint does not need any parameter.

### Return type

[**SandboxList**](../Models/SandboxList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

