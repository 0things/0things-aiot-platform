// Package devicev1 owns the device HTTP contract and mirrors device-service.
package devicev1

import (
	"encoding/json"
	"time"
)

type CreateDeviceRequest struct {
	Name      string          `json:"name" binding:"required"`
	ProductID int64           `json:"productId" binding:"required"`
	Enabled   bool            `json:"enabled"`
	Metadata  json.RawMessage `json:"metadata"`
}

type UpdateDeviceRequest struct {
	Name     string          `json:"name"`
	State    string          `json:"state"`
	Metadata json.RawMessage `json:"metadata"`
}

type SetDeviceEnabledRequest struct {
	Enabled bool `json:"enabled"`
}

type SetDeviceTagsRequest struct {
	Tags map[string]string `json:"tags" binding:"required"`
}

type DeleteDeviceTagsRequest struct {
	Keys []string `json:"keys" binding:"required"`
}

type UpdateDesiredShadowRequest struct {
	Desired map[string]any `json:"desired" binding:"required"`
	Version int64          `json:"version"`
}

type UpdateReportedShadowRequest struct {
	Reported map[string]any `json:"reported" binding:"required"`
	Version  int64          `json:"version"`
}

type ClearDesiredShadowRequest struct {
	Version int64 `json:"version"`
}

type SimulatePushRequest struct {
	Payload string `json:"payload"`
}

type MockKafkaRequest struct {
	Topic string `json:"topic" binding:"required"`
	Data  string `json:"data"`
}

type Device struct {
	ID              int64      `json:"id"`
	DeviceKey       string     `json:"deviceKey"`
	Name            string     `json:"name"`
	ProductID       int64      `json:"productId"`
	ProductKey      string     `json:"productKey"`
	ProductName     string     `json:"productName"`
	State           string     `json:"state"`
	Enabled         bool       `json:"enabled"`
	LastOnlineTime  *int64     `json:"lastOnlineTime"`
	LastOfflineTime *int64     `json:"lastOfflineTime"`
	Metadata        string     `json:"metadata"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	DeletedAt       *time.Time `json:"deletedAt,omitempty"`
}

type CreateDeviceResponse struct {
	Device Device `json:"device"`
}

type GetDeviceResponse struct {
	Device Device `json:"device"`
}

type UpdateDeviceResponse struct {
	Device Device `json:"device"`
}

type ActivateDeviceResponse struct {
	Device Device `json:"device"`
}

type SetDeviceEnabledResponse struct {
	Device Device `json:"device"`
}

type RestoreDeviceResponse struct {
	Device Device `json:"device"`
}

type ListDevicesResponse struct {
	Devices  []Device `json:"devices"`
	Total    int64    `json:"total"`
	Page     int      `json:"page"`
	PageSize int      `json:"pageSize"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

type TelemetryResponse struct {
	Telemetry string `json:"telemetry"`
}

type DeviceStatisticsResponse struct {
	TotalDevices     int64 `json:"totalDevices"`
	ActivatedDevices int64 `json:"activatedDevices"`
	OnlineDevices    int64 `json:"onlineDevices"`
	OfflineDevices   int64 `json:"offlineDevices"`
	InactiveDevices  int64 `json:"inactiveDevices"`
}

type MQTTParametersResponse struct {
	ClientID    string `json:"clientId"`
	Username    string `json:"username"`
	MQTTHostURL string `json:"mqttHostUrl"`
	Password    string `json:"password"`
	Port        int32  `json:"port"`
}

type Shadow struct {
	Desired   any            `json:"desired"`
	Reported  any            `json:"reported"`
	Delta     map[string]any `json:"delta"`
	Metadata  any            `json:"metadata"`
	Version   int64          `json:"version"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

type DeviceTag struct {
	ID        int64      `json:"id"`
	DeviceID  int64      `json:"deviceId"`
	Key       string     `json:"key"`
	Value     string     `json:"value"`
	Source    string     `json:"source"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

type ListDeviceTagsResponse struct {
	Tags []DeviceTag `json:"tags"`
}

type DeviceShadowHistory struct {
	ID        int64     `json:"id"`
	DeviceID  int64     `json:"deviceId"`
	Version   int64     `json:"version"`
	Source    string    `json:"source"`
	Desired   any       `json:"desired"`
	Reported  any       `json:"reported"`
	CreatedAt time.Time `json:"createdAt"`
}

type ListDeviceShadowHistoryResponse struct {
	History []DeviceShadowHistory `json:"history"`
}

type PushRecord struct {
	ID            int64     `json:"id"`
	DeviceID      int64     `json:"deviceId"`
	OperationType string    `json:"operationType"`
	OperationName string    `json:"operationName"`
	Payload       string    `json:"payload"`
	Status        string    `json:"status"`
	ErrorMessage  string    `json:"errorMessage"`
	CreatedBy     string    `json:"createdBy"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type SimulatePushResponse struct {
	PushRecordID string `json:"pushRecordId"`
	Timestamp    int64  `json:"timestamp"`
	Status       string `json:"status"`
	Message      string `json:"message"`
}

type ListPushRecordsResponse struct {
	Records  []PushRecord `json:"records"`
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
	PageSize int          `json:"pageSize"`
}

type GetPushRecordResponse struct {
	Record PushRecord `json:"record"`
}

type ClearPushRecordsResponse struct {
	DeletedCount int64 `json:"deletedCount"`
	Success      bool  `json:"success"`
}

type BatchUploadError struct {
	Row        int    `json:"row"`
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
	Error      string `json:"error"`
}

type BatchUploadDevicesResponse struct {
	SuccessCount int                `json:"successCount"`
	FailureCount int                `json:"failureCount"`
	Errors       []BatchUploadError `json:"errors"`
}

type MockKafkaResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
