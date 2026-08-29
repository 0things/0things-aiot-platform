package repository

import (
	"context"
	"testing"

	"aiot-backend/internal/model"
	"aiot-backend/internal/tenant"

	"github.com/stretchr/testify/require"
)

func TestProductRepository(t *testing.T) {
	store := newRepositoryTestDB(t, &model.Product{}, &model.Category{})
	repo := NewProductRepository(store)
	ctx := context.Background()

	categoryID := int64(10)
	require.NoError(t, store.WithContext(ctx).Create(&model.Category{ID: categoryID, Name: "Sensors"}).Error)
	product := &model.Product{ProductKey: "P001", Name: "sensor", CategoryID: &categoryID, Status: "active", OrganizationID: 1}
	require.NoError(t, repo.Create(ctx, product))
	found, err := repo.FindByKey(ctx, "P001")
	require.NoError(t, err)
	require.Equal(t, "sensor", found.Name)
	found.Name = "updated sensor"
	require.NoError(t, repo.Save(ctx, found))
	items, total, err := repo.List(ctx, 1, 20, "10", "active", "updated")
	require.NoError(t, err)
	require.EqualValues(t, 1, total)
	require.Len(t, items, 1)
	require.Equal(t, "Sensors", items[0].CategoryName)

	require.NoError(t, repo.Create(ctx, &model.Product{ProductKey: "P002", Name: "updated other tenant", OrganizationID: 2}))
	items, total, err = repo.List(ctx, 1, 20, "", "", "updated")
	require.NoError(t, err)
	require.EqualValues(t, 1, total)
	require.Len(t, items, 1)

	items, total, err = repo.List(tenant.WithTenant(ctx, 2), 1, 20, "", "", "updated")
	require.NoError(t, err)
	require.EqualValues(t, 1, total)
	require.Len(t, items, 1)
}
