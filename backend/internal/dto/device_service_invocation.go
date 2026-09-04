package dto

import "time"

type ListDeviceServiceInvocationsQuery struct {
	DeviceKey         string
	DeviceID          int64
	ServiceIdentifier string
	StartAt           *time.Time
	EndAt             *time.Time
	Page              int
	PageSize          int
}
