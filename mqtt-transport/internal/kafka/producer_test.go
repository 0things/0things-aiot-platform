package kafka

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"mqtt-transport/internal/model"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func TestNewProducer_MockMode(t *testing.T) {
	v := viper.New()
	logger := zap.NewNop()

	producer, cleanup, err := NewProducer(v, logger)
	if err != nil {
		t.Fatalf("unexpected error initializing producer without brokers: %v", err)
	}
	defer cleanup()

	msg := model.DeviceMessage{
		DeviceKey:   "dev_mock_01",
		Transport:   "mqtt",
		MessageType: "telemetry",
		Payload:     json.RawMessage(`{"temp": 25.5}`),
		Timestamp:   time.Now(),
	}

	// 在无 broker 模式下应平滑返回 nil
	err = producer.SendDeviceMessage(context.Background(), msg)
	if err != nil {
		t.Errorf("SendDeviceMessage in mock mode returned error: %v", err)
	}
}
