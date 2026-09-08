package event

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

type Handler[T any] func(ctx context.Context, data *T, meta map[string]string) error

type Consumer interface {
	Subscriber() message.Subscriber
	Close() error
}

type watermillConsumer struct {
	sub    message.Subscriber
	logger *zap.Logger
}

func NewConsumer(sub message.Subscriber, logger *zap.Logger) Consumer {
	return &watermillConsumer{
		sub:    sub,
		logger: logger,
	}
}

func (c *watermillConsumer) Subscriber() message.Subscriber {
	return c.sub
}

func (c *watermillConsumer) Close() error {
	if c.sub != nil {
		return c.sub.Close()
	}
	return nil
}

// Subscribe provides strongly-typed subscription with automatic JSON unmarshaling,
// Ack/Nack management, and panic recovery.
func Subscribe[T any](ctx context.Context, consumer Consumer, topic Topic, handler Handler[T]) error {
	if consumer == nil || consumer.Subscriber() == nil {
		return errors.New("event: consumer or subscriber is nil")
	}

	messages, err := consumer.Subscriber().Subscribe(ctx, topic.String())
	if err != nil {
		return fmt.Errorf("event: failed to subscribe to topic %s: %w", topic, err)
	}

	go func() {
		for msg := range messages {
			processMessage(msg, handler)
		}
	}()

	return nil
}

func processMessage[T any](msg *message.Message, handler Handler[T]) {
	defer func() {
		if r := recover(); r != nil {
			msg.Nack()
		}
	}()

	var target T
	if err := json.Unmarshal(msg.Payload, &target); err != nil {
		msg.Nack()
		return
	}

	if err := handler(msg.Context(), &target, msg.Metadata); err != nil {
		msg.Nack()
		return
	}

	msg.Ack()
}
