# Device

Device message

## Properties

| Name               | Type        | Description | Notes                             |
| ------------------ | ----------- | ----------- | --------------------------------- |
| **id**             | **string**  |             | [optional] [default to undefined] |
| **deviceKey**      | **string**  |             | [optional] [default to undefined] |
| **name**           | **string**  |             | [optional] [default to undefined] |
| **productId**      | **string**  |             | [optional] [default to undefined] |
| **productKey**     | **string**  |             | [optional] [default to undefined] |
| **productName**    | **string**  |             | [optional] [default to undefined] |
| **state**          | **string**  |             | [optional] [default to undefined] |
| **enabled**        | **boolean** |             | [optional] [default to undefined] |
| **lastOnlineTime** | **string**  |             | [optional] [default to undefined] |
| **metadata**       | **string**  |             | [optional] [default to undefined] |
| **createdAt**      | **string**  |             | [optional] [default to undefined] |
| **updatedAt**      | **string**  |             | [optional] [default to undefined] |
| **deletedAt**      | **string**  |             | [optional] [default to undefined] |

## Example

```typescript
import { Device } from '@api/device-service'

const instance: Device = {
  id,
  deviceKey,
  name,
  productId,
  productKey,
  productName,
  state,
  enabled,
  lastOnlineTime,
  metadata,
  createdAt,
  updatedAt,
  deletedAt,
}
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
