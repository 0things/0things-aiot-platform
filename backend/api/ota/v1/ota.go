// Package otav1 owns OTA package HTTP contracts and mirrors device-service.
package otav1

import (
	"time"
)

type OTAPackageRequest struct {
	PackageName  string `json:"packageName" binding:"required"`
	Version      string `json:"version"`
	ProductID    int64  `json:"productId"`
	PackageType  string `json:"packageType"`
	Status       string `json:"status"`
	UploadType   string `json:"uploadType"`
	FileURL      string `json:"fileUrl"`
	FileSize     int64  `json:"fileSize"`
	Checksum     string `json:"checksum"`
	Description  string `json:"description"`
	ReleaseNotes string `json:"releaseNotes"`
}//@name OtaOTAPackageRequest

type CreateOTAPackageRequest struct {
	PackageName  string `json:"packageName" binding:"required"`
	Version      string `json:"version"`
	ProductKey   string `json:"product_key" binding:"required"`
	PackageType  string `json:"packageType"`
	Status       string `json:"status"`
	UploadType   string `json:"uploadType"`
	FileURL      string `json:"fileUrl"`
	FileSize     int64  `json:"fileSize"`
	Checksum     string `json:"checksum"`
	Description  string `json:"description"`
	ReleaseNotes string `json:"releaseNotes"`
}//@name OtaCreateOTAPackageRequest

type DeployOTAPackageRequest struct {
	DeviceKeys []string `json:"deviceKeys" binding:"required"`
}//@name OtaDeployOTAPackageRequest

type BatchUpgradeRequest struct {
	DeviceKeys []string `json:"deviceKeys" binding:"required"`
}//@name OtaBatchUpgradeRequest

type ReportOTAStatusRequest struct {
	DeviceKey string `json:"deviceKey" binding:"required"`
	Status    string `json:"status" binding:"required"`
}//@name OtaReportOTAStatusRequest

type SuccessResponse struct {
	Success bool `json:"success"`
}//@name OtaSuccessResponse

type OTAPackage struct {
	ID           int64      `json:"id"`
	UUID         string     `json:"uuid"`
	PackageName  string     `json:"packageName"`
	Version      string     `json:"version"`
	ProductID    int64      `json:"productId"`
	ProductKey   string     `json:"productKey,omitempty"`
	ProductName  string     `json:"productName,omitempty"`
	PackageType  string     `json:"packageType"`
	Status       string     `json:"status"`
	UploadType   string     `json:"uploadType"`
	FileURL      string     `json:"fileUrl"`
	FileSize     int64      `json:"fileSize"`
	Checksum     string     `json:"checksum"`
	Description  string     `json:"description"`
	ReleaseNotes string     `json:"releaseNotes"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	ReleasedAt   *time.Time `json:"releasedAt,omitempty"`
}//@name OtaOTAPackage

type UpgradeStatistics struct {
	PackageID          string `json:"packageId"`
	TotalTargetDevices int64  `json:"totalTargetDevices"`
	SuccessfulUpgrades int64  `json:"successfulUpgrades"`
	FailedUpgrades     int64  `json:"failedUpgrades"`
	CancelledUpgrades  int64  `json:"cancelledUpgrades"`
	PendingUpgrades    int64  `json:"pendingUpgrades"`
	InProgressUpgrades int64  `json:"inProgressUpgrades"`
}//@name OtaUpgradeStatistics

type UpgradeBatch struct {
	BatchID           string    `json:"batchId"`
	UpgradeStrategy   string    `json:"upgradeStrategy"`
	Status            string    `json:"status"`
	TargetDeviceCount int32     `json:"targetDeviceCount"`
	CreatedAt         time.Time `json:"createdAt"`
}//@name OtaUpgradeBatch

type DeviceDeployment struct {
	DeviceID             int64     `json:"deviceId"`
	DeviceKey            string    `json:"deviceKey"`
	DeviceName           string    `json:"deviceName"`
	ProductID            int64     `json:"productId"`
	ProductKey           string    `json:"productKey"`
	CurrentVersion       string    `json:"currentVersion"`
	UpgradeBatchID       string    `json:"upgradeBatchId"`
	Status               string    `json:"status"`
	LastStatusChangeTime int64     `json:"lastStatusChangeTime"`
	CreatedAt            time.Time `json:"createdAt"`
}//@name OtaDeviceDeployment

type ListOTAPackagesResponse struct {
	OTAPackages []OTAPackage `json:"otaPackages"`
	Total       int64        `json:"total"`
	Page        int          `json:"page"`
	PageSize    int          `json:"pageSize"`
}//@name OtaListOTAPackagesResponse

type GetOTAPackageResponse struct {
	OTAPackage OTAPackage `json:"otaPackage"`
}//@name OtaGetOTAPackageResponse
type CreateOTAPackageResponse struct {
	OTAPackage OTAPackage `json:"otaPackage"`
}//@name OtaCreateOTAPackageResponse
type UpdateOTAPackageResponse struct {
	OTAPackage OTAPackage `json:"otaPackage"`
}//@name OtaUpdateOTAPackageResponse
type GetUpgradeStatisticsResponse struct {
	Statistics UpgradeStatistics `json:"statistics"`
}//@name OtaGetUpgradeStatisticsResponse
type ListUpgradeBatchesResponse struct {
	Items []UpgradeBatch `json:"items"`
}//@name OtaListUpgradeBatchesResponse

type ListDeviceDeploymentsResponse struct {
	Items    []DeviceDeployment `json:"items"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"pageSize"`
}//@name OtaListDeviceDeploymentsResponse
