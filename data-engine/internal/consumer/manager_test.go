package consumer

import (
	"context"
	"testing"

	"0things/pkg/event"
	"data-engine/internal/handler"
	"data-engine/internal/repository"
	"data-engine/internal/service"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func TestManager_Start(t *testing.T) {
	v := viper.New()
	logger := zap.NewNop()
	shadow := repository.NewShadowRepository(v, logger)
	telemetryService := service.NewTelemetryService(v, logger, nil, shadow)
	eventService := service.NewEventService(v, logger)
	otaStore := &reportStoreStub{}
	otaService := service.NewOTAService(otaStore, logger)

	pubSub := gochannel.NewGoChannel(gochannel.Config{}, watermill.NopLogger{})
	defer pubSub.Close()

	eventConsumer := event.NewConsumer(pubSub, logger)
	telemetryConsumer := NewTelemetryConsumer(handler.NewTelemetryHandler(telemetryService, logger))
	eventConsumerHandler := NewEventConsumer(handler.NewEventHandler(eventService, logger))
	otaProgressConsumer := NewOTAProgressConsumer(handler.NewOTAProgressHandler(otaService, logger))

	manager := NewManager(eventConsumer, telemetryConsumer, eventConsumerHandler, otaProgressConsumer, logger)

	err := manager.Start(context.Background())
	if err != nil {
		t.Fatalf("unexpected error starting manager: %v", err)
	}
}
