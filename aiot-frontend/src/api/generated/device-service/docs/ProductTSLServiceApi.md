# ProductTSLServiceApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**productTSLServiceCreateProductTSL**](#producttslservicecreateproducttsl) | **POST** /v1/products/{productKey}/tsl | |
|[**productTSLServiceDeleteProductTSL**](#producttslservicedeleteproducttsl) | **DELETE** /v1/products/{productKey}/tsl | |
|[**productTSLServiceGetProductTSL**](#producttslservicegetproducttsl) | **GET** /v1/products/{productKey}/tsl | |
|[**productTSLServiceUpdateProductTSL**](#producttslserviceupdateproducttsl) | **PUT** /v1/products/{productKey}/tsl | |

# **productTSLServiceCreateProductTSL**
> ProductTslV1CreateProductTSLResponse productTSLServiceCreateProductTSL(productTslV1CreateProductTSLRequest)

CreateProductTSL creates or updates TSL for a product

### Example

```typescript
import {
    ProductTSLServiceApi,
    Configuration,
    ProductTslV1CreateProductTSLRequest
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new ProductTSLServiceApi(configuration);

let productKey: string; // (default to undefined)
let productTslV1CreateProductTSLRequest: ProductTslV1CreateProductTSLRequest; //

const { status, data } = await apiInstance.productTSLServiceCreateProductTSL(
    productKey,
    productTslV1CreateProductTSLRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **productTslV1CreateProductTSLRequest** | **ProductTslV1CreateProductTSLRequest**|  | |
| **productKey** | [**string**] |  | defaults to undefined|


### Return type

**ProductTslV1CreateProductTSLResponse**

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

# **productTSLServiceDeleteProductTSL**
> ProductTslV1DeleteProductTSLResponse productTSLServiceDeleteProductTSL()

DeleteProductTSL deletes TSL for a product

### Example

```typescript
import {
    ProductTSLServiceApi,
    Configuration
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new ProductTSLServiceApi(configuration);

let productKey: string; // (default to undefined)

const { status, data } = await apiInstance.productTSLServiceDeleteProductTSL(
    productKey
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **productKey** | [**string**] |  | defaults to undefined|


### Return type

**ProductTslV1DeleteProductTSLResponse**

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

# **productTSLServiceGetProductTSL**
> ProductTslV1GetProductTSLResponse productTSLServiceGetProductTSL()

GetProductTSL gets TSL for a product

### Example

```typescript
import {
    ProductTSLServiceApi,
    Configuration
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new ProductTSLServiceApi(configuration);

let productKey: string; // (default to undefined)

const { status, data } = await apiInstance.productTSLServiceGetProductTSL(
    productKey
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **productKey** | [**string**] |  | defaults to undefined|


### Return type

**ProductTslV1GetProductTSLResponse**

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

# **productTSLServiceUpdateProductTSL**
> ProductTslV1UpdateProductTSLResponse productTSLServiceUpdateProductTSL(productTslV1UpdateProductTSLRequest)

UpdateProductTSL updates TSL for a product

### Example

```typescript
import {
    ProductTSLServiceApi,
    Configuration,
    ProductTslV1UpdateProductTSLRequest
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new ProductTSLServiceApi(configuration);

let productKey: string; // (default to undefined)
let productTslV1UpdateProductTSLRequest: ProductTslV1UpdateProductTSLRequest; //

const { status, data } = await apiInstance.productTSLServiceUpdateProductTSL(
    productKey,
    productTslV1UpdateProductTSLRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **productTslV1UpdateProductTSLRequest** | **ProductTslV1UpdateProductTSLRequest**|  | |
| **productKey** | [**string**] |  | defaults to undefined|


### Return type

**ProductTslV1UpdateProductTSLResponse**

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

