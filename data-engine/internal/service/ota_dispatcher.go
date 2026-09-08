package service

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"data-engine/internal/model"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
)

// OTADispatcher is data-engine's only MQTT egress path for OTA commands.
type OTADispatcher struct{ client mqtt.Client }

func NewOTADispatcher(config *viper.Viper) (*OTADispatcher, error) {
	broker := config.GetString("mqtt.broker")
	if broker == "" {
		return nil, fmt.Errorf("mqtt.broker is required")
	}
	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID(config.GetString("mqtt.client_id")).SetConnectRetry(true).SetAutoReconnect(true)
	if username := config.GetString("mqtt.username"); username != "" {
		opts.SetUsername(username)
	}
	if password := config.GetString("mqtt.password"); password != "" {
		opts.SetPassword(password)
	}
	tlsConfig, err := otaMQTTTLSConfig(config)
	if err != nil {
		return nil, err
	}
	if tlsConfig != nil {
		opts.SetTLSConfig(tlsConfig)
	}
	client := mqtt.NewClient(opts)
	if token := client.Connect(); !token.WaitTimeout(10 * time.Second) {
		return nil, fmt.Errorf("OTA mqtt connect timeout")
	} else if err := token.Error(); err != nil {
		return nil, fmt.Errorf("connect OTA mqtt: %w", err)
	}
	return &OTADispatcher{client: client}, nil
}

func otaMQTTTLSConfig(config *viper.Viper) (*tls.Config, error) {
	caFile, certFile, keyFile := config.GetString("mqtt.tls.ca_file"), config.GetString("mqtt.tls.cert_file"), config.GetString("mqtt.tls.key_file")
	if caFile == "" && certFile == "" && keyFile == "" {
		return nil, nil
	}
	if certFile == "" || keyFile == "" {
		return nil, fmt.Errorf("OTA MQTT TLS requires cert_file and key_file")
	}
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("load OTA MQTT client certificate: %w", err)
	}
	tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12, Certificates: []tls.Certificate{cert}}
	if caFile != "" {
		pem, err := os.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("read OTA MQTT CA certificate: %w", err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(pem) {
			return nil, fmt.Errorf("parse OTA MQTT CA certificate")
		}
		tlsConfig.RootCAs = pool
	}
	return tlsConfig, nil
}

func (d *OTADispatcher) Dispatch(_ context.Context, command model.OTAUpgradeCommand) error {
	payload, err := json.Marshal(command)
	if err != nil {
		return fmt.Errorf("marshal OTA command: %w", err)
	}
	deviceName := command.DeviceName
	if deviceName == "" {
		deviceName = command.DeviceKey
	}
	topic := fmt.Sprintf("/ota/device/upgrade/%s/%s", command.ProductKey, deviceName)
	token := d.client.Publish(topic, 1, false, payload)
	if !token.WaitTimeout(10 * time.Second) {
		return fmt.Errorf("OTA mqtt publish timeout")
	}
	return token.Error()
}

func (d *OTADispatcher) Close() {
	if d.client != nil && d.client.IsConnected() {
		d.client.Disconnect(250)
	}
}
