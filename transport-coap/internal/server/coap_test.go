package server

import (
	"testing"

	"transport-coap/pkg/log"

	"github.com/spf13/viper"
)

func TestNewCoAPServer(t *testing.T) {
	v := viper.New()
	v.Set("coap.addr", ":5683")
	v.Set("log.mode", "console")
	v.Set("log.log_level", "error")
	logger := log.NewLog(v)

	srv := NewCoAPServer(v, logger)
	if srv == nil {
		t.Fatal("expected non-nil CoAPServer")
	}
	if srv.addr != ":5683" {
		t.Errorf("expected addr :5683, got %s", srv.addr)
	}
}
