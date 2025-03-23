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

type CategoryHandler struct {
	svc    port.CategoryService
	logger *zap.Logger
}

// NewCategoryHandler creates a new CategoryHandler instance
func NewCategoryHandler(svc port.CategoryService, logger *zap.Logger) *CategoryHandler {
	return &CategoryHandler{
		svc:    svc,
		logger: logger,
	}
}

// CreateCategory godoc
//
//	@Summary		Create a new category
//	@Description	Adds a new category to the system
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CategoryRequest	true	"Category request"
//	@Success		200		{object}	util.Response	"Category created successfully"
//	@Failure		400		{object}	util.ErrorResponse	"Invalid request payload"
//	@Failure		500		{object}	util.ErrorResponse	"Internal server error"
//	@Router			/api/v1/categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var request dto.CategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("Failed to bind request JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	h.logger.Info("Creating category", zap.String("name", request.Name))

	if err := h.svc.Store(c.Request.Context(), request); err != nil {
		h.logger.Error("Failed to create category", zap.Error(err))
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	h.logger.Info("Category created successfully", zap.String("name", request.Name))
	c.JSON(http.StatusOK, util.APIResponse("Category created successfully", http.StatusOK, "success", nil))
}

// GetCategory godoc
//
//	@Summary		Get category by ID
//	@Description	Retrieves a category by its ID
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Category ID"
//	@Success		200	{object}	util.Response	"Category found successfully"
//	@Failure		400	{object}	util.ErrorResponse	"Invalid category ID"
//	@Failure		404	{object}	util.ErrorResponse	"Category not found"
//	@Failure		500	{object}	util.ErrorResponse	"Internal server error"
//	@Router			/api/v1/categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	var request dto.GetCategoryRequest
	if err := c.ShouldBindUri(&request); err != nil {
		h.logger.Warn("Invalid category ID in request", zap.Error(err))
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	h.logger.Info("Fetching category", zap.Int("category_id", request.ID))

	category, err := h.svc.FindOne(c.Request.Context(), request.ID)
	if err != nil {
		h.logger.Error("Failed to retrieve category", zap.Int("category_id", request.ID), zap.Error(err))
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	h.logger.Info("Category retrieved successfully", zap.Int("category_id", request.ID))
	c.JSON(http.StatusOK, util.APIResponse("Category found successfully", http.StatusOK, "success", dto.NewCategoryResponse(category)))
}

// ListCategory godoc
//
//	@Summary		List all categories
//	@Description	Fetches all categories available in the system
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	util.Response	"List of categories successfully retrieved"
//	@Failure		500	{object}	util.ErrorResponse	"Internal server error"
//	@Router			/api/v1/categories [get]
func (h *CategoryHandler) ListCategory(c *gin.Context) {
	h.logger.Info("Fetching category list")

	categories, err := h.svc.Finds(c.Request.Context(), dto.ListCategoryRequest{})
	if err != nil {
		h.logger.Error("Failed to fetch category list", zap.Error(err))
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	h.logger.Info("Category list retrieved successfully", zap.Int("count", len(categories)))
	c.JSON(http.StatusOK, util.APIResponse("List Category successfully", http.StatusOK, "success", dto.NewCategoryResponses(categories)))
}

// UpdateCategory godoc
//
//	@Summary		Update category by ID
//	@Description	Updates an existing category by its ID
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int	true	"Category ID"
//	@Param			request	body		dto.CategoryRequest	true	"Updated category data"
//	@Success		200		{object}	util.Response	"Category updated successfully"
//	@Failure		400		{object}	util.ErrorResponse	"Invalid request payload"
//	@Failure		404		{object}	util.ErrorResponse	"Category not found"
//	@Failure		500		{object}	util.ErrorResponse	"Internal server error"
//	@Router			/api/v1/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	var request dto.CategoryRequest
	var requestUri dto.ParamCategoryRequest

	if err := c.ShouldBindUri(&requestUri); err != nil {
		h.logger.Warn("Invalid category ID in request", zap.Error(err))
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("Failed to bind request JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	h.logger.Info("Updating category", zap.Int("category_id", requestUri.ID), zap.String("new_name", request.Name))

	err := h.svc.Update(c.Request.Context(), requestUri.ID, request)
	if err != nil {
		h.logger.Error("Failed to update category", zap.Int("category_id", requestUri.ID), zap.Error(err))
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	h.logger.Info("Category updated successfully", zap.Int("category_id", requestUri.ID))
	c.JSON(http.StatusOK, util.APIResponse("Category updated successfully", http.StatusOK, "success", nil))
}

// DeleteCategory godoc
//
//	@Summary		Delete category by ID
//	@Description	Removes a category from the system by its ID
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Category ID"
//	@Success		200	{object}	util.Response	"Category deleted successfully"
//	@Failure		400	{object}	util.ErrorResponse	"Invalid category ID"
//	@Failure		404	{object}	util.ErrorResponse	"Category not found"
//	@Failure		500	{object}	util.ErrorResponse	"Internal server error"
//	@Router			/api/v1/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	var request dto.ParamCategoryRequest
	if err := c.ShouldBindUri(&request); err != nil {
		h.logger.Warn("Invalid category ID in request", zap.Error(err))
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	h.logger.Info("Deleting category", zap.Int("category_id", request.ID))

	err := h.svc.Delete(c.Request.Context(), request.ID)
	if err != nil {
		h.logger.Error("Failed to delete category", zap.Int("category_id", request.ID), zap.Error(err))
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	h.logger.Info("Category deleted successfully", zap.Int("category_id", request.ID))
	c.JSON(http.StatusOK, util.APIResponse("Category deleted successfully", http.StatusOK, "success", nil))
}
