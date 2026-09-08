package event

import "time"

// OTAUpgradeCommand is published by backend and consumed by data-engine/transports.
type OTAUpgradeCommand struct {
	BatchID       string    `json:"batch_id"`
	PackageID     string    `json:"package_id"`
	ProductKey    string    `json:"product_key"`
	DeviceKey     string    `json:"device_key"`
	DeviceName    string    `json:"device_name,omitempty"`
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

// OTAUpgradeReport is published by transports when a device reports progress/status.
type OTAUpgradeReport struct {
	EventType       string          `json:"event_type,omitempty"`
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
