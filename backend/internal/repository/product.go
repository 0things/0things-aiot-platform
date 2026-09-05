package repository

import (
	"context"
	"errors"
	"strconv"

	"aiot-backend/internal/dto"
	"aiot-backend/internal/enum"
	"aiot-backend/internal/model"
	"aiot-backend/internal/tenant"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) DB(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *ProductRepository) Find(ctx context.Context, id int64) (*model.Product, error) {
	q := useQuery(r.db)
	product, err := q.Product.WithContext(ctx).Where(q.Product.ID.Eq(id), q.Product.OrganizationID.Eq(tenant.GetOrganizationID(ctx))).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if err := r.fillProtocols(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

func (r *ProductRepository) FindByKey(ctx context.Context, key string) (*model.Product, error) {
	q := useQuery(r.db)
	product, err := q.Product.WithContext(ctx).Where(q.Product.ProductKey.Eq(key), q.Product.OrganizationID.Eq(tenant.GetOrganizationID(ctx))).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if err := r.fillProtocols(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

// fillProtocols 为产品详情补充协议组合，产品列表不加载该关联以避免额外查询。
func (r *ProductRepository) fillProtocols(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Where("product_id = ?", product.ID).Order("id").Find(&product.Protocols).Error
}

func (r *ProductRepository) Create(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(product).Error; err != nil {
			return err
		}
		for i := range product.Protocols {
			product.Protocols[i].ProductID = product.ID
			if err := tx.Create(&product.Protocols[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *ProductRepository) Save(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(product).Error; err != nil {
			return err
		}
		if product.Protocols == nil {
			return nil
		}
		if err := tx.Where("product_id = ?", product.ID).Delete(&model.ProductProtocol{}).Error; err != nil {
			return err
		}
		for i := range product.Protocols {
			product.Protocols[i].ProductID = product.ID
			if err := tx.Create(&product.Protocols[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *ProductRepository) CountDevices(ctx context.Context, productID int64) (int64, error) {
	q := useQuery(r.db)
	return q.Device.WithContext(ctx).Where(q.Device.ProductID.Eq(productID), q.Device.OrganizationID.Eq(tenant.GetOrganizationID(ctx))).Count()
}

func (r *ProductRepository) Delete(ctx context.Context, product *model.Product) error {
	_, err := useQuery(r.db).Product.WithContext(ctx).Where(useQuery(r.db).Product.OrganizationID.Eq(tenant.GetOrganizationID(ctx))).Delete(product)
	return err
}

func (r *ProductRepository) List(ctx context.Context, page, size int, category, status, search string) ([]dto.ProductListItem, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Product{}).
		Select("products.*, categories.name AS category_name").
		Joins("LEFT JOIN categories ON categories.id = products.category_id AND categories.deleted_at IS NULL").
		Where("products.organization_id = ?", tenant.GetOrganizationID(ctx))
	if category != "" {
		if catID, err := strconv.ParseInt(category, 10, 64); err == nil {
			query = query.Where("products.category_id = ?", catID)
		}
	}
	if status != "" {
		query = query.Where("products.status = ?", status)
	}
	if search != "" {
		value := "%" + search + "%"
		query = query.Where("(products.product_key LIKE ? OR products.name LIKE ?)", value, value)
	}
	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var result []dto.ProductListItem
	if err := query.Order("products.id DESC").Offset((page - 1) * size).Limit(size).Find(&result).Error; err != nil {
		return nil, 0, err
	}
	return result, total, nil
}

func (r *ProductRepository) ListOptions(ctx context.Context) ([]dto.ProductOption, error) {
	q := useQuery(r.db)
	var result []dto.ProductOption
	err := q.Product.WithContext(ctx).
		Select(q.Product.ID, q.Product.ProductKey, q.Product.Name, q.Product.NodeType).
		Where(
			q.Product.OrganizationID.Eq(tenant.GetOrganizationID(ctx)),
			q.Product.Status.Eq(string(enum.ProductStatusActive)),
		).
		Order(q.Product.ID.Desc()).
		Scan(&result)
	return result, err
}
