package dto

type DepositRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type DepositResponse struct {
	Balance float64 `json:"balance"`
}

type TransferRequest struct {
	RecipientID uint64  `json:"recipient_id" binding:"required,gt=0"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
}

type TransferResponse struct {
	From struct {
		UserID  uint64  `json:"user_id"`
		Balance float64 `json:"balance"`
	} `json:"from"`
	To struct {
		UserID  uint64  `json:"user_id"`
		Balance float64 `json:"balance"`
	} `json:"to"`
}

type WithdrawRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
	// Bank    string  `json:"bank" binding:"required"`
	// Account string  `json:"account" binding:"required"`
}

type BalanceResponse struct {
	Balance float64 `json:"balance"`
}
