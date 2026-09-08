package service

import (
	"context"
	"testing"
	"time"

	"data-engine/internal/model"

	"go.uber.org/zap"
)

type reportStoreStub struct{ report model.OTAUpgradeReport }

func (s *reportStoreStub) RecordReport(_ context.Context, report model.OTAUpgradeReport) error {
	s.report = report
	return nil
}

func TestOTAProcessor_HandleOTAReport(t *testing.T) {
	logger := zap.NewNop()
	store := &reportStoreStub{}
	proc := NewOTAProcessor(store, logger)

	tests := []struct {
		name   string
		report model.OTAUpgradeReport
	}{
		{
			name: "in progress 50%",
			report: model.OTAUpgradeReport{
				BatchID:    "batch_001",
				DeviceKey:  "dev_001",
				EventType:  "progress",
				ReportedAt: time.Now(),
			},
		},
		{
			name: "completed 100%",
			report: model.OTAUpgradeReport{
				BatchID:         "batch_001",
				DeviceKey:       "dev_001",
				EventType:       "inform",
				ReportedVersion: "v2.0.0",
				ReportedAt:      time.Now(),
			},
		},
		{
			name: "failed report",
			report: model.OTAUpgradeReport{
				BatchID:    "batch_001",
				DeviceKey:  "dev_001",
				EventType:  "progress",
				Error:      &model.OTAReportError{Code: "checksum", Message: "checksum mismatch"},
				ReportedAt: time.Now(),
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
