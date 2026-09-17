package service

import (
	"aiot-backend/internal/repository"
	"aiot-backend/pkg/log"
)

type Service struct {
	logger *log.Logger
	tm     repository.Transaction
}

func NewService(
	tm repository.Transaction,
	logger *log.Logger,
) *Service {
	return &Service{
		logger: logger,
		tm:     tm,
	}
}
