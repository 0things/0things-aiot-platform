package model_test

import (
	"testing"

	"aiot-backend/internal/model"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestRuleNodeDefinitionUniqueKeys(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.RuleNodeDefinition{}))

	definition := model.RuleNodeDefinition{
		UUID:          "00000000-0000-0000-0000-000000000001",
		Key:           "filter.value-condition",
		Version:       1,
		Category:      "Filter",
		Name:          "数值条件",
		ConfigSchema:  `{}`,
		DefaultConfig: `{}`,
		InputPorts:    `[]`,
		OutputPorts:   `[]`,
		ExecutorKey:   "filter-value-condition-v1",
		Enabled:       true,
	}
	require.NoError(t, db.Create(&definition).Error)

	duplicateKey := definition
	duplicateKey.ID = 0
	duplicateKey.UUID = "00000000-0000-0000-0000-000000000002"
	require.Error(t, db.Create(&duplicateKey).Error)

	duplicateUUID := definition
	duplicateUUID.ID = 0
	duplicateUUID.Key = "filter.property"
	require.Error(t, db.Create(&duplicateUUID).Error)
}
