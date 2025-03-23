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

// ProductHandler represents the HTTP handler for product-related requests
type ProductHandler struct {
	svc    port.ProductService
	logger *zap.Logger
}

// NewProductHandler creates a new ProductHandler instance
func NewProductHandler(svc port.ProductService, logger *zap.Logger) *ProductHandler {
	return &ProductHandler{
		svc:    svc,
		logger: logger,
	}
}

// CreateProduct godoc
//
//	@Summary		Create a new product
//	@Description	Add a new product to the system
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.ProductRequest	true	"Product request payload"
//	@Success		200		{object}	util.Response	"Product created successfully"
//	@Failure		400		{object}	util.ErrorResponse	"Validation error"
//	@Failure		500		{object}	util.ErrorResponse	"Internal server error"
//	@Router			/api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var request dto.ProductRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("Invalid request payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	err := h.svc.Store(c.Request.Context(), request)
	if err != nil {
		h.logger.Error("Failed to create product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, util.APIResponse(err.Error(), http.StatusInternalServerError, "error", nil))
		return
	}

	h.logger.Info("Product created successfully", zap.String("name", request.Name))
	c.JSON(http.StatusOK, util.APIResponse("Product created successfully", http.StatusOK, "success", nil))
}

// GetProduct godoc
//
//	@Summary		Get a product by ID
//	@Description	Retrieve details of a specific product
//	@Tags			Products
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Product ID"
//	@Success		200	{object}	util.Response	"Product found"
//	@Failure		400	{object}	util.ErrorResponse	"Bad request"
//	@Failure		500	{object}	util.ErrorResponse	"Internal server error"
//	@Router			/api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	var request dto.GetProductRequest

	if err := c.ShouldBindUri(&request); err != nil {
		h.logger.Warn("Invalid product ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	product, err := h.svc.FindOne(c.Request.Context(), request.ID)
	if err != nil {
		h.logger.Error("Failed to retrieve product", zap.Int("id", request.ID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, util.APIResponse(err.Error(), http.StatusInternalServerError, "error", nil))
		return
	}

	h.logger.Info("Product found", zap.Int("id", request.ID))
	c.JSON(http.StatusOK, util.APIResponse("Product found successfully", http.StatusOK, "success", dto.NewProductResponse(product)))
}

// ListProducts godoc
//
//	@Summary		List all products
//	@Description	Retrieve a list of available products
//	@Tags			Products
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	util.Response	"List of products"
//	@Failure		500	{object}	util.ErrorResponse	"Internal server error"
//	@Router			/api/v1/products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	var request dto.ListProductRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		h.logger.Warn("Invalid request parameters", zap.Error(err))
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	products, err := h.svc.Finds(c.Request.Context(), request)
	if err != nil {
		h.logger.Error("Failed to retrieve product list", zap.Error(err))
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	h.logger.Info("Retrieved product list successfully")
	c.JSON(http.StatusOK, util.APIResponse("List Product successfully", http.StatusOK, "success", dto.NewProductsResponse(products)))
}

// UpdateProduct godoc
//
//	@Summary		Update a product
//	@Description	Modify an existing product
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int					true	"Product ID"
//	@Param			request	body		dto.ProductRequest	true	"Updated product data"
//	@Success		200		{object}	util.Response	"Product updated successfully"
//	@Failure		400		{object}	util.ErrorResponse	"Validation error"
//	@Failure		500		{object}	util.ErrorResponse	"Internal server error"
//	@Router			/api/v1/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var (
		request    dto.ProductRequest
		requestUri dto.ParamProductRequest
	)

	if err := c.ShouldBindUri(&requestUri); err != nil {
		h.logger.Warn("Invalid product ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("Invalid update payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	err := h.svc.Update(c.Request.Context(), requestUri.ID, request)
	if err != nil {
		h.logger.Error("Failed to update product", zap.Int("id", requestUri.ID), zap.Error(err))
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	h.logger.Info("Product updated successfully", zap.Int("id", requestUri.ID))
	c.JSON(http.StatusOK, util.APIResponse("Product updated successfully", http.StatusOK, "success", nil))
}

// DeleteProduct godoc
//
//	@Summary		Delete a product
//	@Description	Remove a product from the system
//	@Tags			Products
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path	int	true	"Product ID"
//	@Success		200	{object}	util.Response	"Product deleted successfully"
//	@Failure		400	{object}	util.ErrorResponse	"Validation error"
//	@Failure		500	{object}	util.ErrorResponse	"Internal server error"
//	@Router			/api/v1/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	var request dto.ParamProductRequest

	if err := c.ShouldBindUri(&request); err != nil {
		h.logger.Warn("Invalid product ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	err := h.svc.Delete(c.Request.Context(), request.ID)
	if err != nil {
		h.logger.Error("Failed to delete product", zap.Int("id", request.ID), zap.Error(err))
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	h.logger.Info("Product deleted successfully", zap.Int("id", request.ID))
	c.JSON(http.StatusOK, util.APIResponse("Product deleted successfully", http.StatusOK, "success", nil))
}
