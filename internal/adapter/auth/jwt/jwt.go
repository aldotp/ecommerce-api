package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/config"
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
	"github.com/aldotp/ecommerce-go-api/internal/core/port"
	"github.com/dgrijalva/jwt-go"
)

type JWTToken struct {
}

// New creates a new JWT instance with a specified duration.
func New() (port.TokenInterface, error) {
	return &JWTToken{}, nil
}

// CreateToken generates a new JWT access token for the given user.
func (pt *JWTToken) GenerateAccessToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Minute * time.Duration(config.AccessTokenExpired())).Unix(),
		"role":    user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(config.SecretKey()))
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return signedToken, nil
}

// GenerateRefreshToken creates a long-lived refresh token for the given user.
func (pt *JWTToken) GenerateRefreshToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Minute * time.Duration(config.RefreshTokenExpired())).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(config.RefreshKey()))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return signedToken, nil
}

// VerifyAccessToken validates a JWT access token and returns the decoded payload.
func (pt *JWTToken) VerifyAccessToken(encodedToken string) (*domain.TokenPayload, error) {
	secretKey := config.SecretKey()

	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse access token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid or expired access token")
	}

	return &domain.TokenPayload{
		UserID: int(claims["user_id"].(float64)),
		Email:  claims["email"].(string),
		Role:   domain.UserRole(claims["role"].(string)),
	}, nil
}

// VerifyRefreshToken validates a JWT refresh token and checks its expiration status.
func (pt *JWTToken) VerifyRefreshToken(ctx context.Context, refreshToken string) (*domain.RefreshTokenPayload, bool, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(config.RefreshKey()), nil
	})

	if err != nil {
		return nil, false, fmt.Errorf("failed to parse refresh token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, false, errors.New("invalid or expired refresh token")
	}

	almostExpired := claims.VerifyExpiresAt(time.Now().Add(24*time.Hour).Unix(), true)

	payload := &domain.RefreshTokenPayload{
		UserID: uint64(claims["user_id"].(float64)),
	}

	return payload, almostExpired, nil
}

func (*JWTToken) GenerateActivationToken(userID uint64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.SecretKey()))
}
