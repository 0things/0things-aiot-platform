package task

import (
	"context"

	"aiot-backend/internal/repository"
)

type OTATask interface {
	// DispatchPending 扫描所有仍有 pending 设备升级记录的升级包，
	// 将 pending 记录推进为 in_progress（下发升级命令）。由周期任务定时触发。
	DispatchPending(ctx context.Context) error
}

func NewOTATask(task *Task, otaRepo *repository.OTARepository) OTATask {
	return &otaTask{
		otaRepo: otaRepo,
		Task:    task,
	}
}

type otaTask struct {
	otaRepo *repository.OTARepository
	*Task
}

func (t otaTask) DispatchPending(ctx context.Context) error {
	pkgIDs, err := t.otaRepo.PendingPackageIDs(ctx)
	if err != nil {
		return err
	}
	for _, id := range pkgIDs {
		if _, err := t.otaRepo.DispatchPending(ctx, id); err != nil {
			return err
		}
	}
	return nil
}
