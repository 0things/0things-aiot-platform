package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"

	"aiot-backend/internal/tenant"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var (
	ErrR2NotConfigured = errors.New("R2 storage is not configured")
	ErrR2Unavailable   = errors.New("R2 uploads are unavailable on this server platform")
)

type FileUploadResult struct {
	FileURL  string
	FileSize int64
	Checksum string
}

type FileServiceInterface interface {
	UploadOTA(ctx context.Context, filename string, data []byte) (*FileUploadResult, error)
}

type FileService struct {
	config *viper.Viper
}

func NewFileService(config *viper.Viper) *FileService {
	return &FileService{config: config}
}

func (s *FileService) UploadOTA(ctx context.Context, filename string, data []byte) (*FileUploadResult, error) {
	if !s.isR2Configured() {
		return nil, ErrR2NotConfigured
	}

	objectKey := s.otaObjectKey(ctx, filename)
	if err := s.uploadR2(objectKey, data); err != nil {
		return nil, err
	}

	digest := sha256.Sum256(data)
	return &FileUploadResult{
		FileURL:  strings.TrimRight(s.config.GetString("data.storage.r2.public_base_url"), "/") + "/" + objectKey,
		FileSize: int64(len(data)),
		Checksum: hex.EncodeToString(digest[:]),
	}, nil
}

func (s *FileService) isR2Configured() bool {
	if s.config == nil {
		return false
	}
	for _, key := range []string{
		"data.storage.r2.endpoint",
		"data.storage.r2.bucket",
		"data.storage.r2.access_key_id",
		"data.storage.r2.secret_access_key",
		"data.storage.r2.public_base_url",
	} {
		if s.config.GetString(key) == "" {
			return false
		}
	}
	return true
}

func (s *FileService) otaObjectKey(ctx context.Context, filename string) string {
	prefix := strings.Trim(s.config.GetString("data.storage.r2.prefix"), "/")
	if prefix == "" {
		prefix = "ota"
	}
	extension := strings.ToLower(filepath.Ext(filename))
	return path.Join(prefix, time.Now().UTC().Format("2006/01/02"), fmt.Sprintf("%d", tenant.GetTenantID(ctx)), uuid.NewString()+extension)
}
