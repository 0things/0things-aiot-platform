package model

import "time"

// DeviceCredential stores the MQTT credentials used by a device.
type DeviceCredential struct {
	ID                 int64     `gorm:"column:id;primaryKey" json:"id"`
	DeviceUUID         string    `gorm:"column:device_uuid;size:36;not null;uniqueIndex" json:"deviceUuid"`
	CredentialType     string    `gorm:"column:credential_type;size:32;not null;default:mqtt" json:"credentialType"`
	Username           string    `gorm:"column:username;size:255;not null;uniqueIndex" json:"username"`
	Password           string    `gorm:"column:password;size:255;not null" json:"-"`
	Salt               string    `gorm:"column:salt;size:128;not null;default:''" json:"-"`
	PasswordCiphertext string    `gorm:"column:password_ciphertext;type:text" json:"-"`
	Enabled            bool      `gorm:"column:enabled;not null;default:true" json:"enabled"`
	CreatedAt          time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt          time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (DeviceCredential) TableName() string { return "device_credentials" }
