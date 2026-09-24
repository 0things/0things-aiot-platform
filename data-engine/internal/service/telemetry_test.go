package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"data-engine/internal/event"
	"data-engine/internal/tsdb"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func TestTelemetryService_ProcessPropertyPost(t *testing.T) {
	logger := zap.NewNop()
	v := viper.New()
	tsdbClient := tsdb.NewClient(v, logger)
	defer tsdbClient.Close()

	svc := NewTelemetryService(logger, tsdbClient)

	// 1. 标准归一化遥测数组载荷
	payloadData, err := json.Marshal([]event.DevicePropertyPostPayload{
		{
			Timestamp: 1524448722000,
			Values: map[string]interface{}{
				"Power": "on",
				"WF":    23.6,
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	msg := event.DeviceMessage{
		DeviceKey:   "sensor_test_01",
		Transport:   event.TransportMQTT,
		MessageType: event.MessageTypeTelemetry,
		Payload:     payloadData,
		Timestamp:   time.Now().UnixMilli(),
	}

	if err := svc.ProcessPropertyPost(context.Background(), msg); err != nil {
		t.Errorf("ProcessPropertyPost failed on normalized telemetry: %v", err)
	}

	// 2. 空 Payload 跳过测试
	emptyMsg := event.DeviceMessage{
		DeviceKey:   "sensor_test_02",
		Transport:   event.TransportMQTT,
		MessageType: event.MessageTypeTelemetry,
		Payload:     nil,
		Timestamp:   time.Now().UnixMilli(),
	}
	if err := svc.ProcessPropertyPost(context.Background(), emptyMsg); err != nil {
		t.Errorf("ProcessPropertyPost failed on empty payload: %v", err)
	}

	// 3. 非法 Payload 告警并安全忽略
	invalidMsg := event.DeviceMessage{
		DeviceKey:   "sensor_test_03",
		Transport:   event.TransportMQTT,
		MessageType: event.MessageTypeTelemetry,
		Payload:     json.RawMessage(`{not-valid-json`),
		Timestamp:   time.Now().UnixMilli(),
	}
	if err := svc.ProcessPropertyPost(context.Background(), invalidMsg); err != nil {
		t.Errorf("ProcessPropertyPost failed on invalid payload: %v", err)
	}
}
