package task

import (
	"0things-backend/internal/repository"
	"0things-backend/pkg/jwt"
	"0things-backend/pkg/log"
	"0things-backend/pkg/sid"
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
