package kafka

import "github.com/google/wire"

// ProviderSet is kafka providers.
var ProviderSet = wire.NewSet(NewProducer)
