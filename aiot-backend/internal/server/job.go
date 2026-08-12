package server

import (
	"context"
	"0things-backend/internal/job"
	"0things-backend/pkg/log"
)

type JobServer struct {
	log     *log.Logger
	userJob job.UserJob
	eventConsumer *job.DeviceEventConsumer
}

func NewJobServer(
	log *log.Logger,
	userJob job.UserJob,
	eventConsumer *job.DeviceEventConsumer,
) *JobServer {
	return &JobServer{
		log:     log,
		userJob: userJob,
		eventConsumer: eventConsumer,
	}
}

func (j *JobServer) Start(ctx context.Context) error {
	// Tips: If you want job to start as a separate process, just refer to the task implementation and adjust the code accordingly.

	go j.eventConsumer.Start(ctx)
	return j.userJob.KafkaConsumer(ctx)
}
func (j *JobServer) Stop(ctx context.Context) error {
	j.eventConsumer.Stop()
	return nil
}
