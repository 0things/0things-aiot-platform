package jsoncodec

import (
	"context"
	"encoding/json"
)

// Codec 使用标准 JSON 作为设备应用层协议。
type Codec struct{}

func New() *Codec           { return &Codec{} }
func (*Codec) Name() string { return "json" }

func (*Codec) Decode(_ context.Context, payload []byte) (map[string]any, error) {
	var value map[string]any
	if err := json.Unmarshal(payload, &value); err != nil {
		return nil, err
	}
	return value, nil
}

func (*Codec) Encode(_ context.Context, value map[string]any) ([]byte, error) {
	return json.Marshal(value)
}
