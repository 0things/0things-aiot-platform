## Purpose

Provides a centralized, strongly typed, pluggable event bus infrastructure powered by NATS JetStream and Watermill for cross-microservice asynchronous communication.

## ADDED Requirements

### Requirement: Centralized Event Topic Enumeration
The platform SHALL maintain all event topic names in a dedicated enum type in the shared package `0things/pkg/event` using specific domain-action-version naming to prevent hardcoded topic strings across microservices.

#### Scenario: Valid topic enumeration lookup
- **WHEN** a service references `TopicOTAUpgradeCommand` or `TopicOTAProgressReport`
- **THEN** it receives standard versioned string constants `"ota.upgrade.command.v1"` and `"ota.progress.report.v1"` respectively.

### Requirement: Generic Strongly Typed Event Publishing
The platform SHALL provide a unified `Producer` interface that automatically serializes payloads to JSON and attaches execution context and metadata headers.

#### Scenario: Successful message publishing
- **WHEN** a service calls `producer.Publish(ctx, topic, payload)` with a Go struct
- **THEN** the payload is serialized to JSON and delivered to the configured NATS JetStream stream without throwing an error.

### Requirement: Generic Strongly Typed Event Subscription
The platform SHALL provide a generic `Subscribe[T](ctx, consumer, topic, handler)` helper that automatically deserializes binary payloads into the target struct `*T` and handles Ack/Nack and panic recovery.

#### Scenario: Message successfully handled
- **WHEN** a valid message arrives on the subscribed topic
- **THEN** the consumer deserializes it into `*T`, invokes the handler function, and issues an `Ack` upon successful completion.

#### Scenario: Handler returns error or panics
- **WHEN** the handler returns a non-nil error or panics during message processing
- **THEN** the consumer logs the error/panic, catches the panic, and sends a `Nack` to request message redelivery according to retry policy.
