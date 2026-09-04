package v1

import "time"

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

type ListDeviceGroupsRequest struct {
	PageRequest
	Search string `form:"search"`
	Type   string `form:"type"`
}

type ListDeviceGroupDevicesRequest struct {
	PageRequest
	ProductKey string `form:"productKey"`
	Search     string `form:"search"`
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
	Devices  []Device `json:"devices"`
	Total    int64    `json:"total"`
	Page     int      `json:"page"`
	PageSize int      `json:"pageSize"`
}

type PreviewRequest struct {
	Rule string `json:"rule" binding:"required"`
}

type PreviewResponse struct {
	Total   int64    `json:"total"`
	Devices []Device `json:"devices"`
}
