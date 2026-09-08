package consumer

import (
	"context"
	"testing"

	"0things/pkg/event"
	"data-engine/internal/engine"
	"data-engine/internal/service"
	"data-engine/internal/storage"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func TestManager_Start(t *testing.T) {
	v := viper.New()
	logger := zap.NewNop()
	shadow := storage.NewShadowStore(v, logger)
	ruleProcessor := engine.NewProcessor(v, logger, nil, shadow)
	eventProcessor := service.NewEventProcessor(v, logger)
	otaStore := &reportStoreStub{}
	otaProcessor := service.NewOTAProcessor(otaStore, logger)

	pubSub := gochannel.NewGoChannel(gochannel.Config{}, watermill.NopLogger{})
	defer pubSub.Close()

	eventConsumer := event.NewConsumer(pubSub, logger)
	telemetryConsumer := NewTelemetryConsumer(ruleProcessor, logger)
	eventConsumerHandler := NewEventConsumer(eventProcessor, logger)
	otaProgressConsumer := NewOTAProgressConsumer(otaProcessor, logger)

	manager := NewManager(eventConsumer, telemetryConsumer, eventConsumerHandler, otaProgressConsumer, logger)

	err := manager.Start(context.Background())
	if err != nil {
		t.Fatalf("unexpected error starting manager: %v", err)
	}
}
