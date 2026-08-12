package repository

import (
	"context"
	"testing"

	"0things-backend/internal/model"
	"github.com/stretchr/testify/require"
)

func TestOTARepository(t *testing.T) {
	store := newRepositoryTestDB(t, &model.OTAPackage{}, &model.DeviceUpgradeStatus{}, &model.UpgradeBatch{}, &model.Product{})
	repo := NewOTARepository(store)
	ctx := context.Background()
	pkg := &model.OTAPackage{PackageName: "firmware-1", Version: "1.0.0", Status: "draft"}
	require.NoError(t, repo.Create(ctx, pkg))
	found, err := repo.FindByName(ctx, "firmware-1")
	require.NoError(t, err)
	found.Status = "deploying"
	require.NoError(t, repo.Save(ctx, found))
	require.NoError(t, store.WithContext(ctx).Create(&model.DeviceUpgradeStatus{DeviceID: 1, OTAPackageID: "1", Status: "success"}).Error)
	require.NoError(t, store.WithContext(ctx).Create(&model.UpgradeBatch{BatchID: "B001", OTAPackageID: "1", Status: "pending"}).Error)

	items, total, err := repo.List(ctx, 1, 20)
	require.NoError(t, err)
	require.EqualValues(t, 1, total)
	require.Len(t, items, 1)
	stats, err := repo.Statistics(ctx, pkg.ID)
	require.NoError(t, err)
	require.EqualValues(t, 1, stats.Success)
	batches, err := repo.Batches(ctx, pkg.ID)
	require.NoError(t, err)
	require.Len(t, batches, 1)
	require.NoError(t, repo.Delete(ctx, pkg.ID))
	_, err = repo.Find(ctx, pkg.ID)
	require.ErrorIs(t, err, ErrNotFound)
}
