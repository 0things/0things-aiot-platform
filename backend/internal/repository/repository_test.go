package repository

import (
	"testing"

	"aiot-backend/internal/model"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newRepositoryTestDB(t *testing.T, models ...any) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)
	// 产品查询会加载协议组合，测试库也需具备对应表结构。
	models = append(models, &model.ProductProtocol{})
	require.NoError(t, db.AutoMigrate(models...))
	return db
}
