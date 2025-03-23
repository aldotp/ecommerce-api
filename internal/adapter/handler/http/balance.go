package http

import (
	"net/http"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/dto"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/helper"
	"github.com/aldotp/ecommerce-go-api/internal/core/port"
	"github.com/aldotp/ecommerce-go-api/pkg/consts"
	"github.com/aldotp/ecommerce-go-api/pkg/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// BalanceHandler represents the HTTP handler for balance-related operations
type BalanceHandler struct {
	svc    port.BalanceService
	logger *zap.Logger
}

// NewBalanceHandler creates a new BalanceHandler instance
func NewBalanceHandler(svc port.BalanceService, logger *zap.Logger) *BalanceHandler {
	return &BalanceHandler{
		svc:    svc,
		logger: logger,
	}
}

// Deposit godoc
//
// @Summary	Deposit money to user balance
// @Description	Add funds to a user's balance
// @Tags		Balance
// @Accept	json
// @Produce	json
// @Security	BearerAuth
// @Param	request	body	dto.DepositRequest	true	"Deposit request payload"
// @Success	200	{object}	util.Response	"Deposit successful"
// @Failure	400	{object}	util.ErrorResponse	"Validation error"
// @Failure	500	{object}	util.ErrorResponse	"Internal server error"
// @Router	/api/v1/balance/deposit [post]
func (bh *BalanceHandler) Deposit(c *gin.Context) {
	userSess := util.GetAuthPayload(c, consts.AuthorizationKey) // Ambil user ID dari context
	if userSess.UserID == 0 {
		bh.logger.Error("Unauthorized request")
		c.JSON(http.StatusUnauthorized, util.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil))
		return
	}

	var request dto.DepositRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		bh.logger.Error("Failed to bind deposit request", zap.Error(err))
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	balance, err := bh.svc.Deposit(c.Request.Context(), uint64(userSess.UserID), request.Amount)
	if err != nil {
		bh.logger.Error("Deposit failed", zap.Error(err))
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	bh.logger.Info("Deposit successful", zap.Uint64("user_id", uint64(userSess.UserID)))
	response := util.APIResponse("Deposit successful", http.StatusOK, "success", balance)
	c.JSON(http.StatusOK, response)
}

// Transfer godoc
//
// @Summary	Transfer money to another user
// @Description	Transfer funds from one user to another
// @Tags		Balance
// @Accept	json
// @Produce	json
// @Security	BearerAuth
// @Param	request	body	dto.TransferRequest	true	"Transfer request payload"
// @Success	200	{object}	util.Response	"Transfer successful"
// @Failure	400	{object}	util.ErrorResponse	"Validation error"
// @Failure	500	{object}	util.ErrorResponse	"Internal server error"
// @Router	/api/v1/balance/transfer [post]
func (bh *BalanceHandler) Transfer(c *gin.Context) {
	userSess := util.GetAuthPayload(c, consts.AuthorizationKey) // Ambil user ID dari context

	if userSess.UserID == 0 {
		bh.logger.Error("Unauthorized request")
		c.JSON(http.StatusUnauthorized, util.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil))
		return
	}

	var request dto.TransferRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		bh.logger.Error("Failed to bind transfer request", zap.Error(err))
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	transaction, err := bh.svc.Transfer(c.Request.Context(), uint64(userSess.UserID), request.RecipientID, request.Amount)
	if err != nil {
		bh.logger.Error("Transfer failed", zap.Error(err))
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	bh.logger.Info("Transfer successful", zap.Uint64("from_user_id", uint64(userSess.UserID)), zap.Uint64("to_user_id", request.RecipientID))
	response := util.APIResponse("Transfer successful", http.StatusOK, "success", transaction)
	c.JSON(http.StatusOK, response)
}

// CheckBalance godoc
//
// @Summary	Check user balance
// @Description	Retrieve the current balance of a user
// @Tags		Balance
// @Accept	json
// @Produce	json
// @Security	BearerAuth
// @Param	id	path	uint64	true	"User ID"
// @Success	200	{object}	util.Response	"Balance retrieved"
// @Failure	400	{object}	util.ErrorResponse	"Invalid request parameters"
// @Failure	404	{object}	util.ErrorResponse	"User not found"
// @Failure	500	{object}	util.ErrorResponse	"Internal server error"
// @Router	/api/v1/balance/{id} [get]
func (bh *BalanceHandler) CheckBalance(c *gin.Context) {
	userSess := util.GetAuthPayload(c, consts.AuthorizationKey) // Ambil user ID dari context
	if userSess.UserID == 0 {
		bh.logger.Error("Unauthorized request")
		c.JSON(http.StatusUnauthorized, util.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil))
		return
	}

	balance, err := bh.svc.CheckBalance(c.Request.Context(), uint64(userSess.UserID))
	if err != nil {
		bh.logger.Error("Failed to retrieve balance", zap.Error(err))
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	bh.logger.Info("Balance retrieved successfully", zap.Uint64("user_id", uint64(userSess.UserID)))
	response := util.APIResponse("Balance retrieved successfully", http.StatusOK, "success", balance)
	c.JSON(http.StatusOK, response)
}

// Withdraw godoc
//
//	@Summary		Withdraw balance
//	@Description	Withdraw balance from user account
//	@Tags			Transactions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.WithdrawRequest	true	"Withdraw request"
//	@Success		200		{object}	util.Response		"Withdraw successful"
//	@Failure		400		{object}	util.ErrorResponse	"Invalid request parameters"
//	@Failure		401		{object}	util.ErrorResponse	"Unauthorized error"
//	@Failure		403		{object}	util.ErrorResponse	"Insufficient balance"
//	@Failure		500		{object}	util.ErrorResponse	"Internal server error"
//	@Router			/api/v1/withdraw [post]
func (bh *BalanceHandler) Withdraw(c *gin.Context) {
	var request dto.WithdrawRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		bh.logger.Error("Failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, util.APIResponse(err.Error(), http.StatusBadRequest, "error", nil))
		return
	}

	userSess := util.GetAuthPayload(c, consts.AuthorizationKey)
	if userSess.UserID == 0 {
		bh.logger.Error("Unauthorized request")
		c.JSON(http.StatusUnauthorized, util.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil))
		return
	}

	err := bh.svc.Withdraw(c.Request.Context(), uint64(userSess.UserID), request.Amount)
	if err != nil {
		bh.logger.Error("Failed to withdraw", zap.Error(err))
		statusCode, response := helper.ErrorResponse(err)
		c.JSON(statusCode, response)
		return
	}

	bh.logger.Info("Withdraw successful", zap.Uint64("user_id", uint64(userSess.UserID)), zap.Float64("amount", request.Amount))
	response := util.APIResponse("Withdraw successful", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}
