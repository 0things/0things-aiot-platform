package repository

import (
	"context"
	"testing"

	"aiot-backend/internal/model"
	"aiot-backend/internal/tenant"
	"github.com/stretchr/testify/require"
)

func TestProductRepository(t *testing.T) {
	store := newRepositoryTestDB(t, &model.Product{})
	repo := NewProductRepository(store)
	ctx := context.Background()

	product := &model.Product{ProductKey: "P001", Name: "sensor", Category: "environment", Status: "active", TenantID: 1}
	require.NoError(t, repo.Create(ctx, product))
	found, err := repo.FindByKey(ctx, "P001")
	require.NoError(t, err)
	require.Equal(t, "sensor", found.Name)
	found.Name = "updated sensor"
	require.NoError(t, repo.Save(ctx, found))
	items, total, err := repo.List(ctx, 1, 20, "environment", "active", "updated")
	require.NoError(t, err)
	require.EqualValues(t, 1, total)
	require.Len(t, items, 1)

	require.NoError(t, repo.Create(ctx, &model.Product{ProductKey: "P002", Name: "updated other tenant", TenantID: 2}))
	items, total, err = repo.List(ctx, 1, 20, "", "", "updated")
	require.NoError(t, err)
	require.EqualValues(t, 1, total)
	require.Len(t, items, 1)

	items, total, err = repo.List(tenant.WithTenant(ctx, 2), 1, 20, "", "", "updated")
	require.NoError(t, err)
	require.EqualValues(t, 1, total)
	require.Len(t, items, 1)
}
