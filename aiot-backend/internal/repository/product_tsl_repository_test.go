package repository

import (
	"context"
	"testing"

	"0things-backend/internal/model"
	"github.com/stretchr/testify/require"
)

func TestProductTSLRepository(t *testing.T) {
	store := newRepositoryTestDB(t, &model.ProductTSL{})
	repo := NewProductTSLRepository(store)
	ctx := context.Background()
	productID := int64(1)
	tsl := &model.ProductTSL{ProductID: &productID, TSL: `{"properties":[]}`}
	require.NoError(t, repo.Create(ctx, tsl))
	stored, err := repo.FindByProductID(ctx, productID)
	require.NoError(t, err)
	require.Equal(t, tsl.TSL, stored.TSL)
	require.NoError(t, repo.Delete(ctx, stored))
	_, err = repo.FindByProductID(ctx, productID)
	require.ErrorIs(t, err, ErrNotFound)
}
