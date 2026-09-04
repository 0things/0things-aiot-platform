package repository

import (
	"context"
	"errors"

	"aiot-backend/internal/model"
	"aiot-backend/internal/tenant"

	"gorm.io/gorm"
)

type SceneLinkageRepository struct {
	db *gorm.DB
}

func NewSceneLinkageRepository(db *gorm.DB) *SceneLinkageRepository {
	return &SceneLinkageRepository{db: db}
}

func (r *SceneLinkageRepository) DB(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *SceneLinkageRepository) Find(ctx context.Context, id int64) (*model.SceneLinkage, error) {
	q := useQuery(r.db)
	scene, err := q.SceneLinkage.WithContext(ctx).Where(q.SceneLinkage.ID.Eq(id), q.SceneLinkage.OrganizationID.Eq(tenant.GetOrganizationID(ctx))).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return scene, nil
}

func (r *SceneLinkageRepository) List(ctx context.Context, page, size int, search string, enable int) ([]model.SceneLinkage, int64, error) {
	q := useQuery(r.db)
	scenes := q.SceneLinkage.WithContext(ctx).Where(q.SceneLinkage.OrganizationID.Eq(tenant.GetOrganizationID(ctx)))
	if search != "" {
		scenes = scenes.Where(q.SceneLinkage.Name.Like("%" + search + "%"))
	}
	if enable == 0 || enable == 1 {
		scenes = scenes.Where(q.SceneLinkage.Enable.Eq(enable))
	}
	items, total, err := scenes.Order(q.SceneLinkage.CreatedAt.Desc()).FindByPage((page-1)*size, size)
	if err != nil {
		return nil, 0, err
	}
	result := make([]model.SceneLinkage, len(items))
	for i := range items {
		result[i] = *items[i]
	}
	return result, total, nil
}

func (r *SceneLinkageRepository) Create(ctx context.Context, scene *model.SceneLinkage) error {
	scene.OrganizationID = tenant.GetOrganizationID(ctx)
	return useQuery(r.db).SceneLinkage.WithContext(ctx).Create(scene)
}

func (r *SceneLinkageRepository) Save(ctx context.Context, scene *model.SceneLinkage) error {
	if _, err := r.Find(ctx, scene.ID); err != nil {
		return err
	}
	return useQuery(r.db).SceneLinkage.WithContext(ctx).Save(scene)
}

func (r *SceneLinkageRepository) Delete(ctx context.Context, id int64) error {
	scene, err := r.Find(ctx, id)
	if err != nil {
		return err
	}
	_, err = useQuery(r.db).SceneLinkage.WithContext(ctx).Delete(scene)
	return err
}
