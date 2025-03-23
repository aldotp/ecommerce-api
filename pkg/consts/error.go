package consts

import (
	"errors"
	"net/http"
)

var (
	ErrInternal                     = errors.New("internal server error")
	ErrDataNotFound                 = errors.New("data not found")
	ErrNoUpdatedData                = errors.New("no data to update")
	ErrConflictingData              = errors.New("data already exists")
	ErrEmailAlreadyExist            = errors.New("email already exist")
	ErrInsufficientStock            = errors.New("product stock is not enough")
	ErrInsufficientPayment          = errors.New("total paid is less than total price")
	ErrTokenDuration                = errors.New("invalid token duration format")
	ErrTokenCreation                = errors.New("error creating token")
	ErrExpiredToken                 = errors.New("access token has expired")
	ErrInvalidToken                 = errors.New("access token is invalid")
	ErrInvalidCredentials           = errors.New("invalid email or password")
	ErrEmptyAuthorizationHeader     = errors.New("authorization header is not provided")
	ErrInvalidAuthorizationHeader   = errors.New("authorization header format is invalid")
	ErrInvalidAuthorizationType     = errors.New("authorization type is not supported")
	ErrUnauthorized                 = errors.New("user is unauthorized to access the resource")
	ErrForbidden                    = errors.New("user is forbidden to access the resource")
	ErrEmailNotVerified             = errors.New("email is not verified")
	ErrInvalidSignature             = errors.New("invalid signature")
	ErrNotImplemented               = errors.New("not implemented")
	ErrEmptyCart                    = errors.New("cart is empty")
	ErrInsufficientBalance          = errors.New("insufficient balance")
	ErrCannotSendBalanceSameAccount = errors.New("cannot send balance to the same account")
)

var ErrorToHTTPStatusCode = map[error]int{
	ErrInternal:                   http.StatusInternalServerError,
	ErrDataNotFound:               http.StatusNotFound,
	ErrConflictingData:            http.StatusConflict,
	ErrInvalidCredentials:         http.StatusUnauthorized,
	ErrUnauthorized:               http.StatusUnauthorized,
	ErrEmptyAuthorizationHeader:   http.StatusUnauthorized,
	ErrInvalidAuthorizationHeader: http.StatusUnauthorized,
	ErrInvalidAuthorizationType:   http.StatusUnauthorized,
	ErrInvalidToken:               http.StatusUnauthorized,
	ErrExpiredToken:               http.StatusUnauthorized,
	ErrForbidden:                  http.StatusForbidden,
	ErrNoUpdatedData:              http.StatusBadRequest,
	ErrInsufficientStock:          http.StatusBadRequest,
	ErrInsufficientPayment:        http.StatusBadRequest,
	ErrTokenCreation:              http.StatusInternalServerError,
	ErrTokenDuration:              http.StatusInternalServerError,
}
