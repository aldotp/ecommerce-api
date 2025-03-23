package service

import (
	"context"
	"strings"
	"time"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/dto"
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
	"github.com/aldotp/ecommerce-go-api/internal/core/port"
	"github.com/aldotp/ecommerce-go-api/pkg/consts"
)

type ProductService struct {
	repo  port.ProductRepository
	cache port.CacheInterface
}

func NewProductService(repo port.ProductRepository, cache port.CacheInterface) *ProductService {
	return &ProductService{
		repo:  repo,
		cache: cache,
	}
}

// Store a new product
func (s *ProductService) Store(ctx context.Context, data dto.ProductRequest) error {
	product := &domain.Product{
		Name:        data.Name,
		Description: data.Description,
		Price:       data.Price,
		Stock:       data.Stock,
		CategoryID:  data.CategoryID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.repo.Store(ctx, product)
}

// Find a single product by ID
func (s *ProductService) FindOne(ctx context.Context, productID int) (*domain.Product, error) {
	product, err := s.repo.FindOne(ctx, (productID))
	if err != nil || product == nil {
		return nil, consts.ErrDataNotFound
	}

	return product, nil
}

// Find multiple products
func (s *ProductService) Finds(ctx context.Context, param dto.ListProductRequest) ([]domain.Product, error) {
	filter := make(map[string]interface{})

	if param.Page != 0 {
		filter["page"] = param.Page
	}
	if param.PageSize != 0 {
		filter["page_size"] = param.PageSize
	}

	if param.Search != "" {
		filter["search"] = strings.ToLower(param.Search)
	}

	return s.repo.Finds(ctx, filter)
}

// Update a product by ID
func (s *ProductService) Update(ctx context.Context, id int, data dto.ProductRequest) error {
	updatedData := domain.Product{
		Name:      data.Name,
		Price:     data.Price,
		Stock:     data.Stock,
		UpdatedAt: time.Now(),
	}

	return s.repo.Update(ctx, id, updatedData)
}

// Delete a product by ID
func (s *ProductService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
