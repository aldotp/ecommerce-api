package http

import (
	"net/http"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/dto"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/helper"
	"github.com/aldotp/ecommerce-go-api/internal/core/service"
	"github.com/aldotp/ecommerce-go-api/pkg/consts"
	"github.com/aldotp/ecommerce-go-api/pkg/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CheckoutHandler struct {
	CheckoutService *service.CheckoutService
	Logger          *zap.Logger
}

// NewCheckoutHandler creates a new instance of CheckoutHandler
func NewCheckoutHandler(checkoutService *service.CheckoutService, logger *zap.Logger) *CheckoutHandler {
	return &CheckoutHandler{
		CheckoutService: checkoutService,
		Logger:          logger,
	}
}

// Checkout godoc
//
//	@Summary		Checkout a cart
//	@Description	Completes the checkout process for the user’s cart
//	@Tags			Checkout
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CheckoutRequest	true	"Checkout request"
//	@Success		200		{object}	util.Response	"Checkout successful"
//	@Failure		400		{object}	util.ErrorResponse	"Invalid request payload"
//	@Failure		401		{object}	util.ErrorResponse	"Unauthorized error"
//	@Failure		404		{object}	util.ErrorResponse	"Cart not found"
//	@Failure		500		{object}	util.ErrorResponse	"Internal server error"
//	@Router			/api/v1/checkout [post]
//	@Security		BearerAuth
func (h *CheckoutHandler) Checkout(c *gin.Context) {
	userSess := util.GetAuthPayload(c, consts.AuthorizationKey)

	var request dto.CheckoutRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.Logger.Error("Failed to bind JSON request", zap.Error(err))
		response := util.APIResponse("Invalid request payload", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	h.Logger.Info("Processing checkout",
		zap.Int("userID", userSess.UserID),
		zap.String("paymentMethod", request.PaymentMethod),
	)

	if request.PaymentMethod == "" {
		h.Logger.Warn("Payment method is missing", zap.Int("userID", userSess.UserID))
		response := util.APIResponse("Payment method is required", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	resp, err := h.CheckoutService.Checkout(c.Request.Context(), userSess.UserID, request.PaymentMethod)
	if err != nil {
		h.Logger.Error("Checkout failed",
			zap.Int("userID", userSess.UserID),
			zap.String("paymentMethod", request.PaymentMethod),
			zap.Error(err),
		)
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	h.Logger.Info("Checkout successful",
		zap.Int("userID", userSess.UserID),
		zap.String("paymentMethod", request.PaymentMethod),
	)

	response := util.APIResponse("Checkout successful", http.StatusOK, "success", resp)
	c.JSON(http.StatusOK, response)
}
