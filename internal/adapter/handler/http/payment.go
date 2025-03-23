package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/dto"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/helper"
	"github.com/aldotp/ecommerce-go-api/internal/core/port"
	"github.com/aldotp/ecommerce-go-api/pkg/consts"
	"github.com/aldotp/ecommerce-go-api/pkg/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PaymentHandler struct {
	PaymentService port.PaymentService
	logger         *zap.Logger
}

// NewPaymentHandler initializes a new PaymentHandler
func NewPaymentHandler(paymentService port.PaymentService, logger *zap.Logger) *PaymentHandler {
	return &PaymentHandler{
		PaymentService: paymentService,
		logger:         logger,
	}
}

// Pay godoc
//
//	@Summary		Process Payment
//	@Description	Make a payment for a given order
//	@Tags			Payment
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.PaymentRequest	true	"Payment request payload"
//	@Success		200		{object}	util.Response	"Payment successful"
//	@Failure		400		{object}	util.ErrorResponse	"Bad request, invalid payload"
//	@Failure		500		{object}	util.ErrorResponse	"Internal server error"
//	@Router			/api/v1/payments [post]
//	@Security		BearerAuth
func (h *PaymentHandler) Pay(c *gin.Context) {
	userSess := util.GetAuthPayload(c, consts.AuthorizationKey)
	if userSess.UserID == 0 {
		h.logger.Error("Unauthorized request")
		c.JSON(http.StatusUnauthorized, util.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil))
		return
	}

	var request dto.PaymentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("Invalid request payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Processing payment", zap.String("order_id", fmt.Sprintf("%v", request.OrderID)))

	if err := h.PaymentService.MakePayment(context.Background(), userSess.UserID, request.OrderID); err != nil {
		h.logger.Error("Payment failed", zap.String("order_id", fmt.Sprintf("%v", request.OrderID)), zap.Error(err))
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	h.logger.Info("Payment successful", zap.String("order_id", fmt.Sprintf("%v", request.OrderID)))

	response := util.APIResponse("Payment successful", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}
