package event

import "time"

// OTAUpgradeCommand is published by backend and consumed by data-engine/transports.
type OTAUpgradeCommand struct {
	BatchID       string    `json:"batch_id" validate:"required"`
	PackageID     string    `json:"package_id" validate:"required"`
	ProductKey    string    `json:"product_key,omitempty"`
	DeviceKey     string    `json:"device_key" validate:"required"`
	DeviceName    string    `json:"device_name,omitempty"`
	Transport     string    `json:"transport,omitempty"`
	Module        string    `json:"module,omitempty"`
	TargetVersion string    `json:"target_version" validate:"required"`
	DownloadURL   string    `json:"download_url" validate:"required"`
	FileSize      int64     `json:"file_size" validate:"required,gt=0"`
	SHA256        string    `json:"sha256" validate:"required"`
	ExpiresAt     time.Time `json:"expires_at"`
}

// OTAStatus represents device OTA execution status.
type OTAStatus string

const (
	OTAStatusPending    OTAStatus = "pending"
	OTAStatusSent       OTAStatus = "sent"
	OTAStatusInProgress OTAStatus = "in_progress"
	OTAStatusSuccess    OTAStatus = "success"
	OTAStatusFailed     OTAStatus = "failed"
	OTAStatusTimeout    OTAStatus = "timeout"
)

// OTAProgressPayload represents the inner payload for device OTA upgrade progress report (/ota/device/progress/+/+).
type OTAProgressPayload struct {
	BatchID      string `json:"batch_id" validate:"required"`
	Progress     int32  `json:"step" validate:"gte=0,lte=100"`
	ErrorMessage string `json:"desc,omitempty"`
	Timestamp    int64  `json:"timestamp,omitempty"`
}

// OTAInformPayload represents the inner payload for device firmware version inform report (/ota/device/inform/+/+).
type OTAInformPayload struct {
	ID        string `json:"id,omitempty"`
	Module    string `json:"module,omitempty"`
	Version   string `json:"version" validate:"required"`
	Timestamp int64  `json:"timestamp,omitempty"`
}
