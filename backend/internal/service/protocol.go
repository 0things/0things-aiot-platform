package service

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	protocolV1 "aiot-backend/api/v1"
	"aiot-backend/internal/enum"
	"aiot-backend/internal/repository"
	"github.com/spf13/viper"
)

type ProtocolServiceInterface interface {
	ListDeviceEndpoints(context.Context, string) (*protocolV1.DeviceEndpoints, error)
}

type ProtocolService struct {
	repo   *repository.ProtocolRepository
	config *viper.Viper
}

func NewProtocolService(repo *repository.ProtocolRepository, configs ...*viper.Viper) *ProtocolService {
	return &ProtocolService{repo: repo, config: func() *viper.Viper {
		if len(configs) > 0 {
			return configs[0]
		}
		return viper.New()
	}()}
}

// ListDeviceEndpoints 根据设备的产品协议生成接入端点，不把运行时地址冗余保存到设备表。
func (s *ProtocolService) ListDeviceEndpoints(ctx context.Context, deviceKey string) (*protocolV1.DeviceEndpoints, error) {
	device, err := s.repo.DeviceByKey(ctx, deviceKey)
	if err != nil {
		return nil, err
	}
	protocols, err := s.repo.ProductProtocols(ctx, device.ProductID)
	if err != nil {
		return nil, err
	}
	result := &protocolV1.DeviceEndpoints{}
	for _, protocol := range protocols {
		switch protocol.TransportProtocol {
		case string(enum.TransportMQTT):
			broker := s.config.GetString("device_gateway.mqtt.broker")
			if broker == "" {
				broker = s.config.GetString("data.mqtt.broker")
			}
			parsed, _ := url.Parse(broker)
			host := parsed.Hostname()
			port := parsed.Port()
			if host == "" {
				host = "127.0.0.1"
			}
			if port == "" {
				port = "1883"
			}
			result.MQTT = &protocolV1.MQTTEndpoint{
				Host: host, Port: port,
				TelemetryTopic:           fmt.Sprintf("/v1/devices/%s/telemetry", deviceKey),
				AttributesTopic:          fmt.Sprintf("/v1/devices/%s/attributes", deviceKey),
				AttributesSubscribeTopic: fmt.Sprintf("/v1/devices/%s/attributes", deviceKey),
				RPCSubscribeTopic:        fmt.Sprintf("/v1/devices/%s/rpc/request/+", deviceKey),
			}
		case string(enum.TransportHTTP):
			addr := s.config.GetString("device_gateway.http_addr")
			if addr == "" {
				addr = ":8080"
			}
			if strings.HasPrefix(addr, ":") {
				addr = "127.0.0.1" + addr
			}
			baseURL := "http://" + addr
			result.HTTP = &protocolV1.HTTPEndpoint{
				HTTP:         "curl -v -X POST " + baseURL + "/api/v1/" + deviceKey + "/telemetry --header 'Content-Type:application/json' --data '{\"temperature\":25}'",
				RPCSubscribe: baseURL + "/api/v1/" + deviceKey + "/rpc",
			}
		case string(enum.TransportCoAP):
			addr := s.config.GetString("device_gateway.coap_addr")
			if addr == "" {
				addr = ":5683"
			}
			if strings.HasPrefix(addr, ":") {
				addr = "127.0.0.1" + addr
			}
			endpoint := "coap://" + addr + "/api/v1/" + deviceKey + "/telemetry"
			result.CoAP = &protocolV1.CoAPEndpoint{
				CoAP:         `coap-client -v 6 -m POST -t "application/json" -e "{\"temperature\":25}" ` + endpoint,
				Docker:       &protocolV1.CoAPDockerExample{CoAP: `docker run --rm -it --add-host=host.docker.internal:host-gateway thingsboard/coap-clients coap-client -v 6 -m POST -t "application/json" -e "{\"temperature\":25}" coap://host.docker.internal:5683/api/v1/` + deviceKey + `/telemetry`},
				RPCSubscribe: "coap://" + addr + "/api/v1/" + deviceKey + "/rpc",
			}
		default:
			continue
		}
	}
	return result, nil
}

func supportedTransport(value enum.TransportProtocol) bool {
	switch value {
	case enum.TransportMQTT, enum.TransportHTTP, enum.TransportCoAP, enum.TransportTCP, enum.TransportUDP, enum.TransportModbus, enum.TransportGB28181:
		return true
	default:
		return false
	}
}

func supportedApplication(value enum.ApplicationProtocol) bool {
	switch value {
	case enum.ApplicationJSON, enum.ApplicationProtobuf, enum.ApplicationModbus, enum.ApplicationGB28181, enum.ApplicationJT808, enum.ApplicationJT1078, enum.ApplicationCustom:
		return true
	default:
		return false
	}
}
