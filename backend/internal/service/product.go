package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

	"aiot-backend/internal/dto"
	"aiot-backend/internal/enum"
	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
	"aiot-backend/internal/tenant"
)

type ProductServiceInterface interface {
	Create(ctx context.Context, product *model.Product) (*model.Product, error)
	Get(ctx context.Context, id int64) (*model.Product, error)
	GetByKey(ctx context.Context, key string) (*model.Product, error)
	Save(ctx context.Context, product *model.Product) error
	// Deprecated: 外部删除接口统一使用 productKey，保留供旧测试兼容。
	Delete(ctx context.Context, id int64) error
	DeleteByKey(ctx context.Context, key string) error
	List(ctx context.Context, page, size int, category, status, search string) ([]dto.ProductListItem, int64, error)
	// Deprecated: 恢复路由已移除，保留旧测试和内部迁移兼容。
	Restore(ctx context.Context, id int64) (*model.Product, error)
	// Deprecated: 恢复路由已移除，保留旧测试和内部迁移兼容。
	RestoreByKey(ctx context.Context, key string) (*model.Product, error)
}

// ProductListItem 是产品列表专用结果，包含分类名称但不污染产品领域模型。

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func productKey() string {
	bytes := make([]byte, 8)
	_, _ = rand.Read(bytes)
	return "P" + strings.ToUpper(hex.EncodeToString(bytes))
}

func normalizeProductMetadata(value string) (string, error) {
	if len(value) == 0 {
		return value, nil
	}

	var legacyString string
	if json.Unmarshal([]byte(value), &legacyString) == nil {
		if !json.Valid([]byte(legacyString)) {
			return "", errors.New("invalid metadata")
		}
		return legacyString, nil
	}
	if !json.Valid([]byte(value)) {
		return "", errors.New("invalid metadata")
	}
	return value, nil
}

func (s *ProductService) Create(ctx context.Context, product *model.Product) (*model.Product, error) {
	if product.Name == "" {
		return nil, errors.New("name is required")
	}
	var err error
	if product.Metadata, err = normalizeProductMetadata(product.Metadata); err != nil {
		return nil, err
	}
	if product.ProductKey == "" {
		product.ProductKey = productKey()
	}
	if product.Status == "" {
		product.Status = string(enum.ProductStatusActive)
	}
	if product.Status != string(enum.ProductStatusActive) && product.Status != string(enum.ProductStatusInactive) && product.Status != string(enum.ProductStatusArchived) {
		return nil, errors.New("invalid product status")
	}
	if product.NodeType == "" {
		product.NodeType = "direct"
	}
	product.OrganizationID = tenant.GetOrganizationID(ctx)
	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) Get(ctx context.Context, id int64) (*model.Product, error) {
	return s.repo.Find(ctx, id)
}

func (s *ProductService) GetByKey(ctx context.Context, key string) (*model.Product, error) {
	return s.repo.FindByKey(ctx, key)
}

func (s *ProductService) Save(ctx context.Context, product *model.Product) error {
	if product.Status == "" {
		product.Status = string(enum.ProductStatusActive)
	}
	if product.Status != string(enum.ProductStatusActive) && product.Status != string(enum.ProductStatusInactive) && product.Status != string(enum.ProductStatusArchived) {
		return errors.New("invalid product status")
	}
	var err error
	if product.Metadata, err = normalizeProductMetadata(product.Metadata); err != nil {
		return err
	}
	return s.repo.Save(ctx, product)
}

func (s *ProductService) DeleteByKey(ctx context.Context, key string) error {
	product, err := s.GetByKey(ctx, key)
	if err != nil {
		return err
	}
	count, err := s.repo.CountDevices(ctx, product.ID)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("product has devices")
	}
	return s.repo.Delete(ctx, product)
}

// Deprecated: 外部删除接口统一使用 productKey。
func (s *ProductService) Delete(ctx context.Context, id int64) error {
	product, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	count, err := s.repo.CountDevices(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("product has devices")
	}
	return s.repo.Delete(ctx, product)
}

func (s *ProductService) List(ctx context.Context, page, size int, category, status, search string) ([]dto.ProductListItem, int64, error) {
	return s.repo.List(ctx, page, size, category, status, search)
}

// Deprecated: 产品恢复接口已下线。
func (s *ProductService) Restore(ctx context.Context, id int64) (*model.Product, error) {
	if err := s.repo.Restore(ctx, id); err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

// Deprecated: 产品恢复接口已下线。
func (s *ProductService) RestoreByKey(ctx context.Context, key string) (*model.Product, error) {
	product, err := s.GetByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	return s.Restore(ctx, product.ID)
}
