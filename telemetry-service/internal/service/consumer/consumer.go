package consumer

import "github.com/google/wire"

// ProviderSet is consumer providers.
var ProviderSet = wire.NewSet(
	NewTelemetryConsumer,
	NewEventConsumer,
	NewClientConnectStateConsumer,
)
