package event

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type SampleEvent struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Count   int    `json:"count"`
}

func TestEventProducerAndConsumer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	logger := zap.NewNop()
	wmLogger := watermill.NopLogger{}
	pubSub := gochannel.NewGoChannel(gochannel.Config{
		OutputChannelBuffer: 100,
	}, wmLogger)
	defer pubSub.Close()

	producer := NewProducer(pubSub)
	consumer := NewConsumer(pubSub, logger)

	var wg sync.WaitGroup
	wg.Add(1)

	var received *SampleEvent
	var receivedMeta map[string]string

	err := Subscribe(ctx, consumer, TopicOTAUpgradeCommand, func(ctx context.Context, data *SampleEvent, meta map[string]string) error {
		received = data
		receivedMeta = meta
		wg.Done()
		return nil
	})
	if err != nil {
		t.Fatalf("failed to subscribe: %v", err)
	}

	payload := &SampleEvent{
		ID:      "evt-123",
		Message: "ota test",
		Count:   42,
	}

	err = producer.Publish(ctx, TopicOTAUpgradeCommand, payload, WithDeviceKey("dev-001"), WithTransport("mqtt"))
	if err != nil {
		t.Fatalf("failed to publish: %v", err)
	}

	wg.Wait()

	if received == nil {
		t.Fatal("expected to receive event, got nil")
	}
	if received.ID != "evt-123" || received.Message != "ota test" || received.Count != 42 {
		t.Errorf("unexpected event content: %+v", received)
	}
	if receivedMeta["device_key"] != "dev-001" || receivedMeta["transport"] != "mqtt" {
		t.Errorf("unexpected metadata: %+v", receivedMeta)
	}
}

func TestConsumerPanicRecovery(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	wmLogger := watermill.NopLogger{}
	pubSub := gochannel.NewGoChannel(gochannel.Config{
		OutputChannelBuffer: 10,
	}, wmLogger)
	defer pubSub.Close()

	producer := NewProducer(pubSub)
	consumer := NewConsumer(pubSub, zap.NewNop())

	panicHandled := make(chan struct{})

	err := Subscribe(ctx, consumer, TopicOTAProgressReport, func(ctx context.Context, data *SampleEvent, meta map[string]string) error {
		close(panicHandled)
		panic("simulated panic in handler")
	})
	if err != nil {
		t.Fatalf("failed to subscribe: %v", err)
	}

	err = producer.Publish(ctx, TopicOTAProgressReport, &SampleEvent{ID: "panic-test"})
	if err != nil {
		t.Fatalf("failed to publish: %v", err)
	}

	select {
	case <-panicHandled:
		// success: panic did not crash the process
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for panic handler to execute")
	}
}

func TestConsumerErrorHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	wmLogger := watermill.NopLogger{}
	pubSub := gochannel.NewGoChannel(gochannel.Config{
		OutputChannelBuffer: 10,
	}, wmLogger)
	defer pubSub.Close()

	producer := NewProducer(pubSub)
	consumer := NewConsumer(pubSub, zap.NewNop())

	errHandled := make(chan struct{})

	err := Subscribe(ctx, consumer, TopicDeviceTelemetryReport, func(ctx context.Context, data *SampleEvent, meta map[string]string) error {
		close(errHandled)
		return errors.New("business logic error")
	})
	if err != nil {
		t.Fatalf("failed to subscribe: %v", err)
	}

	err = producer.Publish(ctx, TopicDeviceTelemetryReport, &SampleEvent{ID: "err-test"})
	if err != nil {
		t.Fatalf("failed to publish: %v", err)
	}

	select {
	case <-errHandled:
		// success
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for error handler to execute")
	}
}

func TestNewBusDrivers(t *testing.T) {
	logger := zap.NewNop()

	// 1. Memory driver
	vMem := viper.New()
	vMem.Set("event.driver", "memory")
	pub, sub, err := NewBus(vMem, logger)
	if err != nil {
		t.Fatalf("failed to create memory bus: %v", err)
	}
	_ = pub.Close()
	_ = sub.Close()

	// 2. Unsupported driver
	vInvalid := viper.New()
	vInvalid.Set("event.driver", "unknown_driver")
	_, _, err = NewBus(vInvalid, logger)
	if err == nil {
		t.Fatal("expected error for unsupported driver, got nil")
	}
}

func TestSQLitePubSub(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger := zap.NewNop()

	// Use temporary SQLite file
	dbPath := t.TempDir() + "/test_event_bus.db"
	v := viper.New()
	v.Set("event.driver", "sqlite")
	v.Set("event.sqlite.path", dbPath)
	v.Set("event.sqlite.consumer_group", "test-group")

	pub, sub, err := NewBus(v, logger)
	if err != nil {
		t.Fatalf("failed to create SQLite bus: %v", err)
	}
	defer pub.Close()
	defer sub.Close()

	producer := NewProducer(pub)
	consumer := NewConsumer(sub, logger)

	var wg sync.WaitGroup
	wg.Add(1)

	var received *SampleEvent
	err = Subscribe(ctx, consumer, TopicOTAUpgradeCommand, func(ctx context.Context, data *SampleEvent, meta map[string]string) error {
		received = data
		wg.Done()
		return nil
	})
	if err != nil {
		t.Fatalf("failed to subscribe to SQLite bus: %v", err)
	}

	payload := &SampleEvent{ID: "sqlite-001", Message: "sqlite test", Count: 99}
	if err := producer.Publish(ctx, TopicOTAUpgradeCommand, payload); err != nil {
		t.Fatalf("failed to publish to SQLite bus: %v", err)
	}

	wg.Wait()

	if received == nil || received.ID != "sqlite-001" || received.Count != 99 {
		t.Fatalf("unexpected event from SQLite pub/sub: %+v", received)
	}
}

