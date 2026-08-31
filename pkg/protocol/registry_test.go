package protocol_test

import (
	"context"
	"testing"

	"0things/pkg/protocol"
)

func TestRegistry(t *testing.T) {
	r := protocol.DefaultRegistry()

	// 验证默认注册的常用协议编解码器
	codecs := []string{"json", "modbus", "gb28181"}
	for _, name := range codecs {
		c, ok := r.Get(name)
		if !ok || c == nil {
			t.Fatalf("expected codec %s registered, but got not found", name)
		}
		if c.Name() != name {
			t.Fatalf("expected name %s, got %s", name, c.Name())
		}
	}

	// 验证 JSON 解码
	jsonCodec, _ := r.Get("json")
	res, err := jsonCodec.Decode(context.Background(), []byte(`{"temperature": 28.5}`))
	if err != nil {
		t.Fatalf("json decode failed: %v", err)
	}
	if res["temperature"] != 28.5 {
		t.Fatalf("unexpected value: %v", res)
	}
}
