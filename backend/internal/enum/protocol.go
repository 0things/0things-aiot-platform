package enum

// TransportProtocol identifies how a device connects to the platform.
type TransportProtocol string

const (
	TransportMQTT    TransportProtocol = "mqtt"    // MQTT 长连接消息接入。
	TransportHTTP    TransportProtocol = "http"    // HTTP 请求式设备接入。
	TransportCoAP    TransportProtocol = "coap"    // CoAP 低功耗设备接入。
	TransportTCP     TransportProtocol = "tcp"     // TCP 长连接或透传接入。
	TransportUDP     TransportProtocol = "udp"     // UDP 无连接设备接入。
	TransportModbus  TransportProtocol = "modbus"  // Modbus 网关接入。
	TransportGB28181 TransportProtocol = "gb28181" // GB28181 视频设备接入。
)

// ApplicationProtocol identifies how the payload is encoded.
type ApplicationProtocol string

const (
	ApplicationJSON     ApplicationProtocol = "json"     // JSON/物模型消息编码。
	ApplicationProtobuf ApplicationProtocol = "protobuf" // Protobuf 消息编码。
	ApplicationModbus   ApplicationProtocol = "modbus"   // Modbus 数据编码。
	ApplicationGB28181  ApplicationProtocol = "gb28181"  // GB28181 信令编码。
	ApplicationJT808    ApplicationProtocol = "jt808"    // JT/T 808 车辆终端协议。
	ApplicationJT1078   ApplicationProtocol = "jt1078"   // JT/T 1078 车载视频协议。
	ApplicationCustom   ApplicationProtocol = "custom"   // 由协议插件定义的自定义编码。
)
