package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/afif-musyayyidin/hertz-boilerplate/config"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == nil || !strings.HasPrefix(string(authHeader), "Bearer ") {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(string(authHeader), "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.LoadConfig().JWTSecret), nil
		})

		if err != nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Next(c)
	}
}
