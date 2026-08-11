# OtaV1CreateOTAPackageRequest

CreateOTAPackage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**packageName** | **string** |  | [optional] [default to undefined]
**version** | **string** |  | [optional] [default to undefined]
**productId** | **string** |  | [optional] [default to undefined]
**packageType** | **string** |  | [optional] [default to undefined]
**uploadType** | **string** |  | [optional] [default to undefined]
**fileUrl** | **string** |  | [optional] [default to undefined]
**fileSize** | **string** |  | [optional] [default to undefined]
**checksum** | **string** |  | [optional] [default to undefined]
**description** | **string** |  | [optional] [default to undefined]
**releaseNotes** | **string** |  | [optional] [default to undefined]
**metadata** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { OtaV1CreateOTAPackageRequest } from '@api/device-service';

const instance: OtaV1CreateOTAPackageRequest = {
    packageName,
    version,
    productId,
    packageType,
    uploadType,
    fileUrl,
    fileSize,
    checksum,
    description,
    releaseNotes,
    metadata,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
