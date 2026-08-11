# Telemetry Kafka Consumer

This document describes the telemetry Kafka consumer implementation and how to use it.

## Architecture

The telemetry consumer follows the clean architecture pattern:

```
service/TelemetryConsumer → biz/TelemetryUsecase → data/TelemetryEventRepo
                ↓
        messaging/Consumer
```

## Components

### 1. TelemetryConsumer (service layer)
- Located in: `internal/service/telemetry_consumer.go`
- Subscribes to `iotTelemetry` topic
- Routes messages to the TelemetryUsecase
- Automatically starts when the application starts

### 2. TelemetryUsecase (business logic)
- Located in: `internal/biz/telemetry.go`
- Deserializes JSON event data
- Processes and validates events
- Persists events via the repository

### 3. TelemetryEventRepo (data layer)
- Located in: `internal/data/telemetry.go`
- Stores telemetry events
- Currently logs events (TODO: implement database storage)

## Event Format

Telemetry events are JSON messages with the following structure:

```json
{
  "event_type": "metric.cpu",
  "timestamp": "2024-01-09T12:00:00Z",
  "product_key": "prod-abc123",
  "device_key": "dev-xyz456",
  "session_id": "session-789",
  "data": {
    "value": 75.5,
    "unit": "percent",
    "custom_field": "value"
  }
}
```

### Fields

- `event_type` (string, required): Type of the event (e.g., "metric.cpu", "sensor.read")
- `timestamp` (string, required): ISO 8601 timestamp of when the event occurred
- `product_key` (string, optional): Product identifier
- `device_key` (string, optional): Device identifier
- `session_id` (string, optional): Session identifier
- `data` (object, optional): Additional event-specific data

## How It Works

1. **Startup**: When the application starts, `newApp()` in `main.go` launches the consumer in a goroutine
2. **Subscription**: Consumer subscribes to `iotTelemetry` topic during initialization
3. **Message Processing**:
   - Consumer polls Kafka for messages
   - Invokes handler for each message
   - Handler deserializes JSON and passes to TelemetryUsecase
   - Usecase saves event via repository
4. **Error Handling**: If processing fails, error is logged and consumer continues

## Publishing Events

Use the `TelemetryProducer` to publish events:

```go
import (
    "context"
    "telemetry-service/internal/biz"
    "telemetry-service/internal/service"
)

func publishEvent(producer *service.TelemetryProducer) error {
    event := &biz.TelemetryEvent{
        EventType:  "metric.cpu",
        Timestamp:  "2024-01-09T12:00:00Z",
        ProductKey: "prod-abc123",
        DeviceKey:  "dev-xyz456",
        SessionID:  "session-789",
        Data: map[string]interface{}{
            "value": 75.5,
            "unit":  "percent",
        },
    }

    return producer.PublishEvent(context.Background(), event)
}
```

## Testing the Consumer

### Using kafka-console-producer

```bash
# Start Kafka broker (if not running)
docker run -d --name kafka -p 9092:9092 \
  apache/kafka:latest

# Publish a test message
docker exec -it kafka /opt/kafka/bin/kafka-console-producer.sh \
  --bootstrap-server localhost:9092 \
  --topic iotTelemetry \
  --property "key.separator=:" \
  --property "parse.key=true"

# Then type (key:value format):
prod-abc123:{"event_type":"metric.cpu","timestamp":"2024-01-09T12:00:00Z","product_key":"prod-abc123","device_key":"dev-xyz456"}
```

### Using Go code

Create a simple test client:

```go
package main

import (
    "context"
    "encoding/json"

    "github.com/twmb/franz-go/pkg/kgo"
)

func main() {
    client, _ := kgo.NewClient(kgo.SeedBrokers("localhost:9092"))
    defer client.Close()

    event := map[string]interface{}{
        "event_type":  "metric.cpu",
        "timestamp":   "2024-01-09T12:00:00Z",
        "product_key": "prod-abc123",
        "device_key":  "dev-xyz456",
    }

    value, _ := json.Marshal(event)

    record := &kgo.Record{
        Topic: "iotTelemetry",
        Key:   []byte("prod-abc123"),
        Value: value,
    }

    client.ProduceSync(context.Background(), record)
}
```

## Monitoring

Watch the service logs to see events being processed:

```
INFO  Registered handler for 'iotTelemetry' topic
INFO  Starting telemetry Kafka consumer
INFO  Processing telemetry event: type=metric.cpu, product_key=prod-abc123, device_key=dev-xyz456, timestamp=2024-01-09T12:00:00Z
INFO  Saving telemetry event: type=metric.cpu, product_key=prod-abc123, device_key=dev-xyz456
```

## Configuration

Consumer configuration is in `configs/config.yaml`:

```yaml
data:
  kafka:
    brokers:
      - localhost:9092
    client_id: telemetry-service
    consumer:
      group_id: telemetry-service-group
      topics:
        - iotTelemetry
```

### Adding More Topics

To consume from additional topics:

1. Update `configs/config.yaml`:
   ```yaml
   consumer:
     topics:
       - iotTelemetry
       - iotEvent
   ```

2. Register handlers in `telemetry_consumer.go`:
   ```go
   msg.Consumer.Subscribe("iotTelemetry", messaging.MessageHandlerFunc(tc.handleTelemetryEvent))
   msg.Consumer.Subscribe("iotEvent", messaging.MessageHandlerFunc(tc.handleDeviceEvent))
   ```

## Next Steps

1. **Implement Database Storage**: Update `data/telemetry.go` to save events to PostgreSQL
2. **Add Tests**: Create unit tests for the consumer and use case
3. **Add Metrics**: Instrument with OpenTelemetry metrics
4. **Schema Validation**: Add JSON schema validation for events
5. **Dead Letter Queue**: Handle permanently failed messages

## Troubleshooting

### Consumer not receiving messages

1. Check Kafka broker is running: `docker ps | grep kafka`
2. Verify topic exists: `kafka-topics.sh --list --bootstrap-server localhost:9092`
3. Check consumer logs for connection errors
4. Verify config.yaml has correct broker addresses

### Messages being processed multiple times

- Check consumer group ID is unique per service instance
- Ensure handlers are idempotent
- Consider implementing offset management

### Consumer lag

- Monitor with `kafka-consumer-groups.sh --describe`
- Scale horizontally by adding more service instances
- Optimize handler processing time
