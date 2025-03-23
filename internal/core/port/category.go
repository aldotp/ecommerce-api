package port

import (
	"context"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/dto"
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
)

type CategoryRepository interface {
	Finds(ctx context.Context, filter map[string]interface{}) (response []domain.Category, err error)
	FindOne(ctx context.Context, id int) (response *domain.Category, err error)
	Store(ctx context.Context, data *domain.Category) error
	Update(ctx context.Context, id int, updatedData domain.Category) error
	Delete(ctx context.Context, id int) error
}

type CategoryService interface {
	FindOne(ctx context.Context, CategoryID int) (response *domain.Category, err error)
	Store(ctx context.Context, data dto.CategoryRequest) error
	Finds(ctx context.Context, param dto.ListCategoryRequest) (response []domain.Category, err error)
	Update(ctx context.Context, id int, data dto.CategoryRequest) error
	Delete(ctx context.Context, id int) error
}
