package protocol

import "context"

// Codec 负责应用层载荷编解码，不处理 MQTT、HTTP、CoAP 等连接细节。
type Codec interface {
	Name() string
	Decode(context.Context, []byte) (map[string]any, error)
	Encode(context.Context, map[string]any) ([]byte, error)
}
