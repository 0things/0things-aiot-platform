# OtaV1DeviceDeployment

DeviceDeployment - Individual device upgrade status

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**deviceId** | **string** |  | [optional] [default to undefined]
**deviceKey** | **string** |  | [optional] [default to undefined]
**deviceName** | **string** |  | [optional] [default to undefined]
**productId** | **string** |  | [optional] [default to undefined]
**productKey** | **string** |  | [optional] [default to undefined]
**currentVersion** | **string** |  | [optional] [default to undefined]
**upgradeBatchId** | **string** |  | [optional] [default to undefined]
**status** | **string** |  | [optional] [default to undefined]
**lastStatusChangeTime** | **string** |  | [optional] [default to undefined]
**createdAt** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { OtaV1DeviceDeployment } from '@api/device-service';

const instance: OtaV1DeviceDeployment = {
    deviceId,
    deviceKey,
    deviceName,
    productId,
    productKey,
    currentVersion,
    upgradeBatchId,
    status,
    lastStatusChangeTime,
    createdAt,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
