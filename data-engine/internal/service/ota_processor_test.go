package service

import (
	"context"
	"testing"
	"time"

	"data-engine/internal/model"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func TestOTAProcessor_HandleOTAReport(t *testing.T) {
	logger := zap.NewNop()
	proc := NewOTAProcessor(viper.New(), logger)

	tests := []struct {
		name   string
		report model.OTADeviceUpgradeReportEvent
	}{
		{
			name: "in progress 50%",
			report: model.OTADeviceUpgradeReportEvent{
				BatchID:   "batch_001",
				DeviceKey: "dev_001",
				Progress:  50,
				Status:    "UPGRADING",
				Timestamp: time.Now(),
			},
		},
		{
			name: "completed 100%",
			report: model.OTADeviceUpgradeReportEvent{
				BatchID:   "batch_001",
				DeviceKey: "dev_001",
				Progress:  100,
				Version:   "v2.0.0",
				Timestamp: time.Now(),
			},
		},
		{
			name: "failed report",
			report: model.OTADeviceUpgradeReportEvent{
				BatchID:   "batch_001",
				DeviceKey: "dev_001",
				Progress:  -1,
				Status:    "FAILED",
				Desc:      "checksum mismatch",
				Timestamp: time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := proc.HandleOTAReport(context.Background(), tt.report); err != nil {
				t.Errorf("HandleOTAReport returned error: %v", err)
			}
		})
	}
}
