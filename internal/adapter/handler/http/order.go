package http

import (
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

type OrderHandler struct {
	svc    port.OrderService
	logger *zap.Logger
}

// NewOrderHandler initializes a new OrderHandler
func NewOrderHandler(orderSvc port.OrderService, logger *zap.Logger) *OrderHandler {
	return &OrderHandler{
		svc:    orderSvc,
		logger: logger,
	}
}

// GetOrders godoc
//
//	@Summary		Get User Orders
//	@Description	Retrieve a list of orders for the authenticated user
//	@Tags			Orders
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	util.Response	"Orders retrieved successfully"
//	@Failure		401	{object}	util.ErrorResponse	"Unauthorized"
//	@Failure		500	{object}	util.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/orders [get]
//	@Security		BearerAuth
func (h *OrderHandler) GetOrders(c *gin.Context) {
	userSess := util.GetAuthPayload(c, consts.AuthorizationKey)

	h.logger.Info("Fetching orders", zap.String("user_id", fmt.Sprintf("%v", userSess.UserID)))

	resp, err := h.svc.ListOrders(c.Request.Context(), userSess.UserID)
	if err != nil {
		h.logger.Error("Failed to fetch orders", zap.Error(err))
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	response := util.APIResponse("Get Orders successfully", http.StatusOK, "success", resp)
	c.JSON(http.StatusOK, response)
}

// GetOrderDetail godoc
//
//	@Summary		Get Order Detail
//	@Description	Retrieve details of a specific order by ID
//	@Tags			Orders
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Order ID"
//	@Success		200	{object}	util.Response	"Order details retrieved successfully"
//	@Failure		400	{object}	util.ErrorResponse	"Invalid request parameters"
//	@Failure		401	{object}	util.ErrorResponse	"Unauthorized"
//	@Failure		404	{object}	util.ErrorResponse	"Order not found"
//	@Failure		500	{object}	util.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/orders/{id} [get]
//	@Security		BearerAuth
func (h *OrderHandler) GetOrderDetail(c *gin.Context) {
	userSess := util.GetAuthPayload(c, consts.AuthorizationKey)

	h.logger.Info("Fetching order detail", zap.String("user_id", fmt.Sprintf("%v", userSess.UserID)))

	var request dto.OrderRequest
	if err := c.ShouldBindUri(&request); err != nil {
		h.logger.Warn("Invalid request parameters", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.svc.GetOrder(c.Request.Context(), request.ID, userSess.UserID)
	if err != nil {
		h.logger.Error("Failed to fetch order detail", zap.Error(err))
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	if resp == nil {
		h.logger.Warn("Order not found", zap.String("order_id", fmt.Sprintf("%v", request.ID)))
		response := util.APIResponse("order not found", http.StatusNotFound, "error", nil)
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := util.APIResponse("Get Order Detail successfully", http.StatusOK, "success", resp)
	c.JSON(http.StatusOK, response)
}
