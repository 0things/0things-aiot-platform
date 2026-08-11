# OTAPackageServiceApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**oTAPackageServiceCreateOTAPackage**](#otapackageservicecreateotapackage) | **POST** /v1/ota-packages | |
|[**oTAPackageServiceDeleteOTAPackage**](#otapackageservicedeleteotapackage) | **DELETE** /v1/ota-packages/{id} | |
|[**oTAPackageServiceGetOTAPackage**](#otapackageservicegetotapackage) | **GET** /v1/ota-packages/{id} | |
|[**oTAPackageServiceGetUpgradeStatistics**](#otapackageservicegetupgradestatistics) | **GET** /v1/ota-packages/{packageName}/upgrade-statistics | |
|[**oTAPackageServiceListDeviceDeployments**](#otapackageservicelistdevicedeployments) | **GET** /v1/ota-packages/{packageName}/device-deployments | |
|[**oTAPackageServiceListOTAPackages**](#otapackageservicelistotapackages) | **GET** /v1/ota-packages | |
|[**oTAPackageServiceListUpgradeBatches**](#otapackageservicelistupgradebatches) | **GET** /v1/ota-packages/{packageName}/batches | |
|[**oTAPackageServiceUpdateOTAPackage**](#otapackageserviceupdateotapackage) | **PUT** /v1/ota-packages/{id} | |

# **oTAPackageServiceCreateOTAPackage**
> OtaV1CreateOTAPackageResponse oTAPackageServiceCreateOTAPackage(otaV1CreateOTAPackageRequest)

CreateOTAPackage creates a new OTA package

### Example

```typescript
import {
    OTAPackageServiceApi,
    Configuration,
    OtaV1CreateOTAPackageRequest
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new OTAPackageServiceApi(configuration);

let otaV1CreateOTAPackageRequest: OtaV1CreateOTAPackageRequest; //

const { status, data } = await apiInstance.oTAPackageServiceCreateOTAPackage(
    otaV1CreateOTAPackageRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **otaV1CreateOTAPackageRequest** | **OtaV1CreateOTAPackageRequest**|  | |


### Return type

**OtaV1CreateOTAPackageResponse**

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

# **oTAPackageServiceDeleteOTAPackage**
> OtaV1DeleteOTAPackageResponse oTAPackageServiceDeleteOTAPackage()

DeleteOTAPackage deletes an OTA package

### Example

```typescript
import {
    OTAPackageServiceApi,
    Configuration
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new OTAPackageServiceApi(configuration);

let id: string; // (default to undefined)

const { status, data } = await apiInstance.oTAPackageServiceDeleteOTAPackage(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] |  | defaults to undefined|


### Return type

**OtaV1DeleteOTAPackageResponse**

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

# **oTAPackageServiceGetOTAPackage**
> OtaV1GetOTAPackageResponse oTAPackageServiceGetOTAPackage()

GetOTAPackage gets an OTA package by ID

### Example

```typescript
import {
    OTAPackageServiceApi,
    Configuration
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new OTAPackageServiceApi(configuration);

let id: string; // (default to undefined)

const { status, data } = await apiInstance.oTAPackageServiceGetOTAPackage(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] |  | defaults to undefined|


### Return type

**OtaV1GetOTAPackageResponse**

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

# **oTAPackageServiceGetUpgradeStatistics**
> OtaV1GetUpgradeStatisticsResponse oTAPackageServiceGetUpgradeStatistics()

GetUpgradeStatistics gets upgrade statistics for an OTA package

### Example

```typescript
import {
    OTAPackageServiceApi,
    Configuration
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new OTAPackageServiceApi(configuration);

let packageName: string; // (default to undefined)

const { status, data } = await apiInstance.oTAPackageServiceGetUpgradeStatistics(
    packageName
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **packageName** | [**string**] |  | defaults to undefined|


### Return type

**OtaV1GetUpgradeStatisticsResponse**

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

# **oTAPackageServiceListDeviceDeployments**
> OtaV1ListDeviceDeploymentsResponse oTAPackageServiceListDeviceDeployments()

ListDeviceDeployments lists device deployment status for an OTA package

### Example

```typescript
import {
    OTAPackageServiceApi,
    Configuration
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new OTAPackageServiceApi(configuration);

let packageName: string; // (default to undefined)
let page: number; // (optional) (default to undefined)
let pageSize: number; // (optional) (default to undefined)
let status: string; // (optional) (default to undefined)

const { status, data } = await apiInstance.oTAPackageServiceListDeviceDeployments(
    packageName,
    page,
    pageSize,
    status
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **packageName** | [**string**] |  | defaults to undefined|
| **page** | [**number**] |  | (optional) defaults to undefined|
| **pageSize** | [**number**] |  | (optional) defaults to undefined|
| **status** | [**string**] |  | (optional) defaults to undefined|


### Return type

**OtaV1ListDeviceDeploymentsResponse**

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

# **oTAPackageServiceListOTAPackages**
> OtaV1ListOTAPackagesResponse oTAPackageServiceListOTAPackages()

ListOTAPackages lists OTA packages with pagination and filtering

### Example

```typescript
import {
    OTAPackageServiceApi,
    Configuration
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new OTAPackageServiceApi(configuration);

let page: number; // (optional) (default to undefined)
let pageSize: number; // (optional) (default to undefined)
let productId: string; // (optional) (default to undefined)
let status: string; // (optional) (default to undefined)
let packageType: string; // (optional) (default to undefined)
let uploadType: string; // (optional) (default to undefined)

const { status, data } = await apiInstance.oTAPackageServiceListOTAPackages(
    page,
    pageSize,
    productId,
    status,
    packageType,
    uploadType
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] |  | (optional) defaults to undefined|
| **pageSize** | [**number**] |  | (optional) defaults to undefined|
| **productId** | [**string**] |  | (optional) defaults to undefined|
| **status** | [**string**] |  | (optional) defaults to undefined|
| **packageType** | [**string**] |  | (optional) defaults to undefined|
| **uploadType** | [**string**] |  | (optional) defaults to undefined|


### Return type

**OtaV1ListOTAPackagesResponse**

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

# **oTAPackageServiceListUpgradeBatches**
> OtaV1ListUpgradeBatchesResponse oTAPackageServiceListUpgradeBatches()

ListUpgradeBatches lists upgrade batches for an OTA package

### Example

```typescript
import {
    OTAPackageServiceApi,
    Configuration
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new OTAPackageServiceApi(configuration);

let packageName: string; // (default to undefined)

const { status, data } = await apiInstance.oTAPackageServiceListUpgradeBatches(
    packageName
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **packageName** | [**string**] |  | defaults to undefined|


### Return type

**OtaV1ListUpgradeBatchesResponse**

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

# **oTAPackageServiceUpdateOTAPackage**
> OtaV1UpdateOTAPackageResponse oTAPackageServiceUpdateOTAPackage(otaV1UpdateOTAPackageRequest)

UpdateOTAPackage updates an existing OTA package

### Example

```typescript
import {
    OTAPackageServiceApi,
    Configuration,
    OtaV1UpdateOTAPackageRequest
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new OTAPackageServiceApi(configuration);

let id: string; // (default to undefined)
let otaV1UpdateOTAPackageRequest: OtaV1UpdateOTAPackageRequest; //

const { status, data } = await apiInstance.oTAPackageServiceUpdateOTAPackage(
    id,
    otaV1UpdateOTAPackageRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **otaV1UpdateOTAPackageRequest** | **OtaV1UpdateOTAPackageRequest**|  | |
| **id** | [**string**] |  | defaults to undefined|


### Return type

**OtaV1UpdateOTAPackageResponse**

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

