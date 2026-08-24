package handler

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"aiot-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type fakeFileService struct {
	filename string
	data     []byte
}

func (s *fakeFileService) UploadOTA(_ context.Context, filename string, data []byte) (*service.FileUploadResult, error) {
	s.filename = filename
	s.data = data
	return &service.FileUploadResult{
		FileURL:  "https://files.example.com/ota/1/test.bin",
		FileSize: int64(len(data)),
		Checksum: "checksum",
	}, nil
}

func TestFileHandlerUploadOTAFile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	storage := &fakeFileService{}
	h := NewFileHandler(&Handler{}, storage)
	router := gin.New()
	router.POST("/files/ota", h.UploadOTAFile)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "firmware.bin")
	require.NoError(t, err)
	_, err = part.Write([]byte{0x01, 0x02, 0x03})
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/files/ota", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "firmware.bin", storage.filename)
	require.Equal(t, []byte{0x01, 0x02, 0x03}, storage.data)
}

func TestFileHandlerUploadOTAFileRejectsUnsupportedExtension(t *testing.T) {
	gin.SetMode(gin.TestMode)
	storage := &fakeFileService{}
	h := NewFileHandler(&Handler{}, storage)
	router := gin.New()
	router.POST("/files/ota", h.UploadOTAFile)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "firmware.txt")
	require.NoError(t, err)
	_, err = part.Write([]byte("not firmware"))
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/files/ota", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Empty(t, storage.filename)
}

func TestIsAllowedOTAFile(t *testing.T) {
	for _, filename := range []string{
		"firmware.bin",
		"release.DAV",
		"package.tar",
		"package.gz",
		"package.zip",
		"package.gzip",
		"app.apk",
		"package.tar.gz",
		"package.tar.xz",
		"release.pack",
	} {
		require.Truef(t, isAllowedOTAFile(filename), "%s should be allowed", filename)
	}
	require.False(t, isAllowedOTAFile("package.xz"))
	require.False(t, isAllowedOTAFile("package.hex"))
}
