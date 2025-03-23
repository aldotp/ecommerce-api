package helper

import (
	"net/http"

	"github.com/aldotp/ecommerce-go-api/pkg/consts"
	"github.com/aldotp/ecommerce-go-api/pkg/util"
)

func ErrorResponse(err error) (int, util.Response) {
	statusCode := http.StatusInternalServerError
	message := "Internal server error"

	switch err {
	case consts.ErrDataNotFound:
		statusCode = http.StatusNotFound
		message = err.Error()
	case consts.ErrNoUpdatedData:
		statusCode = http.StatusNotModified
		message = err.Error()
	case consts.ErrConflictingData, consts.ErrEmailAlreadyExist:
		statusCode = http.StatusConflict
		message = err.Error()
	case consts.ErrInsufficientStock, consts.ErrInsufficientPayment:
		statusCode = http.StatusBadRequest
		message = err.Error()
	case consts.ErrTokenDuration, consts.ErrTokenCreation, consts.ErrInvalidToken, consts.ErrExpiredToken:
		statusCode = http.StatusUnauthorized
		message = err.Error()
	case consts.ErrInvalidCredentials:
		statusCode = http.StatusUnauthorized
		message = err.Error()
	case consts.ErrEmptyAuthorizationHeader, consts.ErrInvalidAuthorizationHeader, consts.ErrInvalidAuthorizationType, consts.ErrEmptyCart, consts.ErrInsufficientBalance, consts.ErrCannotSendBalanceSameAccount:
		statusCode = http.StatusBadRequest
		message = err.Error()
	case consts.ErrUnauthorized:
		statusCode = http.StatusUnauthorized
		message = err.Error()
	case consts.ErrForbidden:
		statusCode = http.StatusForbidden
		message = err.Error()
	case consts.ErrEmailNotVerified:
		statusCode = http.StatusForbidden
		message = err.Error()
	case consts.ErrNotImplemented:
		statusCode = http.StatusNotImplemented
		message = err.Error()
	}

	return statusCode, util.APIResponse(message, statusCode, "error", nil)
}
