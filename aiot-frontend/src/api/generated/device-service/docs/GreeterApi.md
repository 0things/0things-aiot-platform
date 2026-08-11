# GreeterApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**greeterSayHello**](#greetersayhello) | **GET** /helloworld/{name} | |

# **greeterSayHello**
> HelloworldV1HelloReply greeterSayHello()

Sends a greeting

### Example

```typescript
import {
    GreeterApi,
    Configuration
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new GreeterApi(configuration);

let name: string; // (default to undefined)

const { status, data } = await apiInstance.greeterSayHello(
    name
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **name** | [**string**] |  | defaults to undefined|


### Return type

**HelloworldV1HelloReply**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

