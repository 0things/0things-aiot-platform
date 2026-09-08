## Purpose

Enables protocol-aware, event-driven OTA firmware upgrade scheduling and device progress tracking across backend, transport, and data-engine services.

## ADDED Requirements

### Requirement: Protocol-Aware OTA Command Dispatch
The platform SHALL distinguish the target device's transport type (MQTT / HTTP / CoAP) and execute protocol-appropriate dispatch strategies when an upgrade batch is initiated.

#### Scenario: MQTT device active push dispatch
- **WHEN** an OTA upgrade command targets an online MQTT device
- **THEN** the system publishes an `OTAUpgradeCommand` event to `TopicOTAUpgradeCommand` with `transport: mqtt` metadata, and `mqtt-transport` pushes the upgrade payload to the device's MQTT topic.

#### Scenario: HTTP device pull-based state staging
- **WHEN** an OTA upgrade command targets an HTTP device
- **THEN** the system stages the target firmware metadata into the device's desired state/shadow, allowing the HTTP device to pull the firmware package on its next polling request without triggering an immediate push event.

### Requirement: Event-Driven OTA Progress Reporting
The transport services (MQTT/HTTP/CoAP) SHALL forward device OTA upgrade progress reports to `TopicOTAProgressReport` on the event bus.

#### Scenario: Device reports upgrade progress
- **WHEN** a device posts progress update payload over its native transport protocol
- **THEN** the respective transport service parses the payload and publishes an `OTAUpgradeReport` event to `TopicOTAProgressReport`.

### Requirement: Atomic OTA Progress Processing
The `data-engine` service SHALL subscribe to `TopicOTAProgressReport` and atomically persist device upgrade state and recalculate batch statistics.

#### Scenario: Device reaches terminal status
- **WHEN** an `OTAUpgradeReport` event with status `SUCCESS` or `FAILED` is processed by `data-engine`
- **THEN** the device status is persisted into the database and batch aggregate counters are updated.
