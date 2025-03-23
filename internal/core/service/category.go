package service

import (
	"context"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/dto"
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
	"github.com/aldotp/ecommerce-go-api/internal/core/port"
	"github.com/aldotp/ecommerce-go-api/pkg/consts"
)

type CategoryService struct {
	repo  port.CategoryRepository
	cache port.CacheInterface
}

func NewCategoryService(repo port.CategoryRepository, cache port.CacheInterface) *CategoryService {
	return &CategoryService{
		repo:  repo,
		cache: cache,
	}
}

// Store a new category
func (s *CategoryService) Store(ctx context.Context, data dto.CategoryRequest) error {
	category := &domain.Category{
		Name: data.Name,
	}

	return s.repo.Store(ctx, category)
}

// Find a single category by ID
func (s *CategoryService) FindOne(ctx context.Context, categoryID int) (*domain.Category, error) {
	data, err := s.repo.FindOne(ctx, (categoryID))
	if err != nil || data == nil {
		return nil, consts.ErrDataNotFound
	}

	return data, nil
}

// Find multiple categorys
func (s *CategoryService) Finds(ctx context.Context, param dto.ListCategoryRequest) ([]domain.Category, error) {
	filter := make(map[string]interface{})

	if param.Page != 0 {
		filter["page"] = param.Page
	}
	if param.PageSize != 0 {
		filter["page_size"] = param.PageSize
	}

	return s.repo.Finds(ctx, filter)
}

// Update a category by ID
func (s *CategoryService) Update(ctx context.Context, id int, data dto.CategoryRequest) error {
	updatedData := domain.Category{
		Name: data.Name,
	}

	return s.repo.Update(ctx, id, updatedData)
}

// Delete a category by ID
func (s *CategoryService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
