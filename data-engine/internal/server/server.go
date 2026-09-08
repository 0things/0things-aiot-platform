package server

import (
	"context"

	"data-engine/internal/consumer"
	"data-engine/pkg/log"
	"data-engine/pkg/server"
)

type DataEngineServer struct {
	consumerManager *consumer.Manager
	logger          *log.Logger
}

var _ server.Server = (*DataEngineServer)(nil)

func NewDataEngineServer(consumerManager *consumer.Manager, logger *log.Logger) *DataEngineServer {
	return &DataEngineServer{
		consumerManager: consumerManager,
		logger:          logger,
	}
}

func (s *DataEngineServer) Start(ctx context.Context) error {
	s.logger.Info("starting data-engine event consumer manager...")
	return s.consumerManager.Start(ctx)
}

func (s *DataEngineServer) Stop(ctx context.Context) error {
	s.logger.Info("stopping data-engine server...")
	return nil
}
