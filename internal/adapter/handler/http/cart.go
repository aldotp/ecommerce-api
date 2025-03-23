package http

import (
	"context"
	"net/http"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/dto"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/helper"
	"github.com/aldotp/ecommerce-go-api/internal/core/service"
	"github.com/aldotp/ecommerce-go-api/pkg/consts"
	"github.com/aldotp/ecommerce-go-api/pkg/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CartHandler handles cart-related operations
type CartHandler struct {
	CartService *service.CartService
	logger      *zap.Logger
}

// NewCartHandler creates a new CartHandler instance
func NewCartHandler(cartService *service.CartService, logger *zap.Logger) *CartHandler {
	return &CartHandler{
		CartService: cartService,
		logger:      logger,
	}
}

// AddToCart godoc
//
//	@Summary		Add product to cart
//	@Description	Add a product to the user's shopping cart
//	@Tags			Cart
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.AddCartRequest	true	"Add to cart request body"
//	@Success		200		{object}	util.Response	"Product added to cart successfully"
//	@Failure		400		{object}	util.ErrorResponse	"Bad request (validation error)"
//	@Failure		401		{object}	util.ErrorResponse	"Unauthorized error"
//	@Failure		500		{object}	util.ErrorResponse	"Internal server error"
//	@Router			/api/v1/carts [post]
func (h *CartHandler) AddToCart(c *gin.Context) {
	userSess := util.GetAuthPayload(c, consts.AuthorizationKey)

	var request dto.AddCartRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("Failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Adding product to cart",
		zap.Int("user_id", userSess.UserID),
		zap.Int("product_id", request.ProductID),
		zap.Int("quantity", request.Quantity),
	)

	if err := h.CartService.AddToCart(c.Request.Context(), userSess.UserID, request.ProductID, request.Quantity); err != nil {
		h.logger.Error("Failed to add product to cart",
			zap.Int("user_id", userSess.UserID),
			zap.Int("product_id", request.ProductID),
			zap.Int("quantity", request.Quantity),
			zap.Error(err),
		)
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	h.logger.Info("Product added to cart successfully", zap.Int("user_id", userSess.UserID))
	response := util.APIResponse("Product added to cart successfully", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

// ViewCart godoc
//
//	@Summary		View user's cart
//	@Description	Get the list of products in the user's shopping cart
//	@Tags			Cart
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200		{object}	util.Response	"Cart retrieved successfully"
//	@Failure		401		{object}	util.ErrorResponse	"Unauthorized error"
//	@Failure		500		{object}	util.ErrorResponse	"Internal server error"
//	@Router			/api/v1/carts [get]
func (h *CartHandler) ViewCart(c *gin.Context) {
	userSess := util.GetAuthPayload(c, consts.AuthorizationKey)

	h.logger.Info("Retrieving cart", zap.Int("user_id", userSess.UserID))

	cart, err := h.CartService.GetCart(context.Background(), userSess.UserID)
	if err != nil {
		h.logger.Error("Failed to retrieve cart", zap.Int("user_id", userSess.UserID), zap.Error(err))
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	h.logger.Info("Cart retrieved successfully", zap.Int("user_id", userSess.UserID))
	response := util.APIResponse("Get Cart successfully", http.StatusOK, "success", cart)
	c.JSON(http.StatusOK, response)
}

func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	userSess := util.GetAuthPayload(c, consts.AuthorizationKey)

	var request dto.RemoveCartRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("Failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Removing product from cart",
		zap.Int("user_id", userSess.UserID),
		zap.Int("product_id", request.ProductID),
	)

	if err := h.CartService.RemoveFromCart(c.Request.Context(), userSess.UserID, request.ProductID); err != nil {
		h.logger.Error("Failed to remove product from cart",
			zap.Int("user_id", userSess.UserID),
			zap.Int("product_id", request.ProductID),
			zap.Error(err),
		)
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	h.logger.Info("Product removed from cart successfully", zap.Int("user_id", userSess.UserID))
	response := util.APIResponse("Product removed from cart successfully", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *CartHandler) UpdateCart(c *gin.Context) {
	userSess := util.GetAuthPayload(c, consts.AuthorizationKey)

	var request dto.UpdateCartRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("Failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Updating cart",
		zap.Int("user_id", userSess.UserID),
		zap.Int("product_id", request.ProductID),
		zap.Int("quantity", request.Quantity),
	)

	if err := h.CartService.UpdateCart(c.Request.Context(), userSess.UserID, request); err != nil {
		h.logger.Error("Failed to update cart",
			zap.Int("user_id", userSess.UserID),
			zap.Int("product_id", request.ProductID),
			zap.Int("quantity", request.Quantity),
			zap.Error(err),
		)
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	h.logger.Info("Cart updated successfully", zap.Int("user_id", userSess.UserID))
	response := util.APIResponse("Cart updated successfully", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}
