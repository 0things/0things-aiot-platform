package model

import (
	"time"

	"gorm.io/gorm"
)

// OTAPackage represents an OTA firmware/software package (table: ota_packages).
type OTAPackage struct {
	ID             int64          `gorm:"primaryKey"`
	UUID           string         `gorm:"column:uuid;type:varchar(64);uniqueIndex"`
	PackageName    string         `gorm:"column:package_name"`
	Version        string         `gorm:"column:version"`
	ProductID      int64          `gorm:"column:product_id"`
	OrganizationID int64          `gorm:"column:organization_id;not null;default:1"`
	ProductKey     string         `gorm:"column:product_key;->" json:"productKey,omitempty"`
	ProductName    string         `gorm:"column:product_name;->" json:"productName,omitempty"`
	PackageType    string         `gorm:"column:package_type"`
	Status         string         `gorm:"column:status"`
	UploadType     string         `gorm:"column:upload_type"`
	FileURL        string         `gorm:"column:file_url"`
	FileSize       int64          `gorm:"column:file_size"`
	Checksum       string         `gorm:"column:checksum"`
	Description    string         `gorm:"column:description"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at"`
	CreatedAt      time.Time      `gorm:"column:created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at"`
	ReleasedAt     *time.Time     `gorm:"column:released_at"`
}

func (OTAPackage) TableName() string { return "ota_packages" }
