package consumer

import (
	"context"
	"testing"
	"time"

	"0things/pkg/event"
	"data-engine/internal/service"

	"go.uber.org/zap"
)

type reportStoreStub struct {
	report event.OTAUpgradeReport
}

func (s *reportStoreStub) RecordReport(_ context.Context, report event.OTAUpgradeReport) error {
	s.report = report
	return nil
}

func TestOTAProgressConsumer_HandleProgressReport(t *testing.T) {
	logger := zap.NewNop()
	store := &reportStoreStub{}
	processor := service.NewOTAProcessor(store, logger)
	consumer := NewOTAProgressConsumer(processor, logger)

	progress := int32(50)
	report := &event.OTAUpgradeReport{
		EventType:  "progress",
		BatchID:    "b-01",
		DeviceKey:  "dev_ota_01",
		Stage:      "DOWNLOADING",
		Progress:   &progress,
		ReportedAt: time.Now().UTC(),
	}

	err := consumer.HandleProgressReport(context.Background(), report, nil)
	if err != nil {
		t.Fatalf("unexpected error in HandleProgressReport: %v", err)
	}

	if store.report.BatchID != "b-01" || store.report.DeviceKey != "dev_ota_01" || store.report.Stage != "DOWNLOADING" {
		t.Errorf("unexpected report stored: %+v", store.report)
	}
}
