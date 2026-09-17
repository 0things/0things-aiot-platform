package task

import (
	"aiot-backend/internal/repository"
	"aiot-backend/pkg/log"
	"aiot-backend/pkg/sid"
)

type Task struct {
	logger *log.Logger
	sid    *sid.Sid
	tm     repository.Transaction
}

func NewTask(
	tm repository.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
) *Task {
	return &Task{
		logger: logger,
		sid:    sid,
		tm:     tm,
	}
}
