package consumer

import (
	"context"
	"testing"

	"0things/pkg/event"
	"0things/pkg/tsdb"
	"data-engine/internal/handler"
	"data-engine/internal/service"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"go.uber.org/zap"
)

func TestManager_Start(t *testing.T) {
	logger := zap.NewNop()
	mockTSDB := tsdb.NewMockClient(logger)
	telemetryService := service.NewTelemetryService(logger, mockTSDB)
	eventService := service.NewEventService(logger, &mockDeviceEventRepository{})
	otaStore := &mockDeviceUpgradeStatusRepository{}
	otaService := service.NewOTAService(otaStore, logger)

	pubSub := gochannel.NewGoChannel(gochannel.Config{}, watermill.NopLogger{})
	defer pubSub.Close()

	eventConsumer := event.NewConsumer(pubSub, logger)
	telemetryConsumer := NewTelemetryConsumer(handler.NewTelemetryHandler(telemetryService, logger))
	eventConsumerHandler := NewEventConsumer(handler.NewEventHandler(eventService, logger))
	otaHandler := handler.NewOTAHandler(otaService, logger)
	otaProgressConsumer := NewOTAProgressConsumer(otaHandler)
	otaDeviceInfoConsumer := NewOTADeviceInfoConsumer(otaHandler)

	manager := NewManager(eventConsumer, telemetryConsumer, eventConsumerHandler, otaProgressConsumer, otaDeviceInfoConsumer, logger)

	err := manager.Start(context.Background())
	if err != nil {
		t.Fatalf("unexpected error starting manager: %v", err)
	}
}
