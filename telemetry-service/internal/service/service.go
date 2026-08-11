package service

import "github.com/google/wire"

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewGreeterService,
	NewTelemetryService,          // HTTP API for publishing telemetry messages
	NewEventService,              // HTTP API for publishing device events
	NewClientConnectStateService, // HTTP API for publishing client connect state events
)
