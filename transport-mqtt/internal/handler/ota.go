package handler

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// HandleOTAProgress delegates OTA progress reporting to the OTA service.
func (h *IngressHandler) HandleOTAProgress(_ mqtt.Client, msg mqtt.Message) {
	if h.otaService != nil {
		_ = h.otaService.HandleProgress(context.Background(), msg)
	}
}

// HandleOTAInform delegates OTA device inform reporting to the OTA service.
func (h *IngressHandler) HandleOTAInform(_ mqtt.Client, msg mqtt.Message) {
	if h.otaService != nil {
		_ = h.otaService.HandleInform(context.Background(), msg)
	}
}
