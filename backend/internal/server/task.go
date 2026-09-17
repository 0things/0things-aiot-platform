package server

import (
	"aiot-backend/pkg/log"
	"context"
)

type TaskServer struct {
	log *log.Logger
}

func NewTaskServer(log *log.Logger) *TaskServer {
	return &TaskServer{
		log: log,
	}
}
func (t *TaskServer) Start(ctx context.Context) error {
	t.log.Info("TaskServer started without local identity tasks")
	return nil
}
func (t *TaskServer) Stop(ctx context.Context) error {
	t.log.Info("TaskServer stop...")
	return nil
}
