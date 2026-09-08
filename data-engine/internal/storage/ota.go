package storage

import (
	"context"
	"fmt"
	"time"

	"0things/pkg/event"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// OTAStore writes data-engine's OTA results to backend-owned tables. It never
// creates or migrates schema.
type OTAStore struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewOTAStore(config *viper.Viper, logger *zap.Logger) (*OTAStore, error) {
	driver, dsn := config.GetString("data.db.aiot.driver"), config.GetString("data.db.aiot.dsn")
	if driver == "" || dsn == "" {
		return nil, fmt.Errorf("data.db.aiot.driver and data.db.aiot.dsn are required")
	}
	var dialect gorm.Dialector
	switch driver {
	case "mysql":
		dialect = mysql.Open(dsn)
	case "postgres", "postgresql":
		dialect = postgres.Open(dsn)
	case "sqlite":
		dialect = sqlite.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported ota database driver %q", driver)
	}
	db, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open ota database: %w", err)
	}
	return &OTAStore{db: db, logger: logger}, nil
}

func (s *OTAStore) RecordReport(ctx context.Context, report event.OTAUpgradeReport) error {
	now := report.ReportedAt
	if now.IsZero() {
		now = time.Now().UTC()
	}
	updates := map[string]interface{}{"last_report_at": now.Unix()}
	if report.EventType == event.OTAEventTypeProgress {
		updates["status"] = string(event.OTAStatusInProgress)
		updates["last_status_change_ts"] = now.Unix()
		updates["first_progress_at"] = gorm.Expr("COALESCE(first_progress_at, ?)", now.Unix())
		if report.Progress != nil {
			updates["progress"] = *report.Progress
		}
	}
	if report.EventType == event.OTAEventTypeInform && report.ReportedVersion != "" {
		updates["current_version"] = report.ReportedVersion
		var targetVersion string
		if err := s.db.WithContext(ctx).Table("ota_device_upgrade_status AS task").Joins("JOIN devices ON devices.id = task.device_id").Where("task.upgrade_batch_id = ? AND devices.device_key = ?", report.BatchID, report.DeviceKey).Pluck("task.target_version", &targetVersion).Error; err != nil {
			return fmt.Errorf("read OTA target version: %w", err)
		}
		if targetVersion == report.ReportedVersion {
			updates["status"] = string(event.OTAStatusSuccess)
			updates["last_status_change_ts"] = now.Unix()
		}
	}
	if report.Error != nil {
		updates["status"] = string(event.OTAStatusFailed)
		updates["last_dispatch_error"] = report.Error.Code + ": " + report.Error.Message
		updates["last_status_change_ts"] = now.Unix()
	}
	return s.updateDeviceTask(ctx, report.BatchID, report.DeviceKey, updates)
}

func (s *OTAStore) updateDeviceTask(ctx context.Context, batchID, deviceKey string, updates map[string]interface{}) error {
	if batchID == "" || deviceKey == "" {
		return fmt.Errorf("batch_id and device_key are required")
	}
	result := s.db.WithContext(ctx).Table("ota_device_upgrade_status AS task").Joins("JOIN devices ON devices.id = task.device_id").Where("task.upgrade_batch_id = ? AND devices.device_key = ?", batchID, deviceKey).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("update ota device task: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("ota device task not found for batch %s and device %s", batchID, deviceKey)
	}
	return s.refreshBatchStatus(ctx, batchID)
}

func (s *OTAStore) refreshBatchStatus(ctx context.Context, batchID string) error {
	var total, succeeded, failed, active int64
	row := s.db.WithContext(ctx).Table("ota_device_upgrade_status").Select("COUNT(*) AS total, SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) AS succeeded, SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) AS failed, SUM(CASE WHEN status IN (?, ?, ?) THEN 1 ELSE 0 END) AS active", string(event.OTAStatusSuccess), string(event.OTAStatusFailed), string(event.OTAStatusPending), string(event.OTAStatusSent), string(event.OTAStatusInProgress)).Where("upgrade_batch_id = ?", batchID).Row()
	if err := row.Scan(&total, &succeeded, &failed, &active); err != nil {
		return fmt.Errorf("read ota batch state: %w", err)
	}
	if total == 0 {
		return nil
	}
	status := string(event.OTAStatusInProgress)
	if active == 0 && succeeded == total {
		status = string(event.OTAStatusSuccess)
	} else if active == 0 && failed == total {
		status = string(event.OTAStatusFailed)
	} else if active == 0 {
		status = "partial"
	}
	return s.db.WithContext(ctx).Table("ota_upgrade_batches").Where("batch_id = ?", batchID).Update("status", status).Error
}
