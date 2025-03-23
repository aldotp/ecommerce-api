package service

import (
	"context"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/dto"
	"github.com/aldotp/ecommerce-go-api/internal/core/port"
	"github.com/aldotp/ecommerce-go-api/pkg/consts"
	"github.com/aldotp/ecommerce-go-api/pkg/util"
	"go.uber.org/zap"
)

type AuthService struct {
	repo port.UserRepository
	ts   port.TokenInterface
	log  *zap.Logger
}

func NewAuthService(repo port.UserRepository, ts port.TokenInterface, log *zap.Logger) *AuthService {
	return &AuthService{
		repo,
		ts,
		log,
	}
}

func (as *AuthService) Login(ctx context.Context, email, password string) (dto.LoginResponse, error) {
	user, err := as.repo.GetUserByEmail(ctx, email)
	if err != nil {
		as.log.Error(err.Error())
		if err == consts.ErrDataNotFound {
			return dto.LoginResponse{}, consts.ErrInvalidCredentials
		}
		return dto.LoginResponse{}, consts.ErrInternal
	}

	err = util.ComparePassword(password, user.Password)
	if err != nil {
		as.log.Error(err.Error())
		return dto.LoginResponse{}, consts.ErrInvalidCredentials
	}

	accessToken, err := as.ts.GenerateAccessToken(user)
	if err != nil {
		as.log.Error(err.Error())
		return dto.LoginResponse{}, consts.ErrTokenCreation
	}

	refreshToken, err := as.ts.GenerateRefreshToken(user)
	if err != nil {
		as.log.Error(err.Error())
		return dto.LoginResponse{}, consts.ErrTokenCreation
	}

	return dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (as *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	payload, _, err := as.ts.VerifyRefreshToken(ctx, refreshToken)
	if err != nil {
		as.log.Error(err.Error())
		return "", consts.ErrInvalidSignature
	}

	user, err := as.repo.GetUserByID(ctx, payload.UserID)
	if err != nil {
		as.log.Error(err.Error())
		if err == consts.ErrDataNotFound {
			return "", consts.ErrInvalidToken
		}
		return "", consts.ErrInternal
	}

	accessToken, err := as.ts.GenerateAccessToken(user)
	if err != nil {
		as.log.Error(err.Error())
		return "", consts.ErrTokenCreation
	}

	return accessToken, nil
}
