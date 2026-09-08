package event

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Producer interface {
	Publish(ctx context.Context, topic Topic, payload any, opts ...PublishOption) error
	Close() error
}

type watermillProducer struct {
	pub message.Publisher
}

func NewProducer(pub message.Publisher) Producer {
	return &watermillProducer{pub: pub}
}

func (p *watermillProducer) Publish(ctx context.Context, topic Topic, payload any, opts ...PublishOption) error {
	var data []byte
	switch v := payload.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		b, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		data = b
	}

	msg := message.NewMessage(watermill.NewUUID(), data)
	msg.SetContext(ctx)

	for _, opt := range opts {
		opt(msg)
	}

	return p.pub.Publish(topic.String(), msg)
}

func (p *watermillProducer) Close() error {
	if p.pub != nil {
		return p.pub.Close()
	}
	return nil
}
