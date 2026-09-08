package task

import (
	"transport-mqtt/internal/repository"
	"transport-mqtt/pkg/jwt"
	"transport-mqtt/pkg/log"
	"transport-mqtt/pkg/sid"
)

type Task struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
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
