package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"data-engine/internal/event"
	"data-engine/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type mockDeviceEventRepository struct {
	createdEvents []*model.DeviceEvent
	createErr     error
}

func (m *mockDeviceEventRepository) Create(ctx context.Context, event *model.DeviceEvent) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.createdEvents = append(m.createdEvents, event)
	return nil
}

func TestEventService_HandleEvent_StructuredPayload(t *testing.T) {
	ctx := context.Background()
	repo := &mockDeviceEventRepository{}
	svc := NewEventService(zap.NewNop(), repo)

	msg := event.DeviceMessage{
		DeviceKey:   "dev_motor_01",
		ProductKey:  "pk_motor_pro",
		Transport:   event.TransportMQTT,
		MessageType: event.MessageTypeEvent,
		Timestamp:   1726315200125,
		Payload: json.RawMessage(`{
			"identifier": "overheat_alarm",
			"type": "alert",
			"timestamp": 1726315200123,
			"data": {
				"temperature": 88.6,
				"threshold": 80.0
			}
		}`),
	}

	err := svc.HandleEvent(ctx, msg)
	require.NoError(t, err)

	// Verify persistence
	require.Len(t, repo.createdEvents, 1)
	saved := repo.createdEvents[0]
	assert.NotEmpty(t, saved.UUID)
	assert.Equal(t, "dev_motor_01", saved.DeviceKey)
	assert.Equal(t, "pk_motor_pro", saved.ProductKey)
	assert.Equal(t, "overheat_alarm", saved.EventIdentifier)
	assert.Equal(t, "alert", saved.EventType)
	assert.Equal(t, int64(1726315200123), saved.EventAt)
	assert.Contains(t, saved.Data, "88.6")
}

func TestEventService_HandleEvent_InfoEventNoAlarm(t *testing.T) {
	ctx := context.Background()
	repo := &mockDeviceEventRepository{}
	svc := NewEventService(zap.NewNop(), repo)

	msg := event.DeviceMessage{
		DeviceKey:   "dev_door_01",
		ProductKey:  "pk_door_pro",
		Transport:   event.TransportMQTT,
		MessageType: event.MessageTypeEvent,
		Timestamp:   1726315200000,
		Payload: json.RawMessage(`{
			"identifier": "door_opened",
			"type": "info",
			"timestamp": 1726315200000,
			"data": {
				"user_id": "usr_123"
			}
		}`),
	}

	err := svc.HandleEvent(ctx, msg)
	require.NoError(t, err)

	require.Len(t, repo.createdEvents, 1)
	assert.Equal(t, "door_opened", repo.createdEvents[0].EventIdentifier)
	assert.Equal(t, "info", repo.createdEvents[0].EventType)
	assert.Equal(t, int64(1726315200000), repo.createdEvents[0].EventAt)
}

func TestEventService_HandleEvent_MissingTimestampInPayload(t *testing.T) {
	ctx := context.Background()
	repo := &mockDeviceEventRepository{}
	svc := NewEventService(zap.NewNop(), repo)

	// Missing timestamp in payload -> rejected by struct validator
	msg := event.DeviceMessage{
		DeviceKey:   "dev_sensor_02",
		ProductKey:  "pk_sensor_02",
		Transport:   event.TransportHTTP,
		MessageType: event.MessageTypeEvent,
		Timestamp:   1726315200555,
		Payload: json.RawMessage(`{
			"identifier": "tamper_switch",
			"type": "error",
			"data": {
				"status": "opened"
			}
		}`),
	}

	err := svc.HandleEvent(ctx, msg)
	assert.NoError(t, err) // Invalid payload logged and skipped gracefully
	assert.Empty(t, repo.createdEvents)
}

func TestEventService_HandleEvent_InvalidEnvelope(t *testing.T) {
	ctx := context.Background()
	repo := &mockDeviceEventRepository{}
	svc := NewEventService(zap.NewNop(), repo)

	// Missing DeviceKey, Transport, MessageType
	msg := event.DeviceMessage{
		Payload: json.RawMessage(`{"identifier":"test","type":"info","timestamp":1726315200000}`),
	}

	err := svc.HandleEvent(ctx, msg)
	assert.NoError(t, err) // Invalid envelope should be logged and skipped gracefully
	assert.Empty(t, repo.createdEvents)
}

func TestEventService_HandleEvent_InvalidPayload(t *testing.T) {
	ctx := context.Background()
	repo := &mockDeviceEventRepository{}
	svc := NewEventService(zap.NewNop(), repo)

	msg := event.DeviceMessage{
		DeviceKey:   "dev_broken",
		Transport:   event.TransportMQTT,
		MessageType: event.MessageTypeEvent,
		Timestamp:   1726315200000,
		Payload:     json.RawMessage(`{invalid-json`),
	}

	err := svc.HandleEvent(ctx, msg)
	assert.NoError(t, err) // Malformed payload should be logged and skipped gracefully
	assert.Empty(t, repo.createdEvents)
}

func TestEventService_HandleEvent_InvalidPayloadFields(t *testing.T) {
	ctx := context.Background()
	repo := &mockDeviceEventRepository{}
	svc := NewEventService(zap.NewNop(), repo)

	// Missing required event identifier and invalid type
	msg := event.DeviceMessage{
		DeviceKey:   "dev_invalid",
		Transport:   event.TransportMQTT,
		MessageType: event.MessageTypeEvent,
		Timestamp:   1726315200000,
		Payload:     json.RawMessage(`{"identifier":"","type":"unsupported_level","timestamp":1726315200000}`),
	}

	err := svc.HandleEvent(ctx, msg)
	assert.NoError(t, err) // Struct validation failure should be logged and skipped gracefully
	assert.Empty(t, repo.createdEvents)
}

func TestEventService_HandleEvent_MissingType(t *testing.T) {
	ctx := context.Background()
	repo := &mockDeviceEventRepository{}
	svc := NewEventService(zap.NewNop(), repo)

	// Missing required type field
	msg := event.DeviceMessage{
		DeviceKey:   "dev_missing_type",
		Transport:   event.TransportMQTT,
		MessageType: event.MessageTypeEvent,
		Timestamp:   1726315200000,
		Payload:     json.RawMessage(`{"identifier":"temp_alarm","timestamp":1726315200000}`),
	}

	err := svc.HandleEvent(ctx, msg)
	assert.NoError(t, err) // Struct validation failure should be logged and skipped gracefully
	assert.Empty(t, repo.createdEvents)
}

func TestEventService_HandleEvent_RepoError(t *testing.T) {
	ctx := context.Background()
	repo := &mockDeviceEventRepository{createErr: errors.New("db connection failure")}
	svc := NewEventService(zap.NewNop(), repo)

	msg := event.DeviceMessage{
		DeviceKey:   "dev_fail",
		Transport:   event.TransportMQTT,
		MessageType: event.MessageTypeEvent,
		Timestamp:   1726315200000,
		Payload:     json.RawMessage(`{"identifier":"test","type":"info","timestamp":1726315200000}`),
	}

	err := svc.HandleEvent(ctx, msg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "persist device event")
}
