package eventv1

import "time"

type DeviceEvent struct {
	ID          int64     `json:"id"`
	DeviceKey   string    `json:"deviceKey"`
	DeviceName  string    `json:"deviceName"`
	EventType   string    `json:"eventType"`
	EventAt     time.Time `json:"eventAt"`
	Data        string    `json:"data"`
}//@name DeviceEvent

type ListDeviceEventsResponse struct {
	Events   []DeviceEvent `json:"events"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"pageSize"`
}//@name DeviceEventListDeviceEventsResponse
