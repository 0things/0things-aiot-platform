package repository

import (
	"context"
	"fmt"
	"time"

	"data-engine/internal/event"

	"gorm.io/gen/field"
)

// DeviceUpgradeStatusRepository manages device OTA upgrade status persistence.
type DeviceUpgradeStatusRepository interface {
	// UpdateProgress updates the device's download/install progress.
	UpdateProgress(ctx context.Context, batchID, deviceKey, errDesc string, progress int32) error
	// UpdateDeviceInfo updates the device's reported version and validates upgrade success against TargetVersion.
	UpdateDeviceInfo(ctx context.Context, deviceKey, version string) error
}

type deviceUpgradeStatusRepository struct {
	*Repository
}

// NewDeviceUpgradeStatusRepository creates a new DeviceUpgradeStatusRepository.
func NewDeviceUpgradeStatusRepository(repo *Repository) DeviceUpgradeStatusRepository {
	return &deviceUpgradeStatusRepository{Repository: repo}
}

// UpdateProgress updates device progress and error description using GORM gen.
func (r *deviceUpgradeStatusRepository) UpdateProgress(ctx context.Context, batchID, deviceKey, errDesc string, progress int32) error {
	q := useQuery(r.db)

	device, err := q.Device.WithContext(ctx).Where(q.Device.DeviceKey.Eq(deviceKey)).First()
	if err != nil {
		return fmt.Errorf("device not found %s: %w", deviceKey, err)
	}

	task, err := q.OTADeviceUpgradeStatus.WithContext(ctx).Where(
		q.OTADeviceUpgradeStatus.UpgradeBatchID.Eq(batchID),
		q.OTADeviceUpgradeStatus.DeviceID.Eq(device.ID),
	).First()
	if err != nil {
		return fmt.Errorf("task not found for batch %s device %s: %w", batchID, deviceKey, err)
	}

	status := task.Status
	if errDesc != "" {
		status = string(event.OTAStatusFailed)
	} else if progress > 0 {
		status = string(event.OTAStatusInProgress)
	}

	now := time.Now().Unix()
	updates := []field.AssignExpr{
		q.OTADeviceUpgradeStatus.Status.Value(status),
		q.OTADeviceUpgradeStatus.Progress.Value(progress),
		q.OTADeviceUpgradeStatus.LastDispatchError.Value(errDesc),
		q.OTADeviceUpgradeStatus.LastStatusChangeTime.Value(now),
		q.OTADeviceUpgradeStatus.LastReportAt.Value(now),
	}
	if task.FirstProgressAt == nil {
		updates = append(updates, q.OTADeviceUpgradeStatus.FirstProgressAt.Value(now))
	}

	_, err = q.OTADeviceUpgradeStatus.WithContext(ctx).Where(q.OTADeviceUpgradeStatus.ID.Eq(task.ID)).UpdateSimple(updates...)
	if err != nil {
		return fmt.Errorf("update device upgrade progress: %w", err)
	}

	return nil
}

// UpdateDeviceInfo updates the reported version from inform, and determines success/failure based on TargetVersion.
func (r *deviceUpgradeStatusRepository) UpdateDeviceInfo(ctx context.Context, deviceKey, version string) error {
	q := useQuery(r.db)

	device, err := q.Device.WithContext(ctx).Where(q.Device.DeviceKey.Eq(deviceKey)).First()
	if err != nil {
		return nil
	}

	task, err := q.OTADeviceUpgradeStatus.WithContext(ctx).Where(
		q.OTADeviceUpgradeStatus.DeviceID.Eq(device.ID),
		q.OTADeviceUpgradeStatus.Status.In(string(event.OTAStatusPending), string(event.OTAStatusSent), string(event.OTAStatusInProgress)),
	).Order(q.OTADeviceUpgradeStatus.ID.Desc()).First()
	if err != nil {
		return nil
	}

	status := string(event.OTAStatusFailed)
	progress := task.Progress
	if task.TargetVersion != "" && version == task.TargetVersion {
		status = string(event.OTAStatusSuccess)
		progress = 100
	}

	now := time.Now().Unix()
	_, err = q.OTADeviceUpgradeStatus.WithContext(ctx).Where(q.OTADeviceUpgradeStatus.ID.Eq(task.ID)).UpdateSimple(
		q.OTADeviceUpgradeStatus.Status.Value(status),
		q.OTADeviceUpgradeStatus.Progress.Value(progress),
		q.OTADeviceUpgradeStatus.CurrentVersion.Value(version),
		q.OTADeviceUpgradeStatus.LastStatusChangeTime.Value(now),
		q.OTADeviceUpgradeStatus.LastReportAt.Value(now),
	)
	if err != nil {
		return fmt.Errorf("update device info status: %w", err)
	}

	return nil
}
