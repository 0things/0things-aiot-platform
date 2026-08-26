package handler

import (
	"errors"
	"net/http"

	otaV1 "aiot-backend/api/ota/v1"
	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
	"aiot-backend/internal/service"

	"github.com/dromara/carbon/v2"
	"github.com/gin-gonic/gin"
)

func otaErrorStatus(err error) int {
	if errors.Is(err, repository.ErrNotFound) {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

type OTAHandler struct {
	*Handler
	svc service.OTAServiceInterface
}

func NewOTAHandler(h *Handler, svc service.OTAServiceInterface) *OTAHandler {
	return &OTAHandler{Handler: h, svc: svc}
}

func otaPackageJSON(pkg model.OTAPackage) otaV1.OTAPackage {
	createdAt := ""
	updatedAt := ""
	var releasedAt *string
	if !pkg.CreatedAt.IsZero() {
		createdAt = carbon.CreateFromStdTime(pkg.CreatedAt).ToDateTimeString()
	}
	if !pkg.UpdatedAt.IsZero() {
		updatedAt = carbon.CreateFromStdTime(pkg.UpdatedAt).ToDateTimeString()
	}
	if pkg.ReleasedAt != nil && !pkg.ReleasedAt.IsZero() {
		formatted := carbon.CreateFromStdTime(*pkg.ReleasedAt).ToDateTimeString()
		releasedAt = &formatted
	}
	return otaV1.OTAPackage{
		ID: pkg.ID, UUID: pkg.UUID, PackageName: pkg.PackageName, Version: pkg.Version,
		ProductID: pkg.ProductID, ProductKey: pkg.ProductKey, ProductName: pkg.ProductName,
		PackageType: pkg.PackageType, Status: pkg.Status, UploadType: pkg.UploadType, FileURL: pkg.FileURL,
		FileSize: pkg.FileSize, Checksum: pkg.Checksum, Description: pkg.Description, ReleaseNotes: pkg.ReleaseNotes,
		CreatedAt: createdAt, UpdatedAt: updatedAt, ReleasedAt: releasedAt,
	}
}

func otaBatchJSON(batch model.UpgradeBatch) otaV1.UpgradeBatch {
	createdAt := ""
	if !batch.CreatedAt.IsZero() {
		createdAt = carbon.CreateFromStdTime(batch.CreatedAt).ToDateTimeString()
	}
	return otaV1.UpgradeBatch{
		BatchID:         batch.BatchID,
		UpgradeStrategy: batch.UpgradeStrategy, Status: batch.Status, TargetDeviceCount: batch.TargetDeviceCount, CreatedAt: createdAt,
	}
}

func otaDeploymentJSON(deployment model.DeviceDeployment) otaV1.DeviceDeployment {
	lastStatusChangeTime := ""
	if deployment.LastStatusChangeTime != 0 {
		lastStatusChangeTime = carbon.CreateFromTimestamp(
			deployment.LastStatusChangeTime,
		).ToDateTimeString()
	}
	return otaV1.DeviceDeployment{
		DeviceID: deployment.DeviceID, DeviceKey: deployment.DeviceKey, DeviceName: deployment.DeviceName,
		ProductID: deployment.ProductID, ProductKey: deployment.ProductKey, CurrentVersion: deployment.CurrentVersion,
		UpgradeBatchID: deployment.UpgradeBatchID, Status: deployment.Status,
		LastStatusChangeTime: lastStatusChangeTime,
		CreatedAt:            carbon.CreateFromStdTime(deployment.CreatedAt).ToDateTimeString(),
	}
}

// ListOTA godoc
// @Summary 获取 OTA 升级包列表
// @Schemes
// @Description 分页获取 OTA 升级包列表
// @Tags OTA 模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} v1.ApiResponse[otaV1.ListOTAPackagesResponse]
// @Router /ota-packages [get]
func (h *OTAHandler) ListOTA(c *gin.Context) {
	pageNumber, pageSize := page(c, 20)
	packages, total, err := h.svc.List(c, pageNumber, pageSize)
	if err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	items := make([]otaV1.OTAPackage, len(packages))
	for i, pkg := range packages {
		items[i] = otaPackageJSON(pkg)
	}
	v1.HandleSuccess(c, otaV1.ListOTAPackagesResponse{OTAPackages: items, Total: total, Page: pageNumber, PageSize: pageSize})
}

// GetOTA godoc
// @Summary 获取 OTA 升级包详情
// @Schemes
// @Description 通过升级包 ID 获取 OTA 升级包详情
// @Tags OTA 模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param uuid path string true "升级包 UUID"
// @Success 200 {object} v1.ApiResponse[otaV1.OTAPackage]
// @Router /ota-packages/{uuid} [get]
func (h *OTAHandler) GetOTA(c *gin.Context) {
	uuid := c.Param("uuid")
	pkg, err := h.svc.Get(c, uuid)
	if err != nil {
		v1.HandleError(c, otaErrorStatus(err), err, nil)
		return
	}
	v1.HandleSuccess(c, otaPackageJSON(*pkg))
}

// CreateOTA godoc
// @Summary 创建 OTA 升级包
// @Schemes
// @Description 创建一个新的 OTA 升级包
// @Tags OTA 模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body otaV1.CreateOTAPackageRequest true "params"
// @Success 200 {object} v1.ApiResponse[otaV1.OTAPackage]
// @Router /ota-packages [post]
func (h *OTAHandler) CreateOTA(c *gin.Context) {
	var req otaV1.CreateOTAPackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	pkg := &model.OTAPackage{PackageName: req.PackageName, Version: req.Version, PackageType: req.PackageType, Status: req.Status, UploadType: req.UploadType, FileURL: req.FileURL, FileSize: req.FileSize, Checksum: req.Checksum, Description: req.Description, ReleaseNotes: req.ReleaseNotes}
	if pkg.Status == "" {
		pkg.Status = "draft"
	}
	if err := h.svc.Create(c, pkg, req.ProductKey); err != nil {
		v1.HandleError(c, otaErrorStatus(err), err, nil)
		return
	}
	v1.HandleSuccess(c, otaPackageJSON(*pkg))
}

// UpdateOTA godoc
// @Summary 更新 OTA 升级包
// @Schemes
// @Description 通过升级包 ID 更新 OTA 升级包
// @Tags OTA 模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param uuid path string true "升级包 UUID"
// @Param request body otaV1.OTAPackageRequest true "params"
// @Success 200 {object} v1.ApiResponse[otaV1.OTAPackage]
// @Router /ota-packages/{uuid} [put]
func (h *OTAHandler) UpdateOTA(c *gin.Context) {
	uuid := c.Param("uuid")
	var req otaV1.OTAPackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	pkg, err := h.svc.Get(c, uuid)
	if err != nil {
		v1.HandleError(c, otaErrorStatus(err), err, nil)
		return
	}
	if req.PackageName != "" {
		pkg.PackageName = req.PackageName
	}
	if req.Version != "" {
		pkg.Version = req.Version
	}
	if req.ProductID != 0 {
		pkg.ProductID = req.ProductID
	}
	if req.PackageType != "" {
		pkg.PackageType = req.PackageType
	}
	if req.Status != "" {
		pkg.Status = req.Status
	}
	if req.UploadType != "" {
		pkg.UploadType = req.UploadType
	}
	if req.FileURL != "" {
		pkg.FileURL = req.FileURL
	}
	if req.FileSize != 0 {
		pkg.FileSize = req.FileSize
	}
	if req.Checksum != "" {
		pkg.Checksum = req.Checksum
	}
	if req.Description != "" {
		pkg.Description = req.Description
	}
	if req.ReleaseNotes != "" {
		pkg.ReleaseNotes = req.ReleaseNotes
	}
	if err := h.svc.Update(c, pkg); err != nil {
		v1.HandleError(c, otaErrorStatus(err), err, nil)
		return
	}
	v1.HandleSuccess(c, otaPackageJSON(*pkg))
}

// DeleteOTA godoc
// @Summary 删除 OTA 升级包
// @Schemes
// @Description 通过升级包 ID 删除 OTA 升级包
// @Tags OTA 模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param uuid path string true "升级包 UUID"
// @Success 200 {object} v1.ApiResponse[otaV1.SuccessResponse]
// @Router /ota-packages/{uuid} [delete]
func (h *OTAHandler) DeleteOTA(c *gin.Context) {
	uuid := c.Param("uuid")
	if err := h.svc.Delete(c, uuid); err != nil {
		v1.HandleError(c, otaErrorStatus(err), err, nil)
		return
	}
	v1.HandleSuccess(c, otaV1.SuccessResponse{Success: true})
}

// BatchUpgradeOTA godoc
// @Summary 批量升级 OTA 升级包
// @Schemes
// @Description 为指定升级包创建一个静态升级批次，并将所选设备加入该批次的升级（状态 pending），升级包状态置为 deploying
// @Tags OTA 模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param uuid path string true "升级包 UUID"
// @Param request body otaV1.BatchUpgradeRequest true "目标设备 deviceKey 列表"
// @Success 200 {object} v1.ApiResponse[otaV1.UpgradeBatch]
// @Router /ota-packages/{uuid}/batch-upgrade [post]
func (h *OTAHandler) BatchUpgradeOTA(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		v1.HandleError(c, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	var req otaV1.BatchUpgradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	if len(req.DeviceKeys) == 0 {
		v1.HandleError(c, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	batch, err := h.svc.BatchUpgrade(c, uuid, req.DeviceKeys)
	if err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(c, otaBatchJSON(*batch))
}

// ReportOTAStatus godoc
// @Summary 上报设备 OTA 升级结果
// @Schemes
// @Description 设备上报对指定升级包的升级结果（in_progress/success/failed），并重新聚合升级包状态
// @Tags OTA 模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param uuid path string true "升级包 UUID"
// @Param request body otaV1.ReportOTAStatusRequest true "params"
// @Success 200 {object} v1.ApiResponse[otaV1.SuccessResponse]
// @Router /ota-packages/{uuid}/report [post]
func (h *OTAHandler) ReportOTAStatus(c *gin.Context) {
	uuid := c.Param("uuid")
	var req otaV1.ReportOTAStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	if err := h.svc.ReportStatus(c, uuid, req.DeviceKey, req.Status); err != nil {
		v1.HandleError(c, otaErrorStatus(err), err, nil)
		return
	}
	v1.HandleSuccess(c, otaV1.SuccessResponse{Success: true})
}

// OTAStats godoc
// @Summary 获取 OTA 升级统计
// @Schemes
// @Description 获取指定 OTA 升级包的升级统计数据
// @Tags OTA 模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param uuid path string true "升级包 UUID"
// @Success 200 {object} v1.ApiResponse[otaV1.GetUpgradeStatisticsResponse]
// @Router /ota-packages/{uuid}/upgrade-statistics [get]
func (h *OTAHandler) OTAStats(c *gin.Context) {
	stats, err := h.svc.Statistics(c, c.Param("uuid"))
	if err != nil {
		v1.HandleError(c, otaErrorStatus(err), err, nil)
		return
	}
	v1.HandleSuccess(c, otaV1.GetUpgradeStatisticsResponse{Statistics: otaV1.UpgradeStatistics{
		PackageID: stats.PackageID, TotalTargetDevices: stats.TotalTargetDevices,
		SuccessfulUpgrades: stats.SuccessfulUpgrades, FailedUpgrades: stats.FailedUpgrades,
		CancelledUpgrades: stats.CancelledUpgrades, PendingUpgrades: stats.PendingUpgrades,
		InProgressUpgrades: stats.InProgressUpgrades,
	}})
}

// OTABatches godoc
// @Summary 获取 OTA 升级批次列表
// @Schemes
// @Description 获取指定 OTA 升级包下的升级批次
// @Tags OTA 模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param uuid path string true "升级包 UUID"
// @Success 200 {object} v1.ApiResponse[otaV1.ListUpgradeBatchesResponse]
// @Router /ota-packages/{uuid}/batches [get]
func (h *OTAHandler) OTABatches(c *gin.Context) {
	batches, err := h.svc.Batches(c, c.Param("uuid"))
	if err != nil {
		v1.HandleError(c, otaErrorStatus(err), err, nil)
		return
	}
	items := make([]otaV1.UpgradeBatch, len(batches))
	for i, batch := range batches {
		items[i] = otaBatchJSON(batch)
	}
	v1.HandleSuccess(c, otaV1.ListUpgradeBatchesResponse{Items: items})
}

// OTADeployments godoc
// @Summary 获取 OTA 设备部署列表
// @Schemes
// @Description 分页获取指定升级包下的设备部署记录
// @Tags OTA 模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param uuid path string true "升级包 UUID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param status query string false "部署状态"
// @Success 200 {object} v1.ApiResponse[otaV1.ListDeviceDeploymentsResponse]
// @Router /ota-packages/{uuid}/device-deployments [get]
func (h *OTAHandler) OTADeployments(c *gin.Context) {
	pageNumber, pageSize := page(c, 100)
	deployments, total, err := h.svc.Deployments(c, c.Param("uuid"), pageNumber, pageSize, c.Query("status"))
	if err != nil {
		v1.HandleError(c, otaErrorStatus(err), err, nil)
		return
	}
	items := make([]otaV1.DeviceDeployment, len(deployments))
	for i, deployment := range deployments {
		items[i] = otaDeploymentJSON(deployment)
	}
	v1.HandleSuccess(c, otaV1.ListDeviceDeploymentsResponse{Items: items, Total: total, Page: pageNumber, PageSize: pageSize})
}
