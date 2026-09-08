package model

import "time"

// OTAUpgradeCommand is created by backend and consumed only by data-engine,
// which owns the final MQTT dispatch.
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

// OTAUpgradeReport is forwarded by mqtt-transport to ota.upgrade.report.v1.
type OTAUpgradeReport struct {
	EventType       string          `json:"event_type"`
	BatchID         string          `json:"batch_id"`
	ProductKey      string          `json:"product_key"`
	DeviceKey       string          `json:"device_key"`
	Module          string          `json:"module"`
	Progress        *int32          `json:"progress,omitempty"`
	Stage           string          `json:"stage,omitempty"`
	ReportedVersion string          `json:"reported_version,omitempty"`
	Error           *OTAReportError `json:"error,omitempty"`
	ReportedAt      time.Time       `json:"reported_at"`
}

type OTAReportError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
