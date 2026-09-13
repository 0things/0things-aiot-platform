package handler

import (
	"context"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// HandleOTAProgress delegates OTA-report processing to the OTA service.
func (h *IngressHandler) HandleOTAProgress(_ mqtt.Client, msg mqtt.Message) {
	if h.otaProgressService != nil {
		_ = h.otaProgressService.Handle(context.Background(), msg)
	}
}
