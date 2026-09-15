package service

import (
	"context"
	"testing"

	"0things/pkg/event"
	"transport-mqtt/pkg/log"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockMessage struct {
	topic   string
	payload []byte
}

func (m *mockMessage) Duplicate() bool   { return false }
func (m *mockMessage) Qos() byte         { return 0 }
func (m *mockMessage) Retained() bool    { return false }
func (m *mockMessage) Topic() string     { return m.topic }
func (m *mockMessage) MessageID() uint16 { return 0 }
func (m *mockMessage) Payload() []byte   { return m.payload }
func (m *mockMessage) Ack()              {}

type mockProducer struct {
	publishedTopic event.Topic
	publishedMsg   interface{}
}

func (p *mockProducer) Publish(_ context.Context, topic event.Topic, msg interface{}, _ ...event.PublishOption) error {
	p.publishedTopic = topic
	p.publishedMsg = msg
	return nil
}

func (p *mockProducer) Close() error { return nil }

func TestOTAService_HandleProgress(t *testing.T) {
	v := viper.New()
	v.Set("log.mode", "console")
	v.Set("log.log_level", "error")
	logger := log.NewLog(v)

	producer := &mockProducer{}
	svc := NewOTAService(producer, logger)

	// 1. Valid progress payload
	msg := &mockMessage{
		topic:   "/ota/device/progress/PK123/DEV001",
		payload: []byte(`{"batch_id":"b_1","step":60,"desc":"downloading"}`),
	}
	err := svc.HandleProgress(context.Background(), msg)
	require.NoError(t, err)
	assert.Equal(t, event.TopicOTAProgressReport, producer.publishedTopic)

	deviceMsg, ok := producer.publishedMsg.(*event.DeviceMessage)
	require.True(t, ok)
	assert.Equal(t, "DEV001", deviceMsg.DeviceKey)
	assert.Equal(t, "PK123", deviceMsg.ProductKey)
	assert.Equal(t, event.TransportMQTT, deviceMsg.Transport)
	assert.Equal(t, event.MessageTypeOTAProgress, deviceMsg.MessageType)
	assert.JSONEq(t, `{"batch_id":"b_1","step":60,"desc":"downloading"}`, string(deviceMsg.Payload))
	assert.Greater(t, deviceMsg.Timestamp, int64(0))
}

func TestOTAService_HandleInform(t *testing.T) {
	v := viper.New()
	v.Set("log.mode", "console")
	v.Set("log.log_level", "error")
	logger := log.NewLog(v)

	producer := &mockProducer{}
	svc := NewOTAService(producer, logger)

	msg := &mockMessage{
		topic:   "/ota/device/inform/PK123/DEV001",
		payload: []byte(`{"id":"uuid_1","version":"v2.0.0"}`),
	}
	err := svc.HandleInform(context.Background(), msg)
	require.NoError(t, err)
	assert.Equal(t, event.TopicOTADeviceInfo, producer.publishedTopic)

	deviceMsg, ok := producer.publishedMsg.(*event.DeviceMessage)
	require.True(t, ok)
	assert.Equal(t, "DEV001", deviceMsg.DeviceKey)
	assert.Equal(t, "PK123", deviceMsg.ProductKey)
	assert.Equal(t, event.TransportMQTT, deviceMsg.Transport)
	assert.Equal(t, event.MessageTypeOTAInform, deviceMsg.MessageType)
	assert.JSONEq(t, `{"id":"uuid_1","version":"v2.0.0"}`, string(deviceMsg.Payload))
	assert.Greater(t, deviceMsg.Timestamp, int64(0))
}
