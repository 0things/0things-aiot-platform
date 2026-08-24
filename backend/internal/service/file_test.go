package service

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestOTAObjectKey_DatePrefixAndTenant(t *testing.T) {
	v := viper.New()
	v.Set("data.storage.r2.prefix", "ota")
	s := NewFileService(v)

	key := s.otaObjectKey(context.Background(), "Firmware.BIN")

	parts := strings.Split(key, "/")
	require.Equal(t, "ota", parts[0])
	require.Len(t, parts[1], 4, "year")
	require.Len(t, parts[2], 2, "month")
	require.Len(t, parts[3], 2, "day")
	require.Equal(t, "1", parts[4], "default tenant id")
	require.True(t, strings.HasSuffix(key, ".bin"), "extension should be lowercased")

	expectedDate := time.Now().UTC().Format("2006/01/02")
	require.Equal(t, expectedDate, strings.Join(parts[1:4], "/"))
}

func TestOTAObjectKey_DefaultPrefix(t *testing.T) {
	v := viper.New()
	s := NewFileService(v)

	key := s.otaObjectKey(context.Background(), "a.bin")
	require.True(t, strings.HasPrefix(key, "ota/"), "prefix defaults to ota, got %s", key)
}

func TestOTAObjectKey_CustomPrefix(t *testing.T) {
	v := viper.New()
	v.Set("data.storage.r2.prefix", "firmware")
	s := NewFileService(v)

	key := s.otaObjectKey(context.Background(), "a.bin")
	require.True(t, strings.HasPrefix(key, "firmware/"), "custom prefix should be used, got %s", key)
}

func TestOTAObjectKey_ExtensionIsLowercased(t *testing.T) {
	v := viper.New()
	v.Set("data.storage.r2.prefix", "ota")
	s := NewFileService(v)

	// filepath.Ext only keeps the final extension, lowercased.
	key := s.otaObjectKey(context.Background(), "package.TAR.GZ")
	require.True(t, strings.HasSuffix(key, ".gz"), "only the last extension is kept and lowercased, got %s", key)
}

func TestIsR2Configured(t *testing.T) {
	v := viper.New()
	s := NewFileService(v)
	require.False(t, s.isR2Configured(), "empty config should be unconfigured")

	for _, key := range []string{
		"data.storage.r2.endpoint",
		"data.storage.r2.bucket",
		"data.storage.r2.access_key_id",
		"data.storage.r2.secret_access_key",
		"data.storage.r2.public_base_url",
	} {
		v.Set(key, "x")
	}
	require.True(t, s.isR2Configured(), "all fields set should be configured")

	v.Set("data.storage.r2.secret_access_key", "")
	require.False(t, s.isR2Configured(), "missing one field should be unconfigured")
}

func TestUploadOTA_NotConfigured(t *testing.T) {
	v := viper.New()
	s := NewFileService(v)

	_, err := s.UploadOTA(context.Background(), "a.bin", []byte("payload"))
	require.ErrorIs(t, err, ErrR2NotConfigured)
}

func TestUploadOTA_Integration(t *testing.T) {
	if os.Getenv("R2_INTEGRATION_TEST") == "" {
		t.Skip("set R2_INTEGRATION_TEST=1 and TEST_CONF to run against real R2")
	}
	cfg := os.Getenv("TEST_CONF")
	if cfg == "" {
		t.Fatal("TEST_CONF must point to a config with R2 credentials")
	}
	v := viper.New()
	v.SetConfigFile(cfg)
	require.NoError(t, v.ReadInConfig())

	s := NewFileService(v)
	res, err := s.UploadOTA(context.Background(), "itest.bin", []byte("integration-payload"))
	require.NoError(t, err)
	require.Contains(t, res.FileURL, "https://")
	require.Equal(t, int64(len("integration-payload")), res.FileSize)
	require.NotEmpty(t, res.Checksum)
	t.Logf("uploaded to %s", res.FileURL)
}
