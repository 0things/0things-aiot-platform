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
	require.NoError(t, db.AutoMigrate(&model.Category{}, &model.RuleNodeDefinition{}))

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

	var definitions []model.RuleNodeDefinition
	require.NoError(t, db.Order("category ASC, sort_order ASC").Find(&definitions).Error)
	require.Len(t, definitions, 102)

	var systemEntry model.RuleNodeDefinition
	require.NoError(t, db.Where("key = ?", "system.entry").First(&systemEntry).Error)
	assert.True(t, systemEntry.IsSystem)

	expectedCategoryCounts := map[string]int64{
		"Filter":         16,
		"Enrichment":     14,
		"Transformation": 11,
		"Action":         33,
		"External":       15,
		"Flow":           10,
		"Analytics":      2,
	}
	for category, expectedCount := range expectedCategoryCounts {
		var categoryCount int64
		require.NoError(t, db.Model(&model.RuleNodeDefinition{}).
			Where("category = ? AND is_system = ?", category, false).
			Count(&categoryCount).Error)
		assert.Equal(t, expectedCount, categoryCount, category)
	}

	var valueCondition model.RuleNodeDefinition
	require.NoError(t, db.Where("key = ?", "filter.value-condition").First(&valueCondition).Error)
	require.Equal(t, "Filter", valueCondition.Category)
	require.NotEmpty(t, valueCondition.ConfigSchema)
	require.NotEmpty(t, valueCondition.InputPorts)
	require.NotEmpty(t, valueCondition.OutputPorts)

	valueCondition.Name = "自定义数值条件"
	valueCondition.ConfigSchema = `{}`
	require.NoError(t, db.Save(&valueCondition).Error)
	require.NoError(t, seeder.SeedDefaults(ctx, db))

	var reseeded model.RuleNodeDefinition
	require.NoError(t, db.Where("key = ?", "filter.value-condition").First(&reseeded).Error)
	assert.Equal(t, "自定义数值条件", reseeded.Name)
	assert.Contains(t, reseeded.ConfigSchema, `"required":["path","operator","value"]`)
}
