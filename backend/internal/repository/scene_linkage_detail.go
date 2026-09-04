package repository

import (
	"context"
	"errors"

	"aiot-backend/internal/model"

	"gorm.io/gorm"
)

type SceneLinkageDetailRepository struct {
	db *gorm.DB
}

func NewSceneLinkageDetailRepository(db *gorm.DB) *SceneLinkageDetailRepository {
	return &SceneLinkageDetailRepository{db: db}
}

func (r *SceneLinkageDetailRepository) FindBySceneID(ctx context.Context, sceneID int64) (*model.SceneLinkageDetail, error) {
	q := useQuery(r.db)
	detail, err := q.SceneLinkageDetail.WithContext(ctx).Where(q.SceneLinkageDetail.SceneID.Eq(sceneID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return detail, nil
}

func (r *SceneLinkageDetailRepository) Create(ctx context.Context, detail *model.SceneLinkageDetail) error {
	return useQuery(r.db).SceneLinkageDetail.WithContext(ctx).Create(detail)
}

func (r *SceneLinkageDetailRepository) Save(ctx context.Context, detail *model.SceneLinkageDetail) error {
	return useQuery(r.db).SceneLinkageDetail.WithContext(ctx).Save(detail)
}
