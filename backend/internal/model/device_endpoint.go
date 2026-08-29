package model

import "time"

// DeviceEndpoint 保存设备实际使用的协议端点和连接状态。
type DeviceEndpoint struct {
	ID                int64      `gorm:"column:id;primaryKey" json:"id"`
	DeviceID          int64      `gorm:"column:device_id;not null;index;uniqueIndex:idx_device_endpoint_protocol" json:"deviceId"`
	ProductProtocolID int64      `gorm:"column:product_protocol_id;not null;index;uniqueIndex:idx_device_endpoint_protocol" json:"productProtocolId"`
	Endpoint          string     `gorm:"column:endpoint;size:512" json:"endpoint"`
	Credentials       string     `gorm:"column:credentials;type:text" json:"-"`
	ProtocolConfig    string     `gorm:"column:protocol_config;type:text" json:"protocolConfig"`
	Enabled           bool       `gorm:"column:enabled;not null;default:true" json:"enabled"`
	Status            string     `gorm:"column:status;size:32;not null;default:inactive" json:"status"`
	LastSeenAt        *time.Time `gorm:"column:last_seen_at" json:"lastSeenAt"`
	LastError         string     `gorm:"column:last_error;type:text" json:"lastError"`
	CreatedAt         time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt         time.Time  `gorm:"column:updated_at" json:"updatedAt"`
}

func (DeviceEndpoint) TableName() string { return "device_endpoints" }
