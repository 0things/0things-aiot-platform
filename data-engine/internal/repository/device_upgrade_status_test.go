package repository

import (
	"context"
	"testing"

	"data-engine/internal/event"
	"data-engine/internal/model"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func newTestDeviceUpgradeStatusRepository(t *testing.T) (DeviceUpgradeStatusRepository, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.Device{}, &model.OTADeviceUpgradeStatus{})
	require.NoError(t, err)

	repo := NewRepository(zap.NewNop(), db)
	return NewDeviceUpgradeStatusRepository(repo), db
}

func TestDeviceUpgradeStatusRepository_UpdateProgressAndDeviceInfo(t *testing.T) {
	ctx := context.Background()
	repo, db := newTestDeviceUpgradeStatusRepository(t)

	// Seed records with TargetVersion = "v2.0.0"
	require.NoError(t, db.Create(&model.Device{ID: 1, DeviceKey: "dev_01"}).Error)
	require.NoError(t, db.Create(&model.OTADeviceUpgradeStatus{ID: 1, DeviceID: 1, UpgradeBatchID: "b_01", Status: string(event.OTAStatusPending), TargetVersion: "v2.0.0"}).Error)

	// 1. Intermediate progress 100% without version report: MUST remain in_progress (waiting for inform reboot)
	err := repo.UpdateProgress(ctx, "b_01", "dev_01", "", 100)
	require.NoError(t, err)

	var task model.OTADeviceUpgradeStatus
	require.NoError(t, db.First(&task, 1).Error)
	assert.Equal(t, string(event.OTAStatusInProgress), task.Status, "100% progress without version must remain in_progress")
	assert.Equal(t, int32(100), task.Progress)

	// 2. Inform report matching TargetVersion ("v2.0.0" == "v2.0.0" -> success)
	err = repo.UpdateDeviceInfo(ctx, "dev_01", "v2.0.0")
	require.NoError(t, err)

	require.NoError(t, db.First(&task, 1).Error)
	assert.Equal(t, string(event.OTAStatusSuccess), task.Status)
	assert.Equal(t, "v2.0.0", task.CurrentVersion)

	// 3. Mismatched version report -> failed
	require.NoError(t, db.Create(&model.Device{ID: 2, DeviceKey: "dev_02"}).Error)
	require.NoError(t, db.Create(&model.OTADeviceUpgradeStatus{ID: 2, DeviceID: 2, UpgradeBatchID: "b_02", Status: string(event.OTAStatusInProgress), TargetVersion: "v2.0.0"}).Error)

	err = repo.UpdateDeviceInfo(ctx, "dev_02", "v1.0.0")
	require.NoError(t, err)
	var task2 model.OTADeviceUpgradeStatus
	require.NoError(t, db.First(&task2, 2).Error)
	assert.Equal(t, string(event.OTAStatusFailed), task2.Status, "mismatched version must be marked failed")

	// 4. Error reporting (errDesc -> failed)
	err = repo.UpdateProgress(ctx, "b_01", "dev_01", "flash error", 50)
	require.NoError(t, err)
	require.NoError(t, db.First(&task, 1).Error)
	assert.Equal(t, string(event.OTAStatusFailed), task.Status)

	// 5. Non-existent device
	err = repo.UpdateProgress(ctx, "b_01", "non_existent", "", 0)
	assert.Error(t, err)
}
