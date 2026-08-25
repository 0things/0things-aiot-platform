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

	// 再次提交相同设备到新批次：为每个设备新增一条独立的批次记录
	require.NoError(t, repo.CreateBatch(ctx, &model.UpgradeBatch{
		BatchID: "B-1-002", OTAPackageID: "1", BatchName: "批量升级-2", Status: "pending",
	}))
	created2, err := repo.CreateBatchDeployments(ctx, pkgID, "B-1-002", []int64{2, 3, 4})
	require.NoError(t, err)
	require.Equal(t, 3, created2)

	// 校验：device 2/3 各保留两条记录，分别属于两个批次；device 1/4 各一条
	var records []model.DeviceUpgradeStatus
	require.NoError(t, store.WithContext(ctx).Where("ota_package_id = ?", "1").Find(&records).Error)
	require.Len(t, records, 6)

	var d2 []model.DeviceUpgradeStatus
	require.NoError(t, store.WithContext(ctx).Where("ota_package_id = ? AND device_id = ?", "1", 2).Find(&d2).Error)
	require.Len(t, d2, 2)
	require.ElementsMatch(t, []string{"B-1-001", "B-1-002"}, []string{d2[0].UpgradeBatchID, d2[1].UpgradeBatchID})
}
