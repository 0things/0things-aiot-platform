package seeder

import (
	"context"

	"aiot-backend/internal/model"

	"gorm.io/gorm"
)

// SeedDefaults initializes default system metadata if the database is empty.
func SeedDefaults(ctx context.Context, db *gorm.DB) error {
	if err := seedDefaultCategories(ctx, db); err != nil {
		return err
	}
	return nil
}

// seedDefaultCategories populates the category table with default categories when empty.
func seedDefaultCategories(ctx context.Context, db *gorm.DB) error {
	var count int64
	if err := db.WithContext(ctx).Model(&model.Category{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	items := []string{"传感器", "执行器", "网关", "控制器", "显示设备", "摄像头", "其他"}
	for i, name := range items {
		if err := db.WithContext(ctx).Create(&model.Category{
			Name:    name,
			Sort:    i,
			Enabled: true,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}
