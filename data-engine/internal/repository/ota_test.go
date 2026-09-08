package repository

import (
	"context"
	"testing"
	"time"

	"0things/pkg/event"
	"data-engine/internal/model"

	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func newTestOTARepository(t *testing.T) (OTARepository, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite in-memory db: %v", err)
	}

	if err := db.AutoMigrate(&model.Device{}, &model.DeviceUpgradeStatus{}, &model.UpgradeBatch{}); err != nil {
		t.Fatalf("failed to automigrate models: %v", err)
	}

	repo := NewRepository(zap.NewNop(), db)
	return NewOTARepository(repo), db
}

func TestOTARepository_RecordReport_Progress(t *testing.T) {
	ctx := context.Background()
	otaRepo, db := newTestOTARepository(t)

	// Seed records
	dev := &model.Device{ID: 1, DeviceKey: "dev_ota_01", Name: "d1"}
	if err := db.Create(dev).Error; err != nil {
		t.Fatalf("seed device failed: %v", err)
	}
	batch := &model.UpgradeBatch{ID: 1, BatchID: "batch-01", Status: string(event.OTAStatusInProgress)}
	if err := db.Create(batch).Error; err != nil {
		t.Fatalf("seed batch failed: %v", err)
	}
	task := &model.DeviceUpgradeStatus{
		ID:             1,
		DeviceID:       1,
		UpgradeBatchID: "batch-01",
		Status:         string(event.OTAStatusPending),
		TargetVersion:  "v2.0.0",
	}
	if err := db.Create(task).Error; err != nil {
		t.Fatalf("seed task failed: %v", err)
	}

	// 1. Report progress 50%
	progress := int32(50)
	rep := event.OTAUpgradeReport{
		EventType:  event.OTAEventTypeProgress,
		BatchID:    "batch-01",
		DeviceKey:  "dev_ota_01",
		Progress:   &progress,
		ReportedAt: time.Now().UTC(),
	}

	if err := otaRepo.RecordReport(ctx, rep); err != nil {
		t.Fatalf("RecordReport failed: %v", err)
	}

	var updatedTask model.DeviceUpgradeStatus
	if err := db.First(&updatedTask, 1).Error; err != nil {
		t.Fatalf("fetch updated task failed: %v", err)
	}
	if updatedTask.Status != string(event.OTAStatusInProgress) || updatedTask.Progress != 50 {
		t.Errorf("unexpected task state: status=%s, progress=%d", updatedTask.Status, updatedTask.Progress)
	}

	// 2. Report inform with target version matching (success)
	informRep := event.OTAUpgradeReport{
		EventType:       event.OTAEventTypeInform,
		BatchID:         "batch-01",
		DeviceKey:       "dev_ota_01",
		ReportedVersion: "v2.0.0",
		ReportedAt:      time.Now().UTC(),
	}
	if err := otaRepo.RecordReport(ctx, informRep); err != nil {
		t.Fatalf("RecordReport inform failed: %v", err)
	}

	if err := db.First(&updatedTask, 1).Error; err != nil {
		t.Fatalf("fetch updated task failed: %v", err)
	}
	if updatedTask.Status != string(event.OTAStatusSuccess) || updatedTask.CurrentVersion != "v2.0.0" {
		t.Errorf("unexpected task state after inform: status=%s, ver=%s", updatedTask.Status, updatedTask.CurrentVersion)
	}

	var updatedBatch model.UpgradeBatch
	if err := db.First(&updatedBatch, 1).Error; err != nil {
		t.Fatalf("fetch updated batch failed: %v", err)
	}
	if updatedBatch.Status != string(event.OTAStatusSuccess) {
		t.Errorf("expected batch status success, got %s", updatedBatch.Status)
	}
}

func TestOTARepository_RecordReport_Error(t *testing.T) {
	ctx := context.Background()
	otaRepo, db := newTestOTARepository(t)

	dev := &model.Device{ID: 2, DeviceKey: "dev_ota_02", Name: "d2"}
	if err := db.Create(dev).Error; err != nil {
		t.Fatalf("seed device failed: %v", err)
	}
	batch := &model.UpgradeBatch{ID: 2, BatchID: "batch-02", Status: string(event.OTAStatusInProgress)}
	if err := db.Create(batch).Error; err != nil {
		t.Fatalf("seed batch failed: %v", err)
	}
	task := &model.DeviceUpgradeStatus{
		ID:             2,
		DeviceID:       2,
		UpgradeBatchID: "batch-02",
		Status:         string(event.OTAStatusInProgress),
		TargetVersion:  "v2.0.0",
	}
	if err := db.Create(task).Error; err != nil {
		t.Fatalf("seed task failed: %v", err)
	}

	errRep := event.OTAUpgradeReport{
		EventType:  event.OTAEventTypeProgress,
		BatchID:    "batch-02",
		DeviceKey:  "dev_ota_02",
		ReportedAt: time.Now().UTC(),
		Error: &event.OTAReportError{
			Code:    "DOWNLOAD_FAILED",
			Message: "checksum mismatch",
		},
	}

	if err := otaRepo.RecordReport(ctx, errRep); err != nil {
		t.Fatalf("RecordReport with error failed: %v", err)
	}

	var updatedTask model.DeviceUpgradeStatus
	if err := db.First(&updatedTask, 2).Error; err != nil {
		t.Fatalf("fetch updated task failed: %v", err)
	}
	if updatedTask.Status != string(event.OTAStatusFailed) {
		t.Errorf("expected status failed, got %s", updatedTask.Status)
	}
	if updatedTask.LastDispatchError != "DOWNLOAD_FAILED: checksum mismatch" {
		t.Errorf("unexpected last dispatch error: %s", updatedTask.LastDispatchError)
	}

	var updatedBatch model.UpgradeBatch
	if err := db.First(&updatedBatch, 2).Error; err != nil {
		t.Fatalf("fetch updated batch failed: %v", err)
	}
	if updatedBatch.Status != string(event.OTAStatusFailed) {
		t.Errorf("expected batch status failed, got %s", updatedBatch.Status)
	}
}

func TestOTARepository_RecordReport_Partial(t *testing.T) {
	ctx := context.Background()
	otaRepo, db := newTestOTARepository(t)

	// Seed 2 devices in 1 batch: dev-1 succeeds, dev-2 fails
	_ = db.Create(&model.Device{ID: 10, DeviceKey: "dev-10", Name: "d10"})
	_ = db.Create(&model.Device{ID: 11, DeviceKey: "dev-11", Name: "d11"})
	_ = db.Create(&model.UpgradeBatch{ID: 10, BatchID: "batch-10", Status: string(event.OTAStatusInProgress)})

	_ = db.Create(&model.DeviceUpgradeStatus{
		ID: 10, DeviceID: 10, UpgradeBatchID: "batch-10", Status: string(event.OTAStatusPending), TargetVersion: "v1.1",
	})
	_ = db.Create(&model.DeviceUpgradeStatus{
		ID: 11, DeviceID: 11, UpgradeBatchID: "batch-10", Status: string(event.OTAStatusPending), TargetVersion: "v1.1",
	})

	// Dev 10 completes
	_ = otaRepo.RecordReport(ctx, event.OTAUpgradeReport{
		EventType: event.OTAEventTypeInform, BatchID: "batch-10", DeviceKey: "dev-10", ReportedVersion: "v1.1",
	})
	// Dev 11 fails
	_ = otaRepo.RecordReport(ctx, event.OTAUpgradeReport{
		EventType: event.OTAEventTypeProgress, BatchID: "batch-10", DeviceKey: "dev-11",
		Error: &event.OTAReportError{Code: "FLASH_ERROR", Message: "hardware write failure"},
	})

	var batch model.UpgradeBatch
	if err := db.First(&batch, 10).Error; err != nil {
		t.Fatalf("fetch batch failed: %v", err)
	}
	if batch.Status != "partial" {
		t.Errorf("expected batch status 'partial', got %s", batch.Status)
	}
}
