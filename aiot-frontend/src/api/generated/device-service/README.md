## @api/device-service@1.0.0

This generator creates TypeScript/JavaScript client that utilizes [axios](https://github.com/axios/axios). The generated Node module can be used in the following environments:

Environment
* Node.js
* Webpack
* Browserify

Language level
* ES5 - you must have a Promises/A+ library installed
* ES6

Module system
* CommonJS
* ES6 module system

It can be used in both TypeScript and JavaScript. In TypeScript, the definition will be automatically resolved via `package.json`. ([Reference](https://www.typescriptlang.org/docs/handbook/declaration-files/consumption.html))

### Building

To build and compile the typescript sources to javascript use:
```
npm install
npm run build
```

### Publishing

First build the package then run `npm publish`

### Consuming

navigate to the folder of your consuming project and run one of the following commands.

_published:_

```
npm install @api/device-service@1.0.0 --save
```

_unPublished (not recommended):_

```
npm install PATH_TO_GENERATED_PACKAGE --save
```

### Documentation for API Endpoints

All URIs are relative to *http://localhost*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*DeviceServiceApi* | [**deviceServiceActivateDevice**](docs/DeviceServiceApi.md#deviceserviceactivatedevice) | **POST** /v1/devices/{id}/activate | 
*DeviceServiceApi* | [**deviceServiceBatchUploadDevices**](docs/DeviceServiceApi.md#deviceservicebatchuploaddevices) | **POST** /v1/devices/batch/upload | 
*DeviceServiceApi* | [**deviceServiceClearPushRecords**](docs/DeviceServiceApi.md#deviceserviceclearpushrecords) | **DELETE** /v1/devices/{deviceKey}/push-records | 
*DeviceServiceApi* | [**deviceServiceCreateDevice**](docs/DeviceServiceApi.md#deviceservicecreatedevice) | **POST** /v1/devices | 
*DeviceServiceApi* | [**deviceServiceDeleteDevice**](docs/DeviceServiceApi.md#deviceservicedeletedevice) | **DELETE** /v1/devices/{id} | 
*DeviceServiceApi* | [**deviceServiceDownloadBatchTemplate**](docs/DeviceServiceApi.md#deviceservicedownloadbatchtemplate) | **GET** /v1/devices/batch/template | 
*DeviceServiceApi* | [**deviceServiceGetDevice**](docs/DeviceServiceApi.md#deviceservicegetdevice) | **GET** /v1/devices/{id} | 
*DeviceServiceApi* | [**deviceServiceGetDeviceByKey**](docs/DeviceServiceApi.md#deviceservicegetdevicebykey) | **GET** /v1/devices/key/{deviceKey} | 
*DeviceServiceApi* | [**deviceServiceGetDeviceStatistics**](docs/DeviceServiceApi.md#deviceservicegetdevicestatistics) | **GET** /v1/device-statistics | 
*DeviceServiceApi* | [**deviceServiceGetDeviceTelemetry**](docs/DeviceServiceApi.md#deviceservicegetdevicetelemetry) | **GET** /v1/devices/{deviceKey}/telemetry | 
*DeviceServiceApi* | [**deviceServiceGetMqttParameters**](docs/DeviceServiceApi.md#deviceservicegetmqttparameters) | **GET** /v1/devices/{deviceKey}/mqtt-parameters | 
*DeviceServiceApi* | [**deviceServiceGetPushRecord**](docs/DeviceServiceApi.md#deviceservicegetpushrecord) | **GET** /v1/devices/push-records/{pushRecordId} | 
*DeviceServiceApi* | [**deviceServiceListDevices**](docs/DeviceServiceApi.md#deviceservicelistdevices) | **GET** /v1/devices | 
*DeviceServiceApi* | [**deviceServiceListPushRecords**](docs/DeviceServiceApi.md#deviceservicelistpushrecords) | **GET** /v1/devices/{deviceKey}/push-records | 
*DeviceServiceApi* | [**deviceServiceMockKafka**](docs/DeviceServiceApi.md#deviceservicemockkafka) | **POST** /v1/devices/mock-kafka | 
*DeviceServiceApi* | [**deviceServiceRestoreDevice**](docs/DeviceServiceApi.md#deviceservicerestoredevice) | **POST** /v1/devices/{id}/restore | 
*DeviceServiceApi* | [**deviceServiceSetDeviceEnabled**](docs/DeviceServiceApi.md#deviceservicesetdeviceenabled) | **POST** /v1/devices/{id}/enabled | 
*DeviceServiceApi* | [**deviceServiceSimulatePush**](docs/DeviceServiceApi.md#deviceservicesimulatepush) | **POST** /v1/devices/{deviceKey}/simulate-push | 
*DeviceServiceApi* | [**deviceServiceUpdateDevice**](docs/DeviceServiceApi.md#deviceserviceupdatedevice) | **PUT** /v1/devices/{id} | 
*GreeterApi* | [**greeterSayHello**](docs/GreeterApi.md#greetersayhello) | **GET** /helloworld/{name} | 
*OTAPackageServiceApi* | [**oTAPackageServiceCreateOTAPackage**](docs/OTAPackageServiceApi.md#otapackageservicecreateotapackage) | **POST** /v1/ota-packages | 
*OTAPackageServiceApi* | [**oTAPackageServiceDeleteOTAPackage**](docs/OTAPackageServiceApi.md#otapackageservicedeleteotapackage) | **DELETE** /v1/ota-packages/{id} | 
*OTAPackageServiceApi* | [**oTAPackageServiceGetOTAPackage**](docs/OTAPackageServiceApi.md#otapackageservicegetotapackage) | **GET** /v1/ota-packages/{id} | 
*OTAPackageServiceApi* | [**oTAPackageServiceGetUpgradeStatistics**](docs/OTAPackageServiceApi.md#otapackageservicegetupgradestatistics) | **GET** /v1/ota-packages/{packageName}/upgrade-statistics | 
*OTAPackageServiceApi* | [**oTAPackageServiceListDeviceDeployments**](docs/OTAPackageServiceApi.md#otapackageservicelistdevicedeployments) | **GET** /v1/ota-packages/{packageName}/device-deployments | 
*OTAPackageServiceApi* | [**oTAPackageServiceListOTAPackages**](docs/OTAPackageServiceApi.md#otapackageservicelistotapackages) | **GET** /v1/ota-packages | 
*OTAPackageServiceApi* | [**oTAPackageServiceListUpgradeBatches**](docs/OTAPackageServiceApi.md#otapackageservicelistupgradebatches) | **GET** /v1/ota-packages/{packageName}/batches | 
*OTAPackageServiceApi* | [**oTAPackageServiceUpdateOTAPackage**](docs/OTAPackageServiceApi.md#otapackageserviceupdateotapackage) | **PUT** /v1/ota-packages/{id} | 
*ProductServiceApi* | [**productServiceCreateProduct**](docs/ProductServiceApi.md#productservicecreateproduct) | **POST** /v1/products | 
*ProductServiceApi* | [**productServiceDeleteProduct**](docs/ProductServiceApi.md#productservicedeleteproduct) | **DELETE** /v1/products/{id} | 
*ProductServiceApi* | [**productServiceGetProduct**](docs/ProductServiceApi.md#productservicegetproduct) | **GET** /v1/products/{id} | 
*ProductServiceApi* | [**productServiceGetProductByKey**](docs/ProductServiceApi.md#productservicegetproductbykey) | **GET** /v1/products/key/{productKey} | 
*ProductServiceApi* | [**productServiceListProducts**](docs/ProductServiceApi.md#productservicelistproducts) | **GET** /v1/products | 
*ProductServiceApi* | [**productServiceRestoreProduct**](docs/ProductServiceApi.md#productservicerestoreproduct) | **POST** /v1/products/{id}/restore | 
*ProductServiceApi* | [**productServiceUpdateProduct**](docs/ProductServiceApi.md#productserviceupdateproduct) | **PUT** /v1/products/{id} | 
*ProductTSLServiceApi* | [**productTSLServiceCreateProductTSL**](docs/ProductTSLServiceApi.md#producttslservicecreateproducttsl) | **POST** /v1/products/{productKey}/tsl | 
*ProductTSLServiceApi* | [**productTSLServiceDeleteProductTSL**](docs/ProductTSLServiceApi.md#producttslservicedeleteproducttsl) | **DELETE** /v1/products/{productKey}/tsl | 
*ProductTSLServiceApi* | [**productTSLServiceGetProductTSL**](docs/ProductTSLServiceApi.md#producttslservicegetproducttsl) | **GET** /v1/products/{productKey}/tsl | 
*ProductTSLServiceApi* | [**productTSLServiceUpdateProductTSL**](docs/ProductTSLServiceApi.md#producttslserviceupdateproducttsl) | **PUT** /v1/products/{productKey}/tsl | 


### Documentation For Models

 - [DeviceV1ActivateDeviceRequest](docs/DeviceV1ActivateDeviceRequest.md)
 - [DeviceV1ActivateDeviceResponse](docs/DeviceV1ActivateDeviceResponse.md)
 - [DeviceV1BatchUploadDevicesRequest](docs/DeviceV1BatchUploadDevicesRequest.md)
 - [DeviceV1BatchUploadDevicesResponse](docs/DeviceV1BatchUploadDevicesResponse.md)
 - [DeviceV1BatchUploadError](docs/DeviceV1BatchUploadError.md)
 - [DeviceV1ClearPushRecordsResponse](docs/DeviceV1ClearPushRecordsResponse.md)
 - [DeviceV1CreateDeviceRequest](docs/DeviceV1CreateDeviceRequest.md)
 - [DeviceV1CreateDeviceResponse](docs/DeviceV1CreateDeviceResponse.md)
 - [DeviceV1DeleteDeviceResponse](docs/DeviceV1DeleteDeviceResponse.md)
 - [DeviceV1Device](docs/DeviceV1Device.md)
 - [DeviceV1GetDeviceResponse](docs/DeviceV1GetDeviceResponse.md)
 - [DeviceV1GetDeviceStatisticsResponse](docs/DeviceV1GetDeviceStatisticsResponse.md)
 - [DeviceV1GetDeviceTelemetryResponse](docs/DeviceV1GetDeviceTelemetryResponse.md)
 - [DeviceV1GetMqttParametersResponse](docs/DeviceV1GetMqttParametersResponse.md)
 - [DeviceV1GetPushRecordResponse](docs/DeviceV1GetPushRecordResponse.md)
 - [DeviceV1ListDevicesResponse](docs/DeviceV1ListDevicesResponse.md)
 - [DeviceV1ListPushRecordsResponse](docs/DeviceV1ListPushRecordsResponse.md)
 - [DeviceV1MockKafkaRequest](docs/DeviceV1MockKafkaRequest.md)
 - [DeviceV1MockKafkaResponse](docs/DeviceV1MockKafkaResponse.md)
 - [DeviceV1PushRecord](docs/DeviceV1PushRecord.md)
 - [DeviceV1RestoreDeviceRequest](docs/DeviceV1RestoreDeviceRequest.md)
 - [DeviceV1RestoreDeviceResponse](docs/DeviceV1RestoreDeviceResponse.md)
 - [DeviceV1SetDeviceEnabledRequest](docs/DeviceV1SetDeviceEnabledRequest.md)
 - [DeviceV1SetDeviceEnabledResponse](docs/DeviceV1SetDeviceEnabledResponse.md)
 - [DeviceV1SimulatePushRequest](docs/DeviceV1SimulatePushRequest.md)
 - [DeviceV1SimulatePushResponse](docs/DeviceV1SimulatePushResponse.md)
 - [DeviceV1UpdateDeviceRequest](docs/DeviceV1UpdateDeviceRequest.md)
 - [DeviceV1UpdateDeviceResponse](docs/DeviceV1UpdateDeviceResponse.md)
 - [HelloworldV1HelloReply](docs/HelloworldV1HelloReply.md)
 - [OtaV1CreateOTAPackageRequest](docs/OtaV1CreateOTAPackageRequest.md)
 - [OtaV1CreateOTAPackageResponse](docs/OtaV1CreateOTAPackageResponse.md)
 - [OtaV1DeleteOTAPackageResponse](docs/OtaV1DeleteOTAPackageResponse.md)
 - [OtaV1DeviceDeployment](docs/OtaV1DeviceDeployment.md)
 - [OtaV1GetOTAPackageResponse](docs/OtaV1GetOTAPackageResponse.md)
 - [OtaV1GetUpgradeStatisticsResponse](docs/OtaV1GetUpgradeStatisticsResponse.md)
 - [OtaV1ListDeviceDeploymentsResponse](docs/OtaV1ListDeviceDeploymentsResponse.md)
 - [OtaV1ListOTAPackagesResponse](docs/OtaV1ListOTAPackagesResponse.md)
 - [OtaV1ListUpgradeBatchesResponse](docs/OtaV1ListUpgradeBatchesResponse.md)
 - [OtaV1OTAPackage](docs/OtaV1OTAPackage.md)
 - [OtaV1UpdateOTAPackageRequest](docs/OtaV1UpdateOTAPackageRequest.md)
 - [OtaV1UpdateOTAPackageResponse](docs/OtaV1UpdateOTAPackageResponse.md)
 - [OtaV1UpgradeBatch](docs/OtaV1UpgradeBatch.md)
 - [ProductTslV1CreateProductTSLRequest](docs/ProductTslV1CreateProductTSLRequest.md)
 - [ProductTslV1CreateProductTSLResponse](docs/ProductTslV1CreateProductTSLResponse.md)
 - [ProductTslV1DeleteProductTSLResponse](docs/ProductTslV1DeleteProductTSLResponse.md)
 - [ProductTslV1GetProductTSLResponse](docs/ProductTslV1GetProductTSLResponse.md)
 - [ProductTslV1ProductTSL](docs/ProductTslV1ProductTSL.md)
 - [ProductTslV1UpdateProductTSLRequest](docs/ProductTslV1UpdateProductTSLRequest.md)
 - [ProductTslV1UpdateProductTSLResponse](docs/ProductTslV1UpdateProductTSLResponse.md)
 - [ProductV1CreateProductRequest](docs/ProductV1CreateProductRequest.md)
 - [ProductV1CreateProductResponse](docs/ProductV1CreateProductResponse.md)
 - [ProductV1DeleteProductResponse](docs/ProductV1DeleteProductResponse.md)
 - [ProductV1GetProductByKeyResponse](docs/ProductV1GetProductByKeyResponse.md)
 - [ProductV1GetProductResponse](docs/ProductV1GetProductResponse.md)
 - [ProductV1ListProductsResponse](docs/ProductV1ListProductsResponse.md)
 - [ProductV1Product](docs/ProductV1Product.md)
 - [ProductV1RestoreProductRequest](docs/ProductV1RestoreProductRequest.md)
 - [ProductV1RestoreProductResponse](docs/ProductV1RestoreProductResponse.md)
 - [ProductV1UpdateProductRequest](docs/ProductV1UpdateProductRequest.md)
 - [ProductV1UpdateProductResponse](docs/ProductV1UpdateProductResponse.md)


<a id="documentation-for-authorization"></a>
## Documentation For Authorization

Endpoints do not require authorization.

