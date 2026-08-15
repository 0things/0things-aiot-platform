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
}//@name DeviceCreateDeviceRequest

type UpdateDeviceRequest struct {
	Name     string          `json:"name"`
	State    string          `json:"state"`
	Metadata json.RawMessage `json:"metadata"`
}//@name DeviceUpdateDeviceRequest

type SetDeviceEnabledRequest struct {
	Enabled bool `json:"enabled"`
}//@name DeviceSetDeviceEnabledRequest

type SetDeviceTagsRequest struct {
	Tags map[string]string `json:"tags" binding:"required"`
}//@name DeviceSetDeviceTagsRequest

type DeleteDeviceTagsRequest struct {
	Keys []string `json:"keys" binding:"required"`
}//@name DeviceDeleteDeviceTagsRequest

type UpdateDesiredShadowRequest struct {
	Desired map[string]any `json:"desired" binding:"required"`
	Version int64          `json:"version"`
}//@name DeviceUpdateDesiredShadowRequest

type UpdateReportedShadowRequest struct {
	Reported map[string]any `json:"reported" binding:"required"`
	Version  int64          `json:"version"`
}//@name DeviceUpdateReportedShadowRequest

type ClearDesiredShadowRequest struct {
	Version int64 `json:"version"`
}//@name DeviceClearDesiredShadowRequest

type SimulatePushRequest struct {
	Payload string `json:"payload"`
}//@name DeviceSimulatePushRequest

type MockKafkaRequest struct {
	Topic string `json:"topic" binding:"required"`
	Data  string `json:"data"`
}//@name DeviceMockKafkaRequest

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
}//@name Device

type CreateDeviceResponse struct {
	Device Device `json:"device"`
}//@name DeviceCreateDeviceResponse

type GetDeviceResponse struct {
	Device Device `json:"device"`
}//@name DeviceGetDeviceResponse

type UpdateDeviceResponse struct {
	Device Device `json:"device"`
}//@name DeviceUpdateDeviceResponse

type ActivateDeviceResponse struct {
	Device Device `json:"device"`
}//@name DeviceActivateDeviceResponse

type SetDeviceEnabledResponse struct {
	Device Device `json:"device"`
}//@name DeviceSetDeviceEnabledResponse

type RestoreDeviceResponse struct {
	Device Device `json:"device"`
}//@name DeviceRestoreDeviceResponse

type ListDevicesResponse struct {
	Devices  []Device `json:"devices"`
	Total    int64    `json:"total"`
	Page     int      `json:"page"`
	PageSize int      `json:"pageSize"`
}//@name DeviceListDevicesResponse

type SuccessResponse struct {
	Success bool `json:"success"`
}//@name DeviceSuccessResponse

type TelemetryResponse struct {
	Telemetry string `json:"telemetry"`
}//@name DeviceTelemetryResponse

type DeviceStatisticsResponse struct {
	TotalDevices     int64 `json:"totalDevices"`
	ActivatedDevices int64 `json:"activatedDevices"`
	OnlineDevices    int64 `json:"onlineDevices"`
	OfflineDevices   int64 `json:"offlineDevices"`
	InactiveDevices  int64 `json:"inactiveDevices"`
}//@name DeviceStatisticsResponse

type MQTTParametersResponse struct {
	ClientID    string `json:"clientId"`
	Username    string `json:"username"`
	MQTTHostURL string `json:"mqttHostUrl"`
	Password    string `json:"password"`
	Port        int32  `json:"port"`
}//@name DeviceMQTTParametersResponse

type Shadow struct {
	Desired   any            `json:"desired"`
	Reported  any            `json:"reported"`
	Delta     map[string]any `json:"delta"`
	Metadata  any            `json:"metadata"`
	Version   int64          `json:"version"`
	UpdatedAt time.Time      `json:"updatedAt"`
}//@name DeviceShadow

type DeviceTag struct {
	ID        int64      `json:"id"`
	DeviceID  int64      `json:"deviceId"`
	Key       string     `json:"key"`
	Value     string     `json:"value"`
	Source    string     `json:"source"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}//@name DeviceTag

type ListDeviceTagsResponse struct {
	Tags []DeviceTag `json:"tags"`
}//@name DeviceListDeviceTagsResponse

type DeviceShadowHistory struct {
	ID        int64     `json:"id"`
	DeviceID  int64     `json:"deviceId"`
	Version   int64     `json:"version"`
	Source    string    `json:"source"`
	Desired   any       `json:"desired"`
	Reported  any       `json:"reported"`
	CreatedAt time.Time `json:"createdAt"`
}//@name DeviceShadowHistory

type ListDeviceShadowHistoryResponse struct {
	History []DeviceShadowHistory `json:"history"`
}//@name DeviceListDeviceShadowHistoryResponse

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
}//@name DevicePushRecord

type SimulatePushResponse struct {
	PushRecordID string `json:"pushRecordId"`
	Timestamp    int64  `json:"timestamp"`
	Status       string `json:"status"`
	Message      string `json:"message"`
}//@name DeviceSimulatePushResponse

type ListPushRecordsResponse struct {
	Records  []PushRecord `json:"records"`
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
	PageSize int          `json:"pageSize"`
}//@name DeviceListPushRecordsResponse

type GetPushRecordResponse struct {
	Record PushRecord `json:"record"`
}//@name DeviceGetPushRecordResponse

type ClearPushRecordsResponse struct {
	DeletedCount int64 `json:"deletedCount"`
	Success      bool  `json:"success"`
}//@name DeviceClearPushRecordsResponse

type BatchUploadError struct {
	Row        int    `json:"row"`
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
	Error      string `json:"error"`
}//@name DeviceBatchUploadError

type BatchUploadDevicesResponse struct {
	SuccessCount int                `json:"successCount"`
	FailureCount int                `json:"failureCount"`
	Errors       []BatchUploadError `json:"errors"`
}//@name DeviceBatchUploadDevicesResponse

type MockKafkaResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}//@name DeviceMockKafkaResponse
