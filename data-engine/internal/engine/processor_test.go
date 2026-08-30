package engine

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"data-engine/internal/model"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func TestProcessor_ProcessMessage(t *testing.T) {
	logger := zap.NewNop()
	proc := NewProcessor(viper.New(), logger)

	// 1. 测试常规温度解析 (低于阈值)
	normalMsg := model.DeviceMessage{
		DeviceKey:   "sensor_test_01",
		Transport:   "mqtt",
		MessageType: "telemetry",
		Payload:     json.RawMessage(`{"temperature": 25.0, "humidity": 50}`),
		Timestamp:   time.Now(),
	}

	if err := proc.ProcessMessage(context.Background(), normalMsg); err != nil {
		t.Errorf("ProcessMessage failed on normal telemetry: %v", err)
	}

	// 2. 测试高温告警触发分支 (高于 70.0°C)
	alarmMsg := model.DeviceMessage{
		DeviceKey:   "sensor_test_02",
		Transport:   "mqtt",
		MessageType: "telemetry",
		Payload:     json.RawMessage(`{"temperature": 85.5}`),
		Timestamp:   time.Now(),
	}

	if err := proc.ProcessMessage(context.Background(), alarmMsg); err != nil {
		t.Errorf("ProcessMessage failed on high temperature: %v", err)
	}

	// 3. 测试 params 嵌套格式
	nestedMsg := model.DeviceMessage{
		DeviceKey:   "sensor_test_03",
		Transport:   "http",
		MessageType: "telemetry",
		Payload:     json.RawMessage(`{"params": {"temperature": 30.2}}`),
		Timestamp:   time.Now(),
	}

	if err := proc.ProcessMessage(context.Background(), nestedMsg); err != nil {
		t.Errorf("ProcessMessage failed on nested params: %v", err)
	}
}
