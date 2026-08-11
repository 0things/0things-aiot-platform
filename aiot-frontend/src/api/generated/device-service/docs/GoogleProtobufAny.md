# GoogleProtobufAny

Contains an arbitrary serialized message along with a @type that describes the type of the serialized message.

## Properties

| Name     | Type       | Description                         | Notes                             |
| -------- | ---------- | ----------------------------------- | --------------------------------- |
| **type** | **string** | The type of the serialized message. | [optional] [default to undefined] |

## Example

```typescript
import { GoogleProtobufAny } from '@api/device-service'

const instance: GoogleProtobufAny = {
  type,
}
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
