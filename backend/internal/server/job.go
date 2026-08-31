package server

import (
	"context"
	"sync"

	"aiot-backend/pkg/log"
)

// JobServer 负责管理 backend 内部的后台定时任务与常驻工作流（如定期清理、定时轮询等）。
type JobServer struct {
	log    *log.Logger
	cancel context.CancelFunc
	mu     sync.Mutex
}

// NewJobServer 初始化 JobServer。
func NewJobServer(log *log.Logger) *JobServer {
	return &JobServer{
		log: log,
	}
}

// Start 启动后台任务调度器并监听系统上下文取消。
func (j *JobServer) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	j.mu.Lock()
	j.cancel = cancel
	j.mu.Unlock()

	j.log.Info("backend JobServer started")
	<-ctx.Done()
	return nil
}

// Stop 优雅停止后台所有正在执行的 Job。
func (j *JobServer) Stop(ctx context.Context) error {
	j.mu.Lock()
	if j.cancel != nil {
		j.cancel()
	}
	j.mu.Unlock()

	j.log.Info("backend JobServer stopped gracefully")
	return nil
}
