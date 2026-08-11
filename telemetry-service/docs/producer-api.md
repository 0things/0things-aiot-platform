# Telemetry Service - Producer API

## Overview

The Telemetry Service provides HTTP/gRPC APIs for publishing events to Kafka. The service automatically generates types and handlers from Protocol Buffer definitions.

## API Endpoints

### 1. Publish Telemetry Event

**Endpoint:** `POST /api/v1/telemetry`

Publishes telemetry events (page views, clicks, API calls, etc.) to the `telemetry` Kafka topic.

**Request Body:**
```json
{
  "type": "metric.cpu",
  "timestamp": 1704792000123,
  "product_key": "prod-abc123",
  "device_key": "dev-xyz456",
  "session_id": "session-789",
  "data": {
    "value": 45.5,
    "unit": "percent"
  }
}
```

**Response:**
```json
{
  "code": 200,
  "message": "Telemetry published successfully"
}
```

**Curl Example:**
```bash
curl -X POST http://localhost:8013/api/v1/telemetry \
  -H "Content-Type: application/json" \
  -d '{
    "type": "metric.cpu",
    "timestamp": 1704792000123,
    "product_key": "prod-abc123",
    "device_key": "dev-xyz456",
    "session_id": "session-789",
    "data": {
      "value": 75.5,
      "unit": "percent"
    }
  }'
```

### 2. Publish Device Event

**Endpoint:** `POST /api/v1/events`

Publishes device events (online, offline, alert, status updates, etc.) to the `events` Kafka topic.

**Request Body:**
```json
{
  "type": "device.online",
  "timestamp": 1704792000123,
  "product_key": "prod-abc123",
  "device_key": "dev-xyz456",
  "data": {
    "status": "online",
    "battery_level": 85,
    "firmware_version": "1.0.5"
  }
}
```

**Response:**
```json
{
  "code": 200,
  "message": "Event published successfully"
}
```

**Curl Example:**
```bash
curl -X POST http://localhost:8013/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "type": "device.online",
    "timestamp": 1704792000123,
    "product_key": "prod-abc123",
    "device_key": "dev-xyz456",
    "data": {
      "status": "online",
      "battery_level": 85
    }
  }'
```

## Testing

### Continuous Test Script

Use the provided script to continuously send test events to both endpoints:

```bash
# Run with default settings (localhost:8013, 2 second interval)
./scripts/continuous-test.sh

# Customize API URL and interval
API_URL=http://localhost:8013 INTERVAL=1 ./scripts/continuous-test.sh
```

The script will:
- Alternate between sending telemetry events and user events
- Generate random user IDs and session IDs
- Display success/failure statistics
- Show progress every 10 requests

**Output Example:**
```
╔════════════════════════════════════════════════════════════╗
║     Telemetry Service - Continuous Test Script            ║
╚════════════════════════════════════════════════════════════╝

API URL:              http://localhost:8013
Telemetry Endpoint:   http://localhost:8013/api/v1/telemetry
Events Endpoint:      http://localhost:8013/api/v1/events
Request Interval:     2s

Press Ctrl+C to stop

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

[✓] Telemetry: page.view (user: user-a1b2c3d4..., session: session-e5f6g7h...)
[✓] Event:     user.login (user: user-i8j9k0l1...)
```

Press `Ctrl+C` to stop the script and view statistics.

### Single Request Testing

#### Test Telemetry Endpoint:
```bash
./scripts/publish-test-event.sh
```

#### Manual Testing:
```bash
# Telemetry event
curl -X POST http://localhost:8013/api/v1/telemetry \
  -H "Content-Type: application/json" \
  -d '{
    "type": "metric.temperature",
    "timestamp": 1704792000123,
    "product_key": "prod-abc123",
    "device_key": "dev-xyz456",
    "session_id": "test-session",
    "data": {"value": 25.5, "unit": "celsius"}
  }'

# Device event
curl -X POST http://localhost:8013/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "type": "device.alert",
    "timestamp": 1704792000123,
    "product_key": "prod-abc123",
    "device_key": "dev-xyz456",
    "data": {"alert_type": "overheating", "value": 85.0}
  }'
```

## Architecture

### Request Flow

```
HTTP Client
    ↓
POST /api/v1/telemetry or /api/v1/events
    ↓
TelemetryService (generated from proto)
    ↓
Kafka Producer (internal/messaging)
    ↓
Kafka Topics: "telemetry" or "events"
    ↓
Kafka Consumers (TelemetryConsumer, EventConsumer)
    ↓
Business Logic (biz layer)
    ↓
Data Layer (data layer - database storage)
```

### Components

1. **API Definition** (`api/telemetry/v1/telemetry.proto`)
   - Protocol Buffer service and message definitions
   - Auto-generates HTTP handlers and Go types

2. **Service Layer** (`internal/service/telemetry_service.go`)
   - Implements the proto-generated interface
   - Validates requests
   - Publishes to Kafka

3. **Messaging Layer** (`internal/messaging/`)
   - Generic Kafka producer/consumer
   - Used by both HTTP API and background consumers

4. **Consumer Services** (`internal/service/*_consumer.go`)
   - Subscribe to Kafka topics
   - Process events asynchronously
   - Route to business logic

## Development

### Regenerate API Code

After modifying `api/telemetry/v1/telemetry.proto`:

```bash
make api
```

This regenerates:
- `api/telemetry/v1/telemetry.pb.go` - Protocol Buffer types
- `api/telemetry/v1/telemetry_http.pb.go` - HTTP handlers
- `api/telemetry/v1/telemetry_grpc.pb.go` - gRPC server interface

### Update Dependencies

```bash
cd cmd/telemetry-service
wire
cd ../..
go build -o ./bin/telemetry-service ./cmd/telemetry-service
```

## Event Types

### Telemetry Events (topic: `iotTelemetry`)

- `metric.cpu` - CPU usage metrics
- `metric.memory` - Memory usage metrics
- `metric.network` - Network traffic metrics
- `metric.temperature` - Temperature readings
- `metric.humidity` - Humidity readings
- `sensor.read` - Generic sensor data
- `data.upload` - Data upload events

### Device Events (topic: `iotEvent`)

- `device.online` - Device comes online
- `device.offline` - Device goes offline
- `device.alert` - Device alert/alarm
- `device.status` - Status update
- `device.heartbeat` - Heartbeat signal
- `device.upgrade` - Firmware upgrade

## Monitoring

Check consumer logs to verify events are being processed:

```bash
# Start service
./bin/telemetry-service -conf ./configs

# In another terminal, send test events
./scripts/continuous-test.sh

# Watch logs for:
# - "Published telemetry event: ..." (producer)
# - "Received message from topic telemetry..." (consumer)
# - "Processing telemetry event: ..." (business logic)
```

## Troubleshooting

### Kafka Not Configured

If you see:
```json
{"code": 503, "message": "Kafka producer not configured"}
```

Ensure Kafka is running and configured in `configs/config.yaml`:

```yaml
data:
  kafka:
    brokers:
      - localhost:9092
    client_id: telemetry-service
    producer:
      max_retries: 3
      timeout: 5s
    consumer:
      group_id: telemetry-service-group
```

### Validation Errors

**Missing type:**
```json
{"code": 400, "message": "type is required"}
```

**Missing product_key:**
```json
{"code": 400, "message": "product_key is required"}
```

**Missing device_key:**
```json
{"code": 400, "message": "device_key is required"}
```
