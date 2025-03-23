package middleware

import (
	"net/http"

	"strings"

	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
	"github.com/aldotp/ecommerce-go-api/internal/core/port"
	"github.com/aldotp/ecommerce-go-api/pkg/consts"
	"github.com/aldotp/ecommerce-go-api/pkg/util"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey = "authorization"
	authorizationType      = "bearer"
)

func AuthMiddleware(token port.TokenInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		isEmpty := len(authorizationHeader) == 0
		if isEmpty {
			err := consts.ErrEmptyAuthorizationHeader
			response := util.APIResponse(err.Error(), http.StatusUnauthorized, "error", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		fields := strings.Fields(authorizationHeader)
		isValid := len(fields) == 2
		if !isValid {
			err := consts.ErrInvalidAuthorizationHeader
			response := util.APIResponse(err.Error(), http.StatusUnauthorized, "error", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		currentAuthorizationType := strings.ToLower(fields[0])
		if currentAuthorizationType != authorizationType {
			err := consts.ErrInvalidAuthorizationType
			response := util.APIResponse(err.Error(), http.StatusUnauthorized, "error", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		accessToken := fields[1]
		payload, err := token.VerifyAccessToken(accessToken)
		if err != nil {
			response := util.APIResponse(err.Error(), http.StatusUnauthorized, "error", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		ctx.Set(consts.AuthorizationKey, payload)
		ctx.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload := util.GetAuthPayload(ctx, consts.AuthorizationKey)

		isAdmin := payload.Role == domain.Admin
		if !isAdmin {
			err := consts.ErrForbidden
			response := util.APIResponse(err.Error(), http.StatusForbidden, "error", nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		ctx.Next()
	}
}
