package repository

import (
	"context"
	"fmt"
	"time"

	"0things/pkg/event"
	"data-engine/internal/dal/query"
	"data-engine/internal/model"

	"gorm.io/gorm"
)

// OTARepository defines the interface for OTA progress and batch status updates.
type OTARepository interface {
	RecordReport(ctx context.Context, report event.OTAUpgradeReport) error
}

type otaRepository struct {
	*Repository
}

func NewOTARepository(repo *Repository) OTARepository {
	return &otaRepository{Repository: repo}
}

func (r *otaRepository) RecordReport(ctx context.Context, report event.OTAUpgradeReport) error {
	now := report.ReportedAt
	if now.IsZero() {
		now = time.Now().UTC()
	}

	updates := map[string]interface{}{
		"last_report_at": now.Unix(),
	}

	if report.Status != "" {
		updates["status"] = string(report.Status)
		updates["last_status_change_ts"] = now.Unix()
	}

	if report.EventType == event.OTAEventTypeProgress {
		if _, ok := updates["status"]; !ok {
			updates["status"] = string(event.OTAStatusInProgress)
			updates["last_status_change_ts"] = now.Unix()
		}
		updates["first_progress_at"] = gorm.Expr("COALESCE(first_progress_at, ?)", now.Unix())
		if report.Progress != nil {
			updates["progress"] = *report.Progress
		}
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		q := useQuery(tx)
		dev, err := q.Device.WithContext(ctx).Where(q.Device.DeviceKey.Eq(report.DeviceKey)).First()
		if err != nil {
			return fmt.Errorf("read device for OTA report: %w", err)
		}

		if report.EventType == event.OTAEventTypeInform && report.ReportedVersion != "" {
			updates["current_version"] = report.ReportedVersion
			task, err := q.DeviceUpgradeStatus.WithContext(ctx).
				Where(q.DeviceUpgradeStatus.UpgradeBatchID.Eq(report.BatchID), q.DeviceUpgradeStatus.DeviceID.Eq(dev.ID)).
				First()
			if err != nil {
				return fmt.Errorf("read OTA task for inform: %w", err)
			}

			if task.TargetVersion == report.ReportedVersion {
				updates["status"] = string(event.OTAStatusSuccess)
				updates["last_status_change_ts"] = now.Unix()
			}
		}

		if report.Error != nil {
			updates["status"] = string(event.OTAStatusFailed)
			updates["last_dispatch_error"] = report.Error.Code + ": " + report.Error.Message
			updates["last_status_change_ts"] = now.Unix()
		}

		result, err := q.DeviceUpgradeStatus.WithContext(ctx).
			Where(q.DeviceUpgradeStatus.UpgradeBatchID.Eq(report.BatchID), q.DeviceUpgradeStatus.DeviceID.Eq(dev.ID)).
			Updates(updates)
		if err != nil {
			return fmt.Errorf("update ota device task: %w", err)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("ota device task not found for batch %s and device %s", report.BatchID, report.DeviceKey)
		}

		return r.refreshBatchStatusTx(ctx, q, report.BatchID)
	})
}

func (r *otaRepository) refreshBatchStatusTx(ctx context.Context, q *query.Query, batchID string) error {
	var total, succeeded, failed, active int64

	row := q.DeviceUpgradeStatus.WithContext(ctx).
		UnderlyingDB().
		Model(&model.DeviceUpgradeStatus{}).
		Select(
			"COUNT(*) AS total, SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) AS succeeded, SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) AS failed, SUM(CASE WHEN status IN (?, ?, ?) THEN 1 ELSE 0 END) AS active",
			string(event.OTAStatusSuccess),
			string(event.OTAStatusFailed),
			string(event.OTAStatusPending),
			string(event.OTAStatusSent),
			string(event.OTAStatusInProgress),
		).
		Where("upgrade_batch_id = ?", batchID).
		Row()

	if err := row.Scan(&total, &succeeded, &failed, &active); err != nil {
		return fmt.Errorf("read ota batch state: %w", err)
	}
	if total == 0 {
		return nil
	}

	status := string(event.OTAStatusInProgress)
	if active == 0 {
		if succeeded == total {
			status = string(event.OTAStatusSuccess)
		} else if failed == total {
			status = string(event.OTAStatusFailed)
		} else {
			status = "partial"
		}
	}

	_, err := q.UpgradeBatch.WithContext(ctx).
		Where(q.UpgradeBatch.BatchID.Eq(batchID)).
		Update(q.UpgradeBatch.Status, status)
	return err
}
