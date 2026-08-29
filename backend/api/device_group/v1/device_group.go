package devicegroupv1

import (
	devicev1 "aiot-backend/api/device/v1"
	"time"
)

type CreateDeviceGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Description string `json:"description"`
	Rule        string `json:"rule"`
}

type UpdateDeviceGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Rule        string `json:"rule"`
}

type DeviceGroup struct {
	GroupUUID   string     `json:"groupUuid"`
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	Description string     `json:"description"`
	Rule        string     `json:"rule,omitempty"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type ListDeviceGroupsResponse struct {
	Items    []DeviceGroup `json:"items"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"pageSize"`
}

type DeviceKeysRequest struct {
	DeviceKeys []string `json:"deviceKeys" binding:"required,min=1"`
}

type DeviceGroupDevicesResponse struct {
	Devices  []devicev1.Device `json:"devices"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
}

type PreviewRequest struct {
	Rule string `json:"rule" binding:"required"`
}

type PreviewResponse struct {
	Total   int64             `json:"total"`
	Devices []devicev1.Device `json:"devices"`
}
