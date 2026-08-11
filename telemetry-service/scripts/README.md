# Quick Start - Testing Telemetry Service

## Prerequisites

1. Kafka running on `localhost:9092`
2. Service built and ready to run

## Step 1: Start Kafka (if not running)

```bash
# Using Docker
docker run -d --name kafka -p 9092:9092 \
  -e KAFKA_CFG_NODE_ID=0 \
  -e KAFKA_CFG_PROCESS_ROLES=controller,broker \
  -e KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093 \
  -e KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT \
  -e KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=0@localhost:9093 \
  -e KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER \
  bitnami/kafka:latest
```

## Step 2: Start the Telemetry Service

```bash
./bin/telemetry-service -conf ./configs
```

You should see:
```
INFO msg=kafka consumer connected to brokers: [localhost:9092], group: telemetry-service-group
INFO msg=Registered handler for 'telemetry' topic
INFO msg=Registered handler for 'events' topic
INFO msg=Starting telemetry Kafka consumer
INFO msg=[HTTP] server listening on: [::]:8013
```

## Step 3: Run Continuous Test Script

In another terminal:

```bash
./scripts/continuous-test.sh
```

You'll see output like:
```
╔════════════════════════════════════════════════════════════╗
║     Telemetry Service - Continuous Test Script            ║
╚════════════════════════════════════════════════════════════╝

[✓] Telemetry: page.view (user: user-a1b2c3d4..., session: session-e5f6g7h...)
[✓] Event:     user.login (user: user-i8j9k0l1...)
[✓] Telemetry: button.click (user: user-m2n3o4p5..., session: session-q6r7s8t...)
```

## Step 4: Check Service Logs

In the service terminal, you should see:
```
INFO msg=Published telemetry event: page.view
INFO msg=Received message from topic telemetry with key: user-a1b2c3d4
INFO msg=Processing telemetry event: page.view
INFO msg=Published event: user.login for user: user-i8j9k0l1
INFO msg=Received message from topic events with key: user-i8j9k0l1
INFO msg=Processing event: user.login for user: user-i8j9k0l1
```

## Manual Test (Single Request)

### Test Telemetry Endpoint:

```bash
curl -X POST http://localhost:8013/api/v1/telemetry \
  -H "Content-Type: application/json" \
  -d '{
    "event_type": "page.view",
    "user_id": "test-user",
    "session_id": "test-session",
    "data": {"page": "/dashboard"}
  }'
```

Expected response:
```json
{"code":200,"message":"Telemetry published successfully"}
```

### Test Events Endpoint:

```bash
curl -X POST http://localhost:8013/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "event_type": "user.login",
    "user_id": "test-user",
    "email": "test@example.com",
    "ip": "127.0.0.1"
  }'
```

Expected response:
```json
{"code":200,"message":"Event published successfully"}
```

## Customizing the Test Script

### Change Request Interval

```bash
# Send every 1 second
INTERVAL=1 ./scripts/continuous-test.sh

# Send every 5 seconds
INTERVAL=5 ./scripts/continuous-test.sh
```

### Use Different API URL

```bash
API_URL=http://prod-server:8013 ./scripts/continuous-test.sh
```

## Stop Testing

Press `Ctrl+C` in the test script terminal to stop and see statistics:

```
╔════════════════════════════════════════════════════════════╗
║                    Test Statistics                         ║
╚════════════════════════════════════════════════════════════╝
Total Requests:          50
Telemetry Success:       25
Telemetry Failed:        0
Events Success:          25
Events Failed:           0
```

## Full Documentation

See [docs/producer-api.md](../docs/producer-api.md) for complete API documentation.
