package port

import (
	"context"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/dto"
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByID(ctx context.Context, id uint64) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id uint64) error
	GetUserByToken(ctx context.Context, token string) (*domain.User, error)
	ExistEmail(ctx context.Context, email string) (bool, error)
}

type UserService interface {
	Register(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUser(ctx context.Context, id uint64) (*domain.User, error)
	ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id uint64) error
	GetProfile(ctx context.Context, id uint64) (*dto.GetProfile, error)
}
