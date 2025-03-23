package http

import (
	"net/http"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/dto"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/helper"
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
	"github.com/aldotp/ecommerce-go-api/internal/core/port"
	"github.com/aldotp/ecommerce-go-api/pkg/consts"
	"github.com/aldotp/ecommerce-go-api/pkg/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserHandler represents the HTTP handler for user-related requests
type UserHandler struct {
	svc    port.UserService
	logger *zap.Logger
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(svc port.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		svc:    svc,
		logger: logger,
	}
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Create a new user account with default role "cashier"
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			registerRequest	body		dto.RegisterRequest	true	"Register request"
//	@Success		200				{object}	dto.UserResponse	"User created"
//	@Failure		400				{object}	util.ErrorResponse	"Validation error"
//	@Failure		409				{object}	util.ErrorResponse	"Data conflict error"
//	@Failure		500				{object}	util.ErrorResponse	"Internal server error"
//	@Router			/api/v1/auth [post]
func (uh *UserHandler) Register(c *gin.Context) {
	var request dto.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		uh.logger.Error("Failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	user := domain.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}

	_, err := uh.svc.Register(c.Request.Context(), &user)
	if err != nil {
		uh.logger.Error("Failed to register user", zap.Error(err))
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	uh.logger.Info("User registered successfully", zap.String("email", user.Email))
	response := util.APIResponse("User created successfully", http.StatusOK, "success", dto.NewUserResponse(&user))
	c.JSON(http.StatusOK, response)
}

// ListUsers godoc
//
//	@Summary		List users
//	@Description	List users with pagination
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200		{object}	util.Response		"Users displayed"
//	@Failure		400		{object}	util.ErrorResponse	"Validation error"
//	@Failure		500		{object}	util.ErrorResponse	"Internal server error"
//	@Router			/api/v1/users [get]
func (uh *UserHandler) ListUsers(c *gin.Context) {
	var request dto.ListUserRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		uh.logger.Error("Invalid query parameters", zap.Error(err))
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	users, err := uh.svc.ListUsers(c.Request.Context(), request.Page, request.Limit)
	if err != nil {
		uh.logger.Error("Failed to list users", zap.Error(err))
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	uh.logger.Info("Users listed successfully")
	response := util.APIResponse("Users displayed successfully", http.StatusOK, "success", users)
	c.JSON(http.StatusOK, response)
}

// GetUser godoc
//
//	@Summary		Get a user
//	@Description	Get a user by ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		uint64				true	"User ID"
//	@Success		200	{object}	dto.UserResponse	"User displayed"
//	@Failure		404	{object}	util.ErrorResponse	"Data not found error"
//	@Failure		500	{object}	util.ErrorResponse	"Internal server error"
//	@Router			/api/v1/users/{id} [get]
//	@Security		BearerAuth
func (uh *UserHandler) GetUser(c *gin.Context) {
	var request dto.GetUserRequest
	if err := c.ShouldBindUri(&request); err != nil {
		uh.logger.Error("Invalid request parameters", zap.Error(err))
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	user, err := uh.svc.GetUser(c.Request.Context(), request.ID)
	if err != nil {
		uh.logger.Error("User not found", zap.Error(err))
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	uh.logger.Info("User retrieved successfully", zap.Uint64("user_id", user.ID))
	response := util.APIResponse("User retrieved successfully", http.StatusOK, "success", dto.NewUserResponse(user))
	c.JSON(http.StatusOK, response)
}

// DeleteUser godoc
//
//	@Summary		Delete a user
//	@Description	Delete a user by ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		uint64				true	"User ID"
//	@Success		200	{object}	util.Response		"User deleted"
//	@Failure		404	{object}	util.ErrorResponse	"Data not found error"
//	@Failure		500	{object}	util.ErrorResponse	"Internal server error"
//	@Router			/api/v1/users/{id} [delete]
//	@Security		BearerAuth
func (uh *UserHandler) DeleteUser(c *gin.Context) {
	var request dto.DeleteUserRequest
	if err := c.ShouldBindUri(&request); err != nil {
		uh.logger.Error("Invalid request parameters", zap.Error(err))
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	err := uh.svc.DeleteUser(c.Request.Context(), request.ID)
	if err != nil {
		uh.logger.Error("Failed to delete user", zap.Error(err))
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	uh.logger.Info("User deleted successfully", zap.Uint64("user_id", request.ID))
	response := util.APIResponse("User deleted successfully", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

// UpdateUser godoc
//
//	@Summary		Update an existing user
//	@Description	Update user details like name, email, password, and role
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"User ID"
//	@Param			request	body		dto.UpdateUserRequest	true	"User update request payload"
//	@Success		200		{object}	util.Response	"User updated successfully"
//	@Failure		400		{object}	util.ErrorResponse	"Invalid request parameters"
//	@Failure		401		{object}	util.ErrorResponse	"Unauthorized error"
//	@Failure		404		{object}	util.ErrorResponse	"User not found"
//	@Failure		500		{object}	util.ErrorResponse	"Internal server error"
//	@Router			/api/v1/users/{id} [put]
func (uh *UserHandler) UpdateUser(c *gin.Context) {
	var request dto.UpdateUserRequest
	if err := c.ShouldBindUri(&request); err != nil {
		uh.logger.Error("Invalid request parameters", zap.Error(err))
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		uh.logger.Error("Failed to bind request JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	idStr := c.Param("id")
	id, err := util.StringToUint64(idStr)
	if err != nil {
		uh.logger.Error("Invalid user ID in request", zap.Error(err))
		response := util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	user := domain.User{
		ID:       id,
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
		Role:     request.Role,
	}

	resp, err := uh.svc.UpdateUser(c.Request.Context(), &user)
	if err != nil {
		uh.logger.Error("Failed to update user", zap.Error(err))
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	uh.logger.Info("User updated successfully", zap.Uint64("user_id", resp.ID))
	response := util.APIResponse("User updated successfully", http.StatusOK, "success", resp)
	c.JSON(http.StatusOK, response)
}

func (uh *UserHandler) GetProfile(c *gin.Context) {
	userSess := util.GetAuthPayload(c, consts.AuthorizationKey)

	profile, err := uh.svc.GetProfile(c.Request.Context(), uint64(userSess.UserID))
	if err != nil {
		uh.logger.Error("Failed to retrieve user", zap.Error(err))
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	uh.logger.Info("User retrieved successfully", zap.Uint64("user_id", profile.ID))
	response := util.APIResponse("User retrieved successfully", http.StatusOK, "success", profile)
	c.JSON(http.StatusOK, response)
}
