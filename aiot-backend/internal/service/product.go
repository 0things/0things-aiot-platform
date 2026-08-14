package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

	"0things-backend/internal/model"
	"0things-backend/internal/repository"
	"0things-backend/internal/tenant"
)

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

func normalizeProductMetadata(value json.RawMessage) (json.RawMessage, error) {
	if len(value) == 0 {
		return value, nil
	}

	var legacyString string
	if json.Unmarshal(value, &legacyString) == nil {
		if !json.Valid([]byte(legacyString)) {
			return nil, errors.New("invalid metadata")
		}
		return json.RawMessage(legacyString), nil
	}
	if !json.Valid(value) {
		return nil, errors.New("invalid metadata")
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
		product.Status = "active"
	}
	if product.NodeType == "" {
		product.NodeType = "direct"
	}
	product.TenantID = tenant.GetTenantID(ctx)
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
	var err error
	if product.Metadata, err = normalizeProductMetadata(product.Metadata); err != nil {
		return err
	}
	return s.repo.Save(ctx, product)
}

func (s *ProductService) Delete(ctx context.Context, id int64) error {
	count, err := s.repo.CountDevices(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("product has devices")
	}
	product, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, product)
}

func (s *ProductService) List(ctx context.Context, page, size int, category, status, search string) ([]model.Product, int64, error) {
	return s.repo.List(ctx, page, size, category, status, search)
}

func (s *ProductService) Restore(ctx context.Context, id int64) (*model.Product, error) {
	if err := s.repo.Restore(ctx, id); err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}
