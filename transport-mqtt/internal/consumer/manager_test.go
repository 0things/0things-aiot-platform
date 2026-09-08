package consumer

import (
	"context"
	"testing"

	"0things/pkg/event"
	"transport-mqtt/pkg/log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func TestManager_Start_NilConsumer(t *testing.T) {
	v := viper.New()
	v.Set("log.mode", "console")
	v.Set("log.log_level", "error")
	logger := log.NewLog(v)

	otaCommandConsumer := NewOTACommandConsumer(nil, logger)
	mgr := NewManager(nil, otaCommandConsumer, logger)

	if err := mgr.Start(context.Background()); err != nil {
		t.Fatalf("expected nil error when consumer is nil, got %v", err)
	}
}

func TestManager_Start_Success(t *testing.T) {
	v := viper.New()
	v.Set("log.mode", "console")
	v.Set("log.log_level", "error")
	logger := log.NewLog(v)

	pubSub := gochannel.NewGoChannel(gochannel.Config{}, watermill.NopLogger{})
	defer pubSub.Close()

	eventConsumer := event.NewConsumer(pubSub, zap.NewNop())
	mockClient := &mockMQTTClient{connected: true}
	otaCommandConsumer := NewOTACommandConsumer(mockClient, logger)

	mgr := NewManager(eventConsumer, otaCommandConsumer, logger)

	if err := mgr.Start(context.Background()); err != nil {
		t.Fatalf("unexpected error starting manager: %v", err)
	}
}
