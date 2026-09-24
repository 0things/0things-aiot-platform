package adaptor

import (
	"encoding/json"

	"transport-mqtt/internal/event"
	"transport-mqtt/internal/dto"

	"github.com/go-playground/validator/v10"
)

// JsonMqttAdaptor normalizes device JSON payloads from MQTT into standard []event.DevicePropertyPostPayload.
type JsonMqttAdaptor struct {
	validate *validator.Validate
}

// NewJsonMqttAdaptor creates a new instance of JsonMqttAdaptor.
func NewJsonMqttAdaptor() *JsonMqttAdaptor {
	return &JsonMqttAdaptor{
		validate: validator.New(),
	}
}

// ConvertToTelemetryPayload converts ThingPropertyPostPayload into []event.DevicePropertyPostPayload.
func (a *JsonMqttAdaptor) ConvertToTelemetryPayload(raw []byte) ([]event.DevicePropertyPostPayload, error) {
	var req dto.ThingPropertyPostPayload
	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, err
	}

	if err := a.validate.Struct(&req); err != nil {
		return nil, err
	}

	// Group metrics by sampling timestamp
	grouped := make(map[int64]map[string]interface{})
	for k, item := range req.Params {
		ts := item.Time
		if grouped[ts] == nil {
			grouped[ts] = make(map[string]interface{})
		}
		grouped[ts][k] = item.Value
	}

	result := make([]event.DevicePropertyPostPayload, 0, len(grouped))
	for ts, vals := range grouped {
		result = append(result, event.DevicePropertyPostPayload{
			Timestamp: ts,
			Values:    vals,
		})
	}
	return result, nil
}
