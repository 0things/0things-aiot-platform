package seeder_test

import (
	"context"
	"testing"

	"aiot-backend/internal/model"
	"aiot-backend/internal/seeder"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestSeedDefaults(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Category{}))

	ctx := context.Background()

	// First run should seed default categories
	err = seeder.SeedDefaults(ctx, db)
	require.NoError(t, err)

	var categories []model.Category
	err = db.Order("sort ASC").Find(&categories).Error
	require.NoError(t, err)
	assert.Equal(t, 7, len(categories))
	assert.Equal(t, "传感器", categories[0].Name)

	// Second run should be idempotent (no duplicates)
	err = seeder.SeedDefaults(ctx, db)
	require.NoError(t, err)

	var count int64
	err = db.Model(&model.Category{}).Count(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(7), count)
}
