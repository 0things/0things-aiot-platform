package handler

import (
	"errors"
	"io"
	"net/http"
	"strings"

	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/service"

	"github.com/gin-gonic/gin"
)

const maxOTAFileSize = 100 << 20

var allowedOTAFileExtensions = []string{
	".tar.gz",
	".tar.xz",
	".gzip",
	".pack",
	".bin",
	".dav",
	".tar",
	".gz",
	".zip",
	".apk",
}

type FileHandler struct {
	*Handler
	svc service.FileServiceInterface
}

func NewFileHandler(h *Handler, svc service.FileServiceInterface) *FileHandler {
	return &FileHandler{Handler: h, svc: svc}
}

// UploadOTAFile godoc
// @Summary Upload OTA update file
// @Description Uploads OTA update file.
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param file formData file true "OTA upgrade file (.bin, .dav, .tar, .gz, .zip, .gzip, .apk, .tar.gz, .tar.xz, .pack; max 100 MB)"
// @Success 200 {object} v1.ApiResponse[v1.UploadOTAFileResponse] "Successful response"
// @Router /files/ota [post]
func (h *FileHandler) UploadOTAFile(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxOTAFileSize+1)
	fileHeader, err := c.FormFile("file")
	if err != nil {
		v1.HandleError(c, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	if fileHeader.Size > maxOTAFileSize || !isAllowedOTAFile(fileHeader.Filename) {
		v1.HandleError(c, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		v1.HandleError(c, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, maxOTAFileSize+1))
	if err != nil || len(data) == 0 || len(data) > maxOTAFileSize {
		v1.HandleError(c, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	result, err := h.svc.UploadOTA(c.Request.Context(), fileHeader.Filename, data)
	if err != nil {
		if errors.Is(err, service.ErrR2NotConfigured) || errors.Is(err, service.ErrR2Unavailable) {
			v1.HandleError(c, http.StatusServiceUnavailable, v1.ErrStorageUnavailable, nil)
			return
		}
		v1.HandleError(c, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(c, v1.UploadOTAFileResponse{
		FileURL:  result.FileURL,
		FileSize: result.FileSize,
		Checksum: result.Checksum,
	})
}

func isAllowedOTAFile(filename string) bool {
	filename = strings.ToLower(filename)
	for _, extension := range allowedOTAFileExtensions {
		if strings.HasSuffix(filename, extension) {
			return true
		}
	}
	return false
}
