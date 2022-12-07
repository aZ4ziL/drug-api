package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/aZ4ziL/drug-api/auth"
	"github.com/gin-gonic/gin"
)

func AuthenticationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.Request.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "no_token",
				"message": "Please login with token header",
			})
			return
		}

		tokenJWT := strings.Replace(authorizationHeader, "Bearer ", "", -1)
		claims, err := auth.ReadAndVerifyToken(tokenJWT)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "invalid_token",
				"message": "Your token is invalid or expired",
			})
			return
		}

		newContext := context.WithValue(context.Background(), "userInfo", claims)
		ctx.Request.WithContext(newContext)
		ctx.Next()
	}
}
