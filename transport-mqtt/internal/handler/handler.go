package handler

import "transport-mqtt/internal/service"

// IngressHandler handles incoming MQTT traffic from devices and converts it to
// internal events. Each message type is handled in its own file.
type IngressHandler struct {
	telemetryService   *service.TelemetryService
	deviceEventService *service.DeviceEventService
	otaProgressService *service.OTAProgressService
}

// NewIngressHandler creates an MQTT ingress handler.
func NewIngressHandler(telemetryService *service.TelemetryService, deviceEventService *service.DeviceEventService, otaProgressService *service.OTAProgressService) *IngressHandler {
	return &IngressHandler{
		telemetryService:   telemetryService,
		deviceEventService: deviceEventService,
		otaProgressService: otaProgressService,
	}
}
