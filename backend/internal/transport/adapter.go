package transport

import (
	"context"

	"aiot-backend/internal/enum"
)

// Adapter 是线协议与设备业务之间的边界，不得直接更新设备、遥测或 OTA 数据表。
type Adapter interface {
	Name() string
	Transport() enum.TransportProtocol
	Start(context.Context, func(context.Context, DeviceMessage) error) error
	Send(context.Context, Command) error
	Stop(context.Context) error
}
