package handler

import (
	"context"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// HandleTelemetry delegates telemetry processing to the telemetry service.
func (h *IngressHandler) HandleTelemetry(_ mqtt.Client, msg mqtt.Message) {
	if h.telemetryService != nil {
		_ = h.telemetryService.Handle(context.Background(), msg)
	}
}
