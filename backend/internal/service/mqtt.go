package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"aiot-backend/pkg/log"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var ErrMQTTNotConnected = errors.New("mqtt client not connected")

type MQTTServiceInterface interface {
	Publish(ctx context.Context, topic string, qos byte, payload []byte) error
	Subscribe(ctx context.Context, topic string, qos byte, handler mqtt.MessageHandler) error
	Close()
}

type MQTTService struct {
	client        mqtt.Client
	logger        *log.Logger
	mu            sync.RWMutex
	subscriptions map[string]mqtt.MessageHandler
}

func NewMQTTService(config *viper.Viper, logger *log.Logger) (*MQTTService, func(), error) {
	broker := config.GetString("data.mqtt.broker")
	if broker == "" {
		broker = "tcp://127.0.0.1:1883"
	}
	options := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(config.GetString("data.mqtt.client_id"))
	if options.ClientID == "" {
		options.ClientID = "aiot-ota-" + uuid.NewString()
	}
	if username := config.GetString("data.mqtt.username"); username != "" {
		options.SetUsername(username)
	}
	if password := config.GetString("data.mqtt.password"); password != "" {
		options.SetPassword(password)
	}
	options.SetAutoReconnect(true).SetConnectRetry(true).SetConnectRetryInterval(5 * time.Second)
	options.SetConnectTimeout(5 * time.Second)
	options.OnConnectionLost = func(_ mqtt.Client, err error) {
		logger.Warn("MQTT connection lost", zap.Error(err))
	}

	client := mqtt.NewClient(options)
	svc := &MQTTService{client: client, logger: logger, subscriptions: make(map[string]mqtt.MessageHandler)}
	options.OnConnect = func(client mqtt.Client) {
		svc.mu.RLock()
		defer svc.mu.RUnlock()
		for topic, handler := range svc.subscriptions {
			if token := client.Subscribe(topic, 1, handler); token.Wait() && token.Error() != nil {
				logger.Warn("MQTT resubscribe failed", zap.String("topic", topic), zap.Error(token.Error()))
			}
		}
	}
	// 首次连接只等待有限时间，避免 MQTT 不可达时阻塞整个 HTTP 服务启动。
	if token := client.Connect(); !token.WaitTimeout(5 * time.Second) {
		logger.Warn("MQTT broker connection timed out; reconnect will continue in background", zap.String("broker", broker))
	} else if token.Error() != nil {
		logger.Warn("MQTT broker unavailable; reconnect will continue in background", zap.String("broker", broker), zap.Error(token.Error()))
	}
	return svc, func() { svc.Close() }, nil
}

func (s *MQTTService) Publish(ctx context.Context, topic string, qos byte, payload []byte) error {
	if !s.client.IsConnected() {
		return ErrMQTTNotConnected
	}
	token := s.client.Publish(topic, qos, false, payload)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-token.Done():
		if err := token.Error(); err != nil {
			return fmt.Errorf("publish mqtt message: %w", err)
		}
		return nil
	}
}

func (s *MQTTService) Subscribe(ctx context.Context, topic string, qos byte, handler mqtt.MessageHandler) error {
	s.mu.Lock()
	s.subscriptions[topic] = handler
	s.mu.Unlock()
	if !s.client.IsConnected() {
		return nil
	}
	token := s.client.Subscribe(topic, qos, handler)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-token.Done():
		if err := token.Error(); err != nil {
			return err
		}
		return nil
	}
}

func (s *MQTTService) Close() {
	if s.client != nil && s.client.IsConnected() {
		s.client.Disconnect(250)
	}
}
