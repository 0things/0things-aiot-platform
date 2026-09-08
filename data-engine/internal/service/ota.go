package service

import (
	"context"
	"fmt"

	"0things/pkg/event"
	"data-engine/internal/repository"

	"go.uber.org/zap"
)

// OTAService coordinates OTA upgrade reports, state transitions, and database persistence in data-engine.
type OTAService interface {
	HandleOTAReport(ctx context.Context, report event.OTAUpgradeReport) error
}

type otaService struct {
	otaRepo repository.OTARepository
	logger  *zap.Logger
}

// NewOTAService creates an instance of OTAService.
func NewOTAService(otaRepo repository.OTARepository, logger *zap.Logger) OTAService {
	return &otaService{
		otaRepo: otaRepo,
		logger:  logger,
	}
}

// HandleOTAReport processes device OTA progress, inform, and failure events to advance task and batch states.
func (s *otaService) HandleOTAReport(ctx context.Context, report event.OTAUpgradeReport) error {
	eventType := report.EventType
	if eventType == "" {
		eventType = event.OTAEventTypeProgress
		report.EventType = event.OTAEventTypeProgress
	}
	if eventType != event.OTAEventTypeProgress && eventType != event.OTAEventTypeInform {
		return fmt.Errorf("unsupported OTA event_type %q", eventType)
	}
	s.logger.Info("processing OTA device report",
		zap.String("device_key", report.DeviceKey),
		zap.String("batch_id", report.BatchID),
		zap.String("event_type", string(eventType)),
	)
	return s.otaRepo.RecordReport(ctx, report)
}
