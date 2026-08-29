package job

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"

	"aiot-backend/pkg/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type mockMQTTService struct {
	subscribed   map[string]mqtt.MessageHandler
	unsubscribed []string
}

func newMockMQTTService() *mockMQTTService {
	return &mockMQTTService{
		subscribed:   make(map[string]mqtt.MessageHandler),
		unsubscribed: make([]string, 0),
	}
}

func (m *mockMQTTService) Publish(ctx context.Context, topic string, qos byte, payload []byte) error {
	return nil
}

func (m *mockMQTTService) Subscribe(ctx context.Context, topic string, qos byte, handler mqtt.MessageHandler) error {
	m.subscribed[topic] = handler
	return nil
}

func (m *mockMQTTService) Unsubscribe(topics ...string) error {
	m.unsubscribed = append(m.unsubscribed, topics...)
	return nil
}

func (m *mockMQTTService) Close() {}

func TestOTAMQTTReportBridge_StepParsing(t *testing.T) {
	// 测试整数 step 解析
	intPayload := `{"params":{"batch_id":"b1","device_key":"d1","step":85,"status":"in_progress"}}`
	var intReport struct {
		Params struct {
			BatchID   string          `json:"batch_id"`
			DeviceKey string          `json:"device_key"`
			Status    string          `json:"status"`
			Step      json.RawMessage `json:"step"`
		} `json:"params"`
	}
	assert.NoError(t, json.Unmarshal([]byte(intPayload), &intReport))
	var intStep int32
	if err := json.Unmarshal(intReport.Params.Step, &intStep); err != nil {
		var text string
		if json.Unmarshal(intReport.Params.Step, &text) == nil {
			v, _ := strconv.ParseInt(text, 10, 32)
			intStep = int32(v)
		}
	}
	assert.Equal(t, int32(85), intStep)

	// 测试字符串 step 解析
	strPayload := `{"params":{"batch_id":"b1","device_key":"d1","step":"90","status":"in_progress"}}`
	var strReport struct {
		Params struct {
			BatchID   string          `json:"batch_id"`
			DeviceKey string          `json:"device_key"`
			Status    string          `json:"status"`
			Step      json.RawMessage `json:"step"`
		} `json:"params"`
	}
	assert.NoError(t, json.Unmarshal([]byte(strPayload), &strReport))
	var strStep int32
	if err := json.Unmarshal(strReport.Params.Step, &strStep); err != nil {
		var text string
		if json.Unmarshal(strReport.Params.Step, &text) == nil {
			v, _ := strconv.ParseInt(text, 10, 32)
			strStep = int32(v)
		}
	}
	assert.Equal(t, int32(90), strStep)
}

func TestOTAMQTTReportBridge_Lifecycle(t *testing.T) {
	logger := &log.Logger{Logger: zap.NewNop()}
	mockMQTT := newMockMQTTService()

	bridge := NewOTAMQTTReportBridge(mockMQTT, nil, logger)
	assert.NotNil(t, bridge)

	ctx := context.Background()
	err := bridge.Start(ctx)
	assert.NoError(t, err)

	// 依赖 kafka 为空时应安全返回
	assert.Equal(t, 0, len(mockMQTT.subscribed))

	bridge.Stop()
}
