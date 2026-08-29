// device-gateway is the dedicated process boundary for device protocol
// adapters. Protocol listeners are added here without coupling them to the
// management HTTP server.
package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"aiot-backend/internal/enum"
	"aiot-backend/internal/service"
	"aiot-backend/internal/transport"
	coaptransport "aiot-backend/internal/transport/coap"
	httptransport "aiot-backend/internal/transport/http"
	kafkatransport "aiot-backend/internal/transport/kafka"
	mqtttransport "aiot-backend/internal/transport/mqtt"
	tcptransport "aiot-backend/internal/transport/tcp"
	"aiot-backend/pkg/config"
	applog "aiot-backend/pkg/log"
)

func main() {
	confPath := flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*confPath)
	logger := applog.NewLog(conf)
	kafka, cleanupKafka, err := service.NewKafkaService(conf, logger)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanupKafka()
	addr := conf.GetString("device_gateway.http_addr")
	if addr == "" {
		addr = ":8081"
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	mqttAdapter := mqtttransport.NewAdapter(conf)
	adapters := []transport.Adapter{httptransport.NewAdapter(addr), mqttAdapter}
	downlinkAdapters := []transport.Adapter{mqttAdapter}
	if coapAddr := conf.GetString("device_gateway.coap_addr"); coapAddr != "" {
		coapAdapter := coaptransport.NewAdapter(coapAddr)
		adapters = append(adapters, coapAdapter)
		downlinkAdapters = append(downlinkAdapters, coapAdapter)
	}
	if tcpAddr := conf.GetString("device_gateway.tcp_addr"); tcpAddr != "" {
		tcpAdapter := tcptransport.NewAdapter(tcpAddr)
		adapters = append(adapters, tcpAdapter)
		downlinkAdapters = append(downlinkAdapters, tcpAdapter)
	}
	registry, err := transport.NewRegistry(downlinkAdapters...)
	if err != nil {
		log.Fatal(err)
	}
	commandConsumer, err := kafkatransport.NewOTACommandConsumer(conf, registry)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		if err := commandConsumer.Start(ctx); err != nil && ctx.Err() == nil {
			log.Print(err)
		}
	}()
	gateway := transport.NewGateway(adapters, func(ctx context.Context, message transport.DeviceMessage) error {
		return kafka.ProduceJSON(ctx, enum.KafkaTopicDeviceMessageV1, message.DeviceKey, message)
	})
	if err := gateway.Start(ctx); err != nil {
		log.Fatal(err)
	}
	<-ctx.Done()
	commandConsumer.Stop()
	_ = gateway.Stop(context.Background())
}
