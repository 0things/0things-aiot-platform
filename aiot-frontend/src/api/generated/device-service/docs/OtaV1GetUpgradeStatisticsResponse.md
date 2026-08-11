# OtaV1GetUpgradeStatisticsResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**packageId** | **string** |  | [optional] [default to undefined]
**totalTargetDevices** | **number** |  | [optional] [default to undefined]
**successfulUpgrades** | **number** |  | [optional] [default to undefined]
**failedUpgrades** | **number** |  | [optional] [default to undefined]
**cancelledUpgrades** | **number** |  | [optional] [default to undefined]
**pendingUpgrades** | **number** |  | [optional] [default to undefined]
**inProgressUpgrades** | **number** |  | [optional] [default to undefined]

## Example

```typescript
import { OtaV1GetUpgradeStatisticsResponse } from '@api/device-service';

const instance: OtaV1GetUpgradeStatisticsResponse = {
    packageId,
    totalTargetDevices,
    successfulUpgrades,
    failedUpgrades,
    cancelledUpgrades,
    pendingUpgrades,
    inProgressUpgrades,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
