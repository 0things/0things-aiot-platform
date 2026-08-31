package protocol

import "context"

// Codec 负责应用层载荷编解码，将原始字节流（JSON/Modbus/GB28181等）转换为通用的物模型属性 Map，反之亦然。
type Codec interface {
	Name() string
	Decode(ctx context.Context, payload []byte) (map[string]any, error)
	Encode(ctx context.Context, value map[string]any) ([]byte, error)
}
