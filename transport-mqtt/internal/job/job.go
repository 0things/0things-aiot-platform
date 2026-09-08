package job

import (
	"transport-mqtt/internal/repository"
	"transport-mqtt/pkg/jwt"
	"transport-mqtt/pkg/log"
	"transport-mqtt/pkg/sid"
)

type Job struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     repository.Transaction
}

func NewJob(
	tm repository.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
) *Job {
	return &Job{
		logger: logger,
		sid:    sid,
		tm:     tm,
	}
}
