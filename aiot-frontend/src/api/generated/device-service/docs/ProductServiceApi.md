# ProductServiceApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**productServiceCreateProduct**](#productservicecreateproduct) | **POST** /v1/products | |
|[**productServiceDeleteProduct**](#productservicedeleteproduct) | **DELETE** /v1/products/{id} | |
|[**productServiceGetProduct**](#productservicegetproduct) | **GET** /v1/products/{id} | |
|[**productServiceGetProductByKey**](#productservicegetproductbykey) | **GET** /v1/products/key/{productKey} | |
|[**productServiceListProducts**](#productservicelistproducts) | **GET** /v1/products | |
|[**productServiceRestoreProduct**](#productservicerestoreproduct) | **POST** /v1/products/{id}/restore | |
|[**productServiceUpdateProduct**](#productserviceupdateproduct) | **PUT** /v1/products/{id} | |

# **productServiceCreateProduct**
> ProductV1CreateProductResponse productServiceCreateProduct(productV1CreateProductRequest)

CreateProduct creates a new product

### Example

```typescript
import {
    ProductServiceApi,
    Configuration,
    ProductV1CreateProductRequest
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new ProductServiceApi(configuration);

let productV1CreateProductRequest: ProductV1CreateProductRequest; //

const { status, data } = await apiInstance.productServiceCreateProduct(
    productV1CreateProductRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **productV1CreateProductRequest** | **ProductV1CreateProductRequest**|  | |


### Return type

**ProductV1CreateProductResponse**

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

# **productServiceDeleteProduct**
> ProductV1DeleteProductResponse productServiceDeleteProduct()

DeleteProduct deletes a product

### Example

```typescript
import {
    ProductServiceApi,
    Configuration
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new ProductServiceApi(configuration);

let id: string; // (default to undefined)

const { status, data } = await apiInstance.productServiceDeleteProduct(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] |  | defaults to undefined|


### Return type

**ProductV1DeleteProductResponse**

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

# **productServiceGetProduct**
> ProductV1GetProductResponse productServiceGetProduct()

GetProduct gets a product by ID or ProductKey

### Example

```typescript
import {
    ProductServiceApi,
    Configuration
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new ProductServiceApi(configuration);

let id: string; // (default to undefined)
let productKey: string; // (optional) (default to undefined)

const { status, data } = await apiInstance.productServiceGetProduct(
    id,
    productKey
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] |  | defaults to undefined|
| **productKey** | [**string**] |  | (optional) defaults to undefined|


### Return type

**ProductV1GetProductResponse**

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

# **productServiceGetProductByKey**
> ProductV1GetProductByKeyResponse productServiceGetProductByKey()

GetProductByKey gets a product by product key

### Example

```typescript
import {
    ProductServiceApi,
    Configuration
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new ProductServiceApi(configuration);

let productKey: string; // (default to undefined)

const { status, data } = await apiInstance.productServiceGetProductByKey(
    productKey
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **productKey** | [**string**] |  | defaults to undefined|


### Return type

**ProductV1GetProductByKeyResponse**

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

# **productServiceListProducts**
> ProductV1ListProductsResponse productServiceListProducts()

ListProducts lists products with pagination and filtering

### Example

```typescript
import {
    ProductServiceApi,
    Configuration
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new ProductServiceApi(configuration);

let page: number; // (optional) (default to undefined)
let pageSize: number; // (optional) (default to undefined)
let category: string; // (optional) (default to undefined)
let status: string; // (optional) (default to undefined)
let searchText: string; // (optional) (default to undefined)

const { status, data } = await apiInstance.productServiceListProducts(
    page,
    pageSize,
    category,
    status,
    searchText
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] |  | (optional) defaults to undefined|
| **pageSize** | [**number**] |  | (optional) defaults to undefined|
| **category** | [**string**] |  | (optional) defaults to undefined|
| **status** | [**string**] |  | (optional) defaults to undefined|
| **searchText** | [**string**] |  | (optional) defaults to undefined|


### Return type

**ProductV1ListProductsResponse**

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

# **productServiceRestoreProduct**
> ProductV1RestoreProductResponse productServiceRestoreProduct(productV1RestoreProductRequest)

RestoreProduct restores a soft-deleted product

### Example

```typescript
import {
    ProductServiceApi,
    Configuration,
    ProductV1RestoreProductRequest
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new ProductServiceApi(configuration);

let id: string; // (default to undefined)
let productV1RestoreProductRequest: ProductV1RestoreProductRequest; //

const { status, data } = await apiInstance.productServiceRestoreProduct(
    id,
    productV1RestoreProductRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **productV1RestoreProductRequest** | **ProductV1RestoreProductRequest**|  | |
| **id** | [**string**] |  | defaults to undefined|


### Return type

**ProductV1RestoreProductResponse**

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

# **productServiceUpdateProduct**
> ProductV1UpdateProductResponse productServiceUpdateProduct(productV1UpdateProductRequest)

UpdateProduct updates an existing product

### Example

```typescript
import {
    ProductServiceApi,
    Configuration,
    ProductV1UpdateProductRequest
} from '@api/device-service';

const configuration = new Configuration();
const apiInstance = new ProductServiceApi(configuration);

let id: string; // (default to undefined)
let productV1UpdateProductRequest: ProductV1UpdateProductRequest; //

const { status, data } = await apiInstance.productServiceUpdateProduct(
    id,
    productV1UpdateProductRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **productV1UpdateProductRequest** | **ProductV1UpdateProductRequest**|  | |
| **id** | [**string**] |  | defaults to undefined|


### Return type

**ProductV1UpdateProductResponse**

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

