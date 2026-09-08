package event

import (
	"github.com/ThreeDotsLabs/watermill/message"
)

type PublishOption func(msg *message.Message)

func WithMetadata(key, value string) PublishOption {
	return func(msg *message.Message) {
		msg.Metadata.Set(key, value)
	}
}

func WithDeviceKey(deviceKey string) PublishOption {
	return func(msg *message.Message) {
		msg.Metadata.Set("device_key", deviceKey)
	}
}

func WithTransport(transport string) PublishOption {
	return func(msg *message.Message) {
		msg.Metadata.Set("transport", transport)
	}
}
