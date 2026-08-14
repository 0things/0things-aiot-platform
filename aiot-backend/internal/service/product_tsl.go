package service

import (
	"context"

	"0things-backend/internal/model"
	"0things-backend/internal/repository"
)

type ProductTSLServiceInterface interface {
	Get(ctx context.Context, productKey string) (*model.ProductTSL, error)
	Upsert(ctx context.Context, productKey, content string) error
	Delete(ctx context.Context, productKey string) error
}

type ProductTSLService struct {
	products *repository.ProductRepository
	tsls     *repository.ProductTSLRepository
}

func NewProductTSLService(products *repository.ProductRepository, tsls *repository.ProductTSLRepository) *ProductTSLService {
	return &ProductTSLService{products: products, tsls: tsls}
}

func (s *ProductTSLService) Get(ctx context.Context, productKey string) (*model.ProductTSL, error) {
	product, err := s.products.FindByKey(ctx, productKey)
	if err != nil {
		return nil, err
	}
	return s.tsls.FindByProductID(ctx, product.ID)
}

func (s *ProductTSLService) Upsert(ctx context.Context, productKey, content string) error {
	product, err := s.products.FindByKey(ctx, productKey)
	if err != nil {
		return err
	}
	tsl, err := s.tsls.FindByProductID(ctx, product.ID)
	if err == repository.ErrNotFound {
		productID := product.ID
		return s.tsls.Create(ctx, &model.ProductTSL{ProductID: &productID, TSL: content})
	}
	if err != nil {
		return err
	}
	tsl.TSL = content
	return s.tsls.Save(ctx, tsl)
}

func (s *ProductTSLService) Delete(ctx context.Context, productKey string) error {
	tsl, err := s.Get(ctx, productKey)
	if err != nil {
		return err
	}
	return s.tsls.Delete(ctx, tsl)
}
