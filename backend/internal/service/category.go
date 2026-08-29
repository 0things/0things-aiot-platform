package service

import (
	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
	"context"
)

type CategoryServiceInterface interface {
	Tree(context.Context) ([]model.Category, error)
}
type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Tree(ctx context.Context) ([]model.Category, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	children := make(map[int64][]model.Category)
	var roots []model.Category
	for _, item := range items {
		if item.ParentID == nil {
			roots = append(roots, item)
		} else {
			children[*item.ParentID] = append(children[*item.ParentID], item)
		}
	}
	var attach func([]model.Category) []model.Category
	attach = func(nodes []model.Category) []model.Category {
		for i := range nodes {
			nodes[i].Children = attach(children[nodes[i].ID])
		}
		return nodes
	}
	return attach(roots), nil
}
