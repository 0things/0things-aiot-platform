package server

import (
	"aiot-backend/internal/job"
	"aiot-backend/pkg/log"
	"context"
	"go.uber.org/zap"
	"sync"
)

type JobServer struct {
	log           *log.Logger
	userJob       job.UserJob
	eventConsumer *job.DeviceEventConsumer
	otaConsumer   *job.OTACommandConsumer
	cancel        context.CancelFunc
	mu            sync.Mutex
}

func NewJobServer(
	log *log.Logger,
	userJob job.UserJob,
	eventConsumer *job.DeviceEventConsumer,
	otaConsumer *job.OTACommandConsumer,
) *JobServer {
	return &JobServer{
		log:           log,
		userJob:       userJob,
		eventConsumer: eventConsumer,
		otaConsumer:   otaConsumer,
	}
}

func (j *JobServer) Start(ctx context.Context) error {
	// 为所有后台消费者派生独立上下文，停止服务时统一取消，避免消费者泄漏。

	ctx, cancel := context.WithCancel(ctx)
	j.mu.Lock()
	j.cancel = cancel
	j.mu.Unlock()
	go j.eventConsumer.Start(ctx)
	go func() {
		if err := j.otaConsumer.Start(ctx); err != nil && ctx.Err() == nil {
			j.log.Error("OTA command consumer stopped", zap.Error(err))
		}
	}()
	return j.userJob.KafkaConsumer(ctx)
}
func (j *JobServer) Stop(ctx context.Context) error {
	j.mu.Lock()
	if j.cancel != nil {
		j.cancel()
	}
	j.mu.Unlock()
	j.eventConsumer.Stop()
	j.otaConsumer.Stop()
	return nil
}
