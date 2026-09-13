package handler

import (
	"context"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// HandleDeviceEvent delegates device-event processing to the device-event service.
func (h *IngressHandler) HandleDeviceEvent(_ mqtt.Client, msg mqtt.Message) {
	if h.deviceEventService != nil {
		_ = h.deviceEventService.Handle(context.Background(), msg)
	}
}
