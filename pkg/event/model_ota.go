package event

import "time"

// OTAUpgradeCommand is published by backend and consumed by data-engine/transports.
type OTAUpgradeCommand struct {
	BatchID       string    `json:"batch_id"`
	PackageID     string    `json:"package_id"`
	ProductKey    string    `json:"product_key"`
	DeviceKey     string    `json:"device_key"`
	DeviceName    string    `json:"device_name,omitempty"`
	Transport     string    `json:"transport,omitempty"`
	Module        string    `json:"module"`
	TargetVersion string    `json:"target_version"`
	DownloadURL   string    `json:"download_url"`
	FileSize      int64     `json:"file_size"`
	SHA256        string    `json:"sha256"`
	ExpiresAt     time.Time `json:"expires_at"`
}

// OTAReportError details an error during device OTA update.
type OTAReportError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// OTAEventType represents standard OTA reporting event types.
type OTAEventType string

const (
	OTAEventTypeProgress OTAEventType = "progress"
	OTAEventTypeInform   OTAEventType = "inform"
)

// OTAStatus represents device OTA execution status.
type OTAStatus string

const (
	OTAStatusPending    OTAStatus = "pending"
	OTAStatusSent       OTAStatus = "sent"
	OTAStatusInProgress OTAStatus = "in_progress"
	OTAStatusSuccess    OTAStatus = "success"
	OTAStatusFailed     OTAStatus = "failed"
)

// OTAUpgradeReport is published by transports when a device reports progress/status.
type OTAUpgradeReport struct {
	EventType       OTAEventType    `json:"event_type,omitempty"`
	BatchID         string          `json:"batch_id"`
	ProductKey      string          `json:"product_key,omitempty"`
	DeviceKey       string          `json:"device_key"`
	Module          string          `json:"module,omitempty"`
	Status          string          `json:"status,omitempty"`
	Progress        *int32          `json:"progress,omitempty"`
	Stage           string          `json:"stage,omitempty"`
	ReportedVersion string          `json:"reported_version,omitempty"`
	Error           *OTAReportError `json:"error,omitempty"`
	ReportedAt      time.Time       `json:"reported_at"`
}
