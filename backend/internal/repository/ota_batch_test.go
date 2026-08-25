package repository

import (
	"context"
	"testing"

	"aiot-backend/internal/model"
	"github.com/stretchr/testify/require"
)

func TestOTABatchRepository(t *testing.T) {
	store := newRepositoryTestDB(t, &model.OTAPackage{}, &model.DeviceUpgradeStatus{}, &model.UpgradeBatch{})
	repo := NewOTARepository(store)
	ctx := context.Background()

	pkg := &model.OTAPackage{PackageName: "firmware-batch", Version: "2.0.0", ProductID: 1, Status: "draft"}
	require.NoError(t, repo.Create(ctx, pkg))
	const pkgID = int64(1) // first inserted row id in the in-memory db

	// 第一批：3 台设备
	batch := &model.UpgradeBatch{
		BatchID:           "B-1-001",
		OTAPackageID:      "1",
		BatchName:         "批量升级-1",
		UpgradeStrategy:   "immediate",
		Status:            "pending",
		TargetDeviceCount: 3,
	}
	require.NoError(t, repo.CreateBatch(ctx, batch))

	created, err := repo.CreateBatchDeployments(ctx, pkgID, "B-1-001", []int64{1, 2, 3})
	require.NoError(t, err)
	require.Equal(t, 3, created)

	// 再次提交相同设备到新批次：已存在记录应更新为 pending 并关联新批次，不重复插入
	require.NoError(t, repo.CreateBatch(ctx, &model.UpgradeBatch{
		BatchID: "B-1-002", OTAPackageID: "1", BatchName: "批量升级-2", Status: "pending",
	}))
	created2, err := repo.CreateBatchDeployments(ctx, pkgID, "B-1-002", []int64{2, 3, 4})
	require.NoError(t, err)
	require.Equal(t, 3, created2)

	// 校验：每台设备仅一条记录，且 device 2/3 已关联到第二个批次
	var records []model.DeviceUpgradeStatus
	require.NoError(t, store.WithContext(ctx).Where("ota_package_id = ?", "1").Find(&records).Error)
	require.Len(t, records, 4)

	var d2 model.DeviceUpgradeStatus
	require.NoError(t, store.WithContext(ctx).Where("ota_package_id = ? AND device_id = ?", "1", 2).First(&d2).Error)
	require.Equal(t, "B-1-002", d2.UpgradeBatchID)
	require.Equal(t, "pending", d2.Status)
}
