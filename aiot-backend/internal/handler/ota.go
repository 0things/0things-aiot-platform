package handler

import (
	otaV1 "0things-backend/api/ota/v1"
	"0things-backend/internal/model"
	"0things-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type OTAHandler struct {
	*Handler
	svc *service.OTAService
}

func NewOTAHandler(h *Handler, svc *service.OTAService) *OTAHandler {
	return &OTAHandler{Handler: h, svc: svc}
}

func otaPackageJSON(pkg model.OTAPackage) otaV1.OTAPackage {
	return otaV1.OTAPackage{
		ID: pkg.ID, PackageName: pkg.PackageName, Version: pkg.Version, ProductID: pkg.ProductID,
		PackageType: pkg.PackageType, Status: pkg.Status, UploadType: pkg.UploadType, FileURL: pkg.FileURL,
		FileSize: pkg.FileSize, Checksum: pkg.Checksum, Description: pkg.Description, ReleaseNotes: pkg.ReleaseNotes,
		Metadata: string(pkg.Metadata), CreatedAt: pkg.CreatedAt, UpdatedAt: pkg.UpdatedAt, ReleasedAt: pkg.ReleasedAt,
	}
}

func otaBatchJSON(batch model.UpgradeBatch) otaV1.UpgradeBatch {
	return otaV1.UpgradeBatch{
		BatchID: batch.BatchID, BatchName: batch.BatchName, BatchType: batch.BatchType,
		UpgradeStrategy: batch.UpgradeStrategy, Status: batch.Status, TargetDeviceCount: batch.TargetDeviceCount, CreatedAt: batch.CreatedAt,
	}
}

func otaDeploymentJSON(deployment model.DeviceDeployment) otaV1.DeviceDeployment {
	return otaV1.DeviceDeployment{
		DeviceID: deployment.DeviceID, DeviceKey: deployment.DeviceKey, DeviceName: deployment.DeviceName,
		ProductID: deployment.ProductID, ProductKey: deployment.ProductKey, CurrentVersion: deployment.CurrentVersion,
		UpgradeBatchID: deployment.UpgradeBatchID, Status: deployment.Status,
		LastStatusChangeTime: deployment.LastStatusChangeTime, CreatedAt: deployment.CreatedAt,
	}
}

func (h *OTAHandler) ListOTA(c *gin.Context) {
	pageNumber, pageSize := page(c, 20)
	packages, total, err := h.svc.List(c, pageNumber, pageSize)
	if err != nil {
		deviceError(c, err)
		return
	}
	items := make([]otaV1.OTAPackage, len(packages))
	for i, pkg := range packages {
		items[i] = otaPackageJSON(pkg)
	}
	c.JSON(200, otaV1.ListOTAPackagesResponse{OTAPackages: items, Total: total, Page: pageNumber, PageSize: pageSize})
}

func (h *OTAHandler) GetOTA(c *gin.Context) {
	packageID, err := id(c)
	if err != nil {
		deviceError(c, err)
		return
	}
	pkg, err := h.svc.Get(c, packageID)
	if err != nil {
		deviceError(c, err)
		return
	}
	c.JSON(200, otaV1.GetOTAPackageResponse{OTAPackage: otaPackageJSON(*pkg)})
}

func (h *OTAHandler) CreateOTA(c *gin.Context) {
	var req otaV1.CreateOTAPackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		deviceError(c, err)
		return
	}
	pkg := &model.OTAPackage{PackageName: req.PackageName, Version: req.Version, PackageType: req.PackageType, Status: req.Status, UploadType: req.UploadType, FileURL: req.FileURL, FileSize: req.FileSize, Checksum: req.Checksum, Description: req.Description, ReleaseNotes: req.ReleaseNotes, Metadata: req.Metadata}
	if pkg.Status == "" {
		pkg.Status = "draft"
	}
	if err := h.svc.Create(c, pkg, req.ProductKey); err != nil {
		deviceError(c, err)
		return
	}
	c.JSON(200, otaV1.CreateOTAPackageResponse{OTAPackage: otaPackageJSON(*pkg)})
}

func (h *OTAHandler) UpdateOTA(c *gin.Context) {
	packageID, err := id(c)
	if err != nil {
		deviceError(c, err)
		return
	}
	var req otaV1.OTAPackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		deviceError(c, err)
		return
	}
	pkg, err := h.svc.Get(c, packageID)
	if err != nil {
		deviceError(c, err)
		return
	}
	pkg.PackageName, pkg.Version, pkg.ProductID, pkg.PackageType, pkg.Status, pkg.UploadType = req.PackageName, req.Version, req.ProductID, req.PackageType, req.Status, req.UploadType
	pkg.FileURL, pkg.FileSize, pkg.Checksum, pkg.Description, pkg.ReleaseNotes, pkg.Metadata = req.FileURL, req.FileSize, req.Checksum, req.Description, req.ReleaseNotes, req.Metadata
	if err := h.svc.Update(c, pkg); err != nil {
		deviceError(c, err)
		return
	}
	c.JSON(200, otaV1.UpdateOTAPackageResponse{OTAPackage: otaPackageJSON(*pkg)})
}

func (h *OTAHandler) DeleteOTA(c *gin.Context) {
	packageID, err := id(c)
	if err == nil {
		err = h.svc.Delete(c, packageID)
	}
	if err != nil {
		deviceError(c, err)
		return
	}
	c.JSON(200, otaV1.SuccessResponse{Success: true})
}

func (h *OTAHandler) OTAStats(c *gin.Context) {
	stats, err := h.svc.Statistics(c, c.Param("id"))
	if err != nil {
		deviceError(c, err)
		return
	}
	c.JSON(200, otaV1.GetUpgradeStatisticsResponse{Statistics: otaV1.UpgradeStatistics{
		PackageID: stats.PackageID, TotalTargetDevices: stats.TotalTargetDevices,
		SuccessfulUpgrades: stats.SuccessfulUpgrades, FailedUpgrades: stats.FailedUpgrades,
		CancelledUpgrades: stats.CancelledUpgrades, PendingUpgrades: stats.PendingUpgrades,
		InProgressUpgrades: stats.InProgressUpgrades,
	}})
}

func (h *OTAHandler) OTABatches(c *gin.Context) {
	batches, err := h.svc.Batches(c, c.Param("id"))
	if err != nil {
		deviceError(c, err)
		return
	}
	items := make([]otaV1.UpgradeBatch, len(batches))
	for i, batch := range batches {
		items[i] = otaBatchJSON(batch)
	}
	c.JSON(200, otaV1.ListUpgradeBatchesResponse{Batches: items})
}

func (h *OTAHandler) OTADeployments(c *gin.Context) {
	pageNumber, pageSize := page(c, 100)
	deployments, total, err := h.svc.Deployments(c, c.Param("id"), pageNumber, pageSize, c.Query("status"))
	if err != nil {
		deviceError(c, err)
		return
	}
	items := make([]otaV1.DeviceDeployment, len(deployments))
	for i, deployment := range deployments {
		items[i] = otaDeploymentJSON(deployment)
	}
	c.JSON(200, otaV1.ListDeviceDeploymentsResponse{Deployments: items, Total: total, Page: pageNumber, PageSize: pageSize})
}
