package repository

import (
	"context"
	"errors"

	"0things-backend/internal/model"
	"0things-backend/internal/tenant"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *IoTDB
}

func NewProductRepository(db *IoTDB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) DB(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *ProductRepository) Find(ctx context.Context, id int64) (*model.Product, error) {
	q := useIoTQuery(r.db)
	product, err := q.Product.WithContext(ctx).Where(q.Product.ID.Eq(id), q.Product.TenantID.Eq(tenant.GetTenantID(ctx))).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return product, nil
}

func (r *ProductRepository) FindByKey(ctx context.Context, key string) (*model.Product, error) {
	q := useIoTQuery(r.db)
	product, err := q.Product.WithContext(ctx).Where(q.Product.ProductKey.Eq(key), q.Product.TenantID.Eq(tenant.GetTenantID(ctx))).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return product, nil
}

func (r *ProductRepository) Create(ctx context.Context, product *model.Product) error {
	return useIoTQuery(r.db).Product.WithContext(ctx).Create(product)
}

func (r *ProductRepository) Save(ctx context.Context, product *model.Product) error {
	return useIoTQuery(r.db).Product.WithContext(ctx).Save(product)
}

func (r *ProductRepository) CountDevices(ctx context.Context, productID int64) (int64, error) {
	q := useIoTQuery(r.db)
	return q.Device.WithContext(ctx).Where(q.Device.ProductID.Eq(productID), q.Device.TenantID.Eq(tenant.GetTenantID(ctx))).Count()
}

func (r *ProductRepository) Delete(ctx context.Context, product *model.Product) error {
	_, err := useIoTQuery(r.db).Product.WithContext(ctx).Where(useIoTQuery(r.db).Product.TenantID.Eq(tenant.GetTenantID(ctx))).Delete(product)
	return err
}

func (r *ProductRepository) List(ctx context.Context, page, size int, category, status, search string) ([]model.Product, int64, error) {
	q := useIoTQuery(r.db)
	products := q.Product.WithContext(ctx).Where(q.Product.TenantID.Eq(tenant.GetTenantID(ctx)))
	if category != "" {
		products = products.Where(q.Product.Category.Eq(category))
	}
	if status != "" {
		products = products.Where(q.Product.Status.Eq(status))
	}
	if search != "" {
		products = products.Where(field.Or(q.Product.ProductKey.Like("%"+search+"%"), q.Product.Name.Like("%"+search+"%")))
	}
	items, total, err := products.Order(q.Product.CreatedAt.Desc()).FindByPage((page-1)*size, size)
	if err != nil {
		return nil, 0, err
	}
	result := make([]model.Product, len(items))
	for i := range items {
		result[i] = *items[i]
	}
	return result, total, nil
}

func (r *ProductRepository) Restore(ctx context.Context, id int64) error {
	q := useIoTQuery(r.db)
	_, err := q.Product.WithContext(ctx).Unscoped().Where(q.Product.ID.Eq(id), q.Product.TenantID.Eq(tenant.GetTenantID(ctx))).UpdateSimple(q.Product.DeletedAt.Null())
	return err
}
