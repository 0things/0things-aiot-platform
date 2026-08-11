# DeviceServiceApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**deviceServiceActivateDevice**](#deviceserviceactivatedevice) | **POST** /v1/devices/{id}/activate | |
|[**deviceServiceBatchUploadDevices**](#deviceservicebatchuploaddevices) | **POST** /v1/devices/batch/upload | |
|[**deviceServiceClearPushRecords**](#deviceserviceclearpushrecords) | **DELETE** /v1/devices/{deviceKey}/push-records | |
|[**deviceServiceCreateDevice**](#deviceservicecreatedevice) | **POST** /v1/devices | |
|[**deviceServiceDeleteDevice**](#deviceservicedeletedevice) | **DELETE** /v1/devices/{id} | |
|[**deviceServiceDownloadBatchTemplate**](#deviceservicedownloadbatchtemplate) | **GET** /v1/devices/batch/template | |
|[**deviceServiceGetDevice**](#deviceservicegetdevice) | **GET** /v1/devices/{id} | |
|[**deviceServiceGetDeviceByKey**](#deviceservicegetdevicebykey) | **GET** /v1/devices/key/{deviceKey} | |
|[**deviceServiceGetDeviceStatistics**](#deviceservicegetdevicestatistics) | **GET** /v1/device-statistics | |
|[**deviceServiceGetDeviceTelemetry**](#deviceservicegetdevicetelemetry) | **GET** /v1/devices/{deviceKey}/telemetry | |
|[**deviceServiceGetMqttParameters**](#deviceservicegetmqttparameters) | **GET** /v1/devices/{deviceKey}/mqtt-parameters | |
|[**deviceServiceGetPushRecord**](#deviceservicegetpushrecord) | **GET** /v1/devices/push-records/{pushRecordId} | |
|[**deviceServiceListDevices**](#deviceservicelistdevices) | **GET** /v1/devices | |
|[**deviceServiceListPushRecords**](#deviceservicelistpushrecords) | **GET** /v1/devices/{deviceKey}/push-records | |
|[**deviceServiceMockKafka**](#deviceservicemockkafka) | **POST** /v1/devices/mock-kafka | |
|[**deviceServiceRestoreDevice**](#deviceservicerestoredevice) | **POST** /v1/devices/{id}/restore | |
|[**deviceServiceSetDeviceEnabled**](#deviceservicesetdeviceenabled) | **POST** /v1/devices/{id}/enabled | |
|[**deviceServiceSimulatePush**](#deviceservicesimulatepush) | **POST** /v1/devices/{deviceKey}/simulate-push | |
|[**deviceServiceUpdateDevice**](#deviceserviceupdatedevice) | **PUT** /v1/devices/{id} | |

# **deviceServiceActivateDevice**
> DeviceV1ActivateDeviceResponse deviceServiceActivateDevice(deviceV1ActivateDeviceRequest)

ActivateDevice activates a device

### Example

```typescript
import {
    DeviceServiceApi,
    Configuration,
    DeviceV1ActivateDeviceRequest
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new DeviceServiceApi(configuration);

let id: string; // (default to undefined)
let deviceV1ActivateDeviceRequest: DeviceV1ActivateDeviceRequest; //

const { status, data } = await apiInstance.deviceServiceActivateDevice(
    id,
    deviceV1ActivateDeviceRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **deviceV1ActivateDeviceRequest** | **DeviceV1ActivateDeviceRequest**|  | |
| **id** | [**string**] |  | defaults to undefined|


### Return type

**DeviceV1ActivateDeviceResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deviceServiceBatchUploadDevices**
> DeviceV1BatchUploadDevicesResponse deviceServiceBatchUploadDevices(deviceV1BatchUploadDevicesRequest)

BatchUploadDevices uploads an Excel file to create multiple devices in batch

### Example

```typescript
import {
    DeviceServiceApi,
    Configuration,
    DeviceV1BatchUploadDevicesRequest
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new DeviceServiceApi(configuration);

let deviceV1BatchUploadDevicesRequest: DeviceV1BatchUploadDevicesRequest; //

const { status, data } = await apiInstance.deviceServiceBatchUploadDevices(
    deviceV1BatchUploadDevicesRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **deviceV1BatchUploadDevicesRequest** | **DeviceV1BatchUploadDevicesRequest**|  | |


### Return type

**DeviceV1BatchUploadDevicesResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deviceServiceClearPushRecords**
> DeviceV1ClearPushRecordsResponse deviceServiceClearPushRecords()

ClearPushRecords clears push records for a device

### Example

```typescript
import {
    DeviceServiceApi,
    Configuration
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new DeviceServiceApi(configuration);

let deviceKey: string; // (default to undefined)
let beforeTimestamp: string; //Optional: clear records before this timestamp (milliseconds)  If not set, clear all records older than 7 days (optional) (default to undefined)

const { status, data } = await apiInstance.deviceServiceClearPushRecords(
    deviceKey,
    beforeTimestamp
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **deviceKey** | [**string**] |  | defaults to undefined|
| **beforeTimestamp** | [**string**] | Optional: clear records before this timestamp (milliseconds)  If not set, clear all records older than 7 days | (optional) defaults to undefined|


### Return type

**DeviceV1ClearPushRecordsResponse**

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

# **deviceServiceCreateDevice**
> DeviceV1CreateDeviceResponse deviceServiceCreateDevice(deviceV1CreateDeviceRequest)

CreateDevice creates a new device

### Example

```typescript
import {
    DeviceServiceApi,
    Configuration,
    DeviceV1CreateDeviceRequest
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new DeviceServiceApi(configuration);

let deviceV1CreateDeviceRequest: DeviceV1CreateDeviceRequest; //

const { status, data } = await apiInstance.deviceServiceCreateDevice(
    deviceV1CreateDeviceRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **deviceV1CreateDeviceRequest** | **DeviceV1CreateDeviceRequest**|  | |


### Return type

**DeviceV1CreateDeviceResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deviceServiceDeleteDevice**
> DeviceV1DeleteDeviceResponse deviceServiceDeleteDevice()

DeleteDevice deletes a device

### Example

```typescript
import {
    DeviceServiceApi,
    Configuration
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new DeviceServiceApi(configuration);

let id: string; // (default to undefined)

const { status, data } = await apiInstance.deviceServiceDeleteDevice(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] |  | defaults to undefined|


### Return type

**DeviceV1DeleteDeviceResponse**

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

# **deviceServiceDownloadBatchTemplate**
> deviceServiceDownloadBatchTemplate()

DownloadBatchTemplate downloads an Excel template for batch device upload

### Example

```typescript
import {
    DeviceServiceApi,
    Configuration
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new DeviceServiceApi(configuration);

const { status, data } = await apiInstance.deviceServiceDownloadBatchTemplate();
```

### Parameters
This endpoint does not have any parameters.


### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deviceServiceGetDevice**
> DeviceV1GetDeviceResponse deviceServiceGetDevice()

GetDevice gets a device by ID or DeviceName

### Example

```typescript
import {
    DeviceServiceApi,
    Configuration
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new DeviceServiceApi(configuration);

let id: string; // (default to undefined)
let deviceKey: string; // (optional) (default to undefined)

const { status, data } = await apiInstance.deviceServiceGetDevice(
    id,
    deviceKey
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] |  | defaults to undefined|
| **deviceKey** | [**string**] |  | (optional) defaults to undefined|


### Return type

**DeviceV1GetDeviceResponse**

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

# **deviceServiceGetDeviceByKey**
> DeviceV1GetDeviceResponse deviceServiceGetDeviceByKey()

GetDeviceByKey gets a device by DeviceKey

### Example

```typescript
import {
    DeviceServiceApi,
    Configuration
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new DeviceServiceApi(configuration);

let deviceKey: string; // (default to undefined)

const { status, data } = await apiInstance.deviceServiceGetDeviceByKey(
    deviceKey
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **deviceKey** | [**string**] |  | defaults to undefined|


### Return type

**DeviceV1GetDeviceResponse**

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

# **deviceServiceGetDeviceStatistics**
> DeviceV1GetDeviceStatisticsResponse deviceServiceGetDeviceStatistics()

GetDeviceStatistics gets device statistics

### Example

```typescript
import {
    DeviceServiceApi,
    Configuration
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new DeviceServiceApi(configuration);

const { status, data } = await apiInstance.deviceServiceGetDeviceStatistics();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**DeviceV1GetDeviceStatisticsResponse**

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

# **deviceServiceGetDeviceTelemetry**
> DeviceV1GetDeviceTelemetryResponse deviceServiceGetDeviceTelemetry()

GetDeviceTelemetry gets the latest telemetry data for a device

### Example

```typescript
import {
    DeviceServiceApi,
    Configuration
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new DeviceServiceApi(configuration);

let deviceKey: string; // (default to undefined)

const { status, data } = await apiInstance.deviceServiceGetDeviceTelemetry(
    deviceKey
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **deviceKey** | [**string**] |  | defaults to undefined|


### Return type

**DeviceV1GetDeviceTelemetryResponse**

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

# **deviceServiceGetMqttParameters**
> DeviceV1GetMqttParametersResponse deviceServiceGetMqttParameters()

GetMqttParameters gets MQTT connection parameters for a device

### Example

```typescript
import {
    DeviceServiceApi,
    Configuration
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new DeviceServiceApi(configuration);

let deviceKey: string; // (default to undefined)

const { status, data } = await apiInstance.deviceServiceGetMqttParameters(
    deviceKey
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **deviceKey** | [**string**] |  | defaults to undefined|


### Return type

**DeviceV1GetMqttParametersResponse**

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

# **deviceServiceGetPushRecord**
> DeviceV1GetPushRecordResponse deviceServiceGetPushRecord()

GetPushRecord gets a specific push record

### Example

```typescript
import {
    DeviceServiceApi,
    Configuration
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new DeviceServiceApi(configuration);

let pushRecordId: string; // (default to undefined)

const { status, data } = await apiInstance.deviceServiceGetPushRecord(
    pushRecordId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **pushRecordId** | [**string**] |  | defaults to undefined|


### Return type

**DeviceV1GetPushRecordResponse**

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

# **deviceServiceListDevices**
> DeviceV1ListDevicesResponse deviceServiceListDevices()

ListDevices lists devices with pagination and filtering

### Example

```typescript
import {
    DeviceServiceApi,
    Configuration
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new DeviceServiceApi(configuration);

let page: number; // (optional) (default to undefined)
let pageSize: number; // (optional) (default to undefined)
let productId: string; // (optional) (default to undefined)
let states: Array<string>; // (optional) (default to undefined)
let enabled: boolean; // (optional) (default to undefined)
let searchText: string; // (optional) (default to undefined)

const { status, data } = await apiInstance.deviceServiceListDevices(
    page,
    pageSize,
    productId,
    states,
    enabled,
    searchText
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] |  | (optional) defaults to undefined|
| **pageSize** | [**number**] |  | (optional) defaults to undefined|
| **productId** | [**string**] |  | (optional) defaults to undefined|
| **states** | **Array&lt;string&gt;** |  | (optional) defaults to undefined|
| **enabled** | [**boolean**] |  | (optional) defaults to undefined|
| **searchText** | [**string**] |  | (optional) defaults to undefined|


### Return type

**DeviceV1ListDevicesResponse**

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

# **deviceServiceListPushRecords**
> DeviceV1ListPushRecordsResponse deviceServiceListPushRecords()

ListPushRecords lists push records for a device

### Example

```typescript
import {
    DeviceServiceApi,
    Configuration
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new DeviceServiceApi(configuration);

let deviceKey: string; // (default to undefined)
let page: number; // (optional) (default to undefined)
let pageSize: number; // (optional) (default to undefined)
let operationType: string; // (optional) (default to undefined)
let status: string; // (optional) (default to undefined)

const { status, data } = await apiInstance.deviceServiceListPushRecords(
    deviceKey,
    page,
    pageSize,
    operationType,
    status
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **deviceKey** | [**string**] |  | defaults to undefined|
| **page** | [**number**] |  | (optional) defaults to undefined|
| **pageSize** | [**number**] |  | (optional) defaults to undefined|
| **operationType** | [**string**] |  | (optional) defaults to undefined|
| **status** | [**string**] |  | (optional) defaults to undefined|


### Return type

**DeviceV1ListPushRecordsResponse**

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

# **deviceServiceMockKafka**
> DeviceV1MockKafkaResponse deviceServiceMockKafka(deviceV1MockKafkaRequest)

MockKafka sends a mock message to Kafka topic for testing

### Example

```typescript
import {
    DeviceServiceApi,
    Configuration,
    DeviceV1MockKafkaRequest
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new DeviceServiceApi(configuration);

let deviceV1MockKafkaRequest: DeviceV1MockKafkaRequest; //

const { status, data } = await apiInstance.deviceServiceMockKafka(
    deviceV1MockKafkaRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **deviceV1MockKafkaRequest** | **DeviceV1MockKafkaRequest**|  | |


### Return type

**DeviceV1MockKafkaResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deviceServiceRestoreDevice**
> DeviceV1RestoreDeviceResponse deviceServiceRestoreDevice(deviceV1RestoreDeviceRequest)

RestoreDevice restores a soft-deleted device

### Example

```typescript
import {
    DeviceServiceApi,
    Configuration,
    DeviceV1RestoreDeviceRequest
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new DeviceServiceApi(configuration);

let id: string; // (default to undefined)
let deviceV1RestoreDeviceRequest: DeviceV1RestoreDeviceRequest; //

const { status, data } = await apiInstance.deviceServiceRestoreDevice(
    id,
    deviceV1RestoreDeviceRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **deviceV1RestoreDeviceRequest** | **DeviceV1RestoreDeviceRequest**|  | |
| **id** | [**string**] |  | defaults to undefined|


### Return type

**DeviceV1RestoreDeviceResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deviceServiceSetDeviceEnabled**
> DeviceV1SetDeviceEnabledResponse deviceServiceSetDeviceEnabled(deviceV1SetDeviceEnabledRequest)

SetDeviceEnabled enables or disables a device

### Example

```typescript
import {
    DeviceServiceApi,
    Configuration,
    DeviceV1SetDeviceEnabledRequest
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new DeviceServiceApi(configuration);

let id: string; // (default to undefined)
let deviceV1SetDeviceEnabledRequest: DeviceV1SetDeviceEnabledRequest; //

const { status, data } = await apiInstance.deviceServiceSetDeviceEnabled(
    id,
    deviceV1SetDeviceEnabledRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **deviceV1SetDeviceEnabledRequest** | **DeviceV1SetDeviceEnabledRequest**|  | |
| **id** | [**string**] |  | defaults to undefined|


### Return type

**DeviceV1SetDeviceEnabledResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deviceServiceSimulatePush**
> DeviceV1SimulatePushResponse deviceServiceSimulatePush(deviceV1SimulatePushRequest)

SimulatePush simulates a device push for online debugging

### Example

```typescript
import {
    DeviceServiceApi,
    Configuration,
    DeviceV1SimulatePushRequest
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new DeviceServiceApi(configuration);

let deviceKey: string; // (default to undefined)
let deviceV1SimulatePushRequest: DeviceV1SimulatePushRequest; //

const { status, data } = await apiInstance.deviceServiceSimulatePush(
    deviceKey,
    deviceV1SimulatePushRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **deviceV1SimulatePushRequest** | **DeviceV1SimulatePushRequest**|  | |
| **deviceKey** | [**string**] |  | defaults to undefined|


### Return type

**DeviceV1SimulatePushResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deviceServiceUpdateDevice**
> DeviceV1UpdateDeviceResponse deviceServiceUpdateDevice(deviceV1UpdateDeviceRequest)

UpdateDevice updates an existing device

### Example

```typescript
import {
    DeviceServiceApi,
    Configuration,
    DeviceV1UpdateDeviceRequest
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new DeviceServiceApi(configuration);

let id: string; // (default to undefined)
let deviceV1UpdateDeviceRequest: DeviceV1UpdateDeviceRequest; //

const { status, data } = await apiInstance.deviceServiceUpdateDevice(
    id,
    deviceV1UpdateDeviceRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **deviceV1UpdateDeviceRequest** | **DeviceV1UpdateDeviceRequest**|  | |
| **id** | [**string**] |  | defaults to undefined|


### Return type

**DeviceV1UpdateDeviceResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

