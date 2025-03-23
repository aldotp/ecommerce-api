package util

import (
	"strconv"

	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
	"github.com/gin-gonic/gin"
)

func StringToUint64(str string) (uint64, error) {
	num, err := strconv.ParseUint(str, 10, 64)

	return num, err
}

func GetAuthPayload(ctx *gin.Context, key string) *domain.TokenPayload {
	return ctx.MustGet(key).(*domain.TokenPayload)
}
