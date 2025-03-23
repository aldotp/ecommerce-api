package dto

import (
	"time"

	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
)

type UserResponse struct {
	ID        uint64    `json:"id" example:"1"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"test@example.com"`
	CreatedAt time.Time `json:"created_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"1970-01-01T00:00:00Z"`
}

func NewUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"test@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"12345678" minLength:"8"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required" example:"John Doe"`
	Email    string `json:"email" binding:"required,email" example:"test@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"12345678"`
}

type ListUserRequest struct {
	Page  uint64 `form:"page"`
	Limit uint64 `form:"limit"`
}

type GetUserRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

type UpdateUserRequest struct {
	Name     string          `json:"name" binding:"omitempty,required" example:"John Doe"`
	Email    string          `json:"email" binding:"omitempty,required,email" example:"test@example.com"`
	Password string          `json:"password" binding:"omitempty,required,min=8" example:"12345678"`
	Role     domain.UserRole `json:"role" binding:"omitempty,required,user_role" example:"admin"`
}

type DeleteUserRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

type ActivationRequest struct {
	Token string `uri:"token" binding:"required" form:"token"`
}

type GetProfile struct {
	ID        uint64    `json:"id" example:"1"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"test@example.com"`
	Role      string    `json:"role" example:"admin"`
	CreatedAt time.Time `json:"created_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"1970-01-01T00:00:00Z"`
}
