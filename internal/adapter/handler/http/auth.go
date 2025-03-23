package http

import (
	"net/http"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/dto"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/helper"
	"github.com/aldotp/ecommerce-go-api/internal/core/port"
	"github.com/aldotp/ecommerce-go-api/pkg/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthHandler represents the HTTP handler for authentication-related requests
type AuthHandler struct {
	svc    port.AuthService
	logger *zap.Logger
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(svc port.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		svc:    svc,
		logger: logger,
	}
}

// Login godoc
//
//	@Summary		User Login
//	@Description	Authenticate user and return an access token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginRequest	true	"Login request body"
//	@Success		200		{object}	dto.AuthResponse	"Successfully logged in"
//	@Failure		400		{object}	util.ErrorResponse	"Bad request (validation error)"
//	@Failure		401		{object}	util.ErrorResponse	"Unauthorized error"
//	@Failure		500		{object}	util.ErrorResponse	"Internal server error"
//	@Router			/api/v1/auth/login [post]
func (ah *AuthHandler) Login(c *gin.Context) {
	var request dto.LoginRequest

	// Log request masuk
	ah.logger.Info("Login request received")

	if err := c.ShouldBindJSON(&request); err != nil {
		ah.logger.Warn("Failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	// Log email yang dicoba login
	ah.logger.Info("User attempting login", zap.String("email", request.Email))

	data, err := ah.svc.Login(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		ah.logger.Error("Login failed", zap.String("email", request.Email), zap.Error(err))
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	ah.logger.Info("Login successful", zap.String("email", request.Email))
	c.JSON(http.StatusOK, util.APIResponse("Successfully logged in", http.StatusOK, "success", data))
}

// RefreshToken godoc
//
//	@Summary		Refresh Access Token
//	@Description	Refresh expired access token using refresh token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RefreshTokenRequest	true	"Refresh token request body"
//	@Success		200		{object}	dto.AuthResponse	"Successfully refreshed token"
//	@Failure		400		{object}	util.ErrorResponse	"Bad request (validation error)"
//	@Failure		401		{object}	util.ErrorResponse	"Unauthorized error"
//	@Failure		500		{object}	util.ErrorResponse	"Internal server error"
//	@Router			/api/v1/auth/refresh-token [post]
func (ah *AuthHandler) RefreshToken(c *gin.Context) {
	var request dto.RefreshTokenRequest

	// Log request masuk
	ah.logger.Info("Refresh token request received")

	if err := c.ShouldBindJSON(&request); err != nil {
		ah.logger.Warn("Failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	// Log token yang digunakan untuk refresh
	ah.logger.Info("Attempting token refresh", zap.String("refresh_token", request.RefreshToken))

	token, err := ah.svc.RefreshToken(c.Request.Context(), request.RefreshToken)
	if err != nil {
		ah.logger.Error("Token refresh failed", zap.String("refresh_token", request.RefreshToken), zap.Error(err))
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	ah.logger.Info("Token refreshed successfully")
	c.JSON(http.StatusOK, util.APIResponse("Successfully refreshed token", http.StatusOK, "success", map[string]interface{}{"token": token}))
}
