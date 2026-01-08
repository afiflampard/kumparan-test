package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/afif-musyayyidin/hertz-boilerplate/config"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v5"
)

func GetJwtSecret() []byte {
	return []byte(config.LoadConfig().JWTSecret)
}

type Claims struct {
    AuthorID string `json:"author_id"`
    AuthorName string `json:"author_name"`
    AuthorEmail string `json:"author_email"`
    jwt.RegisteredClaims
}

func GenerateToken(authorID, authorName, authorEmail string, expireDuration time.Duration) (string, error) {
    claims := Claims{
        AuthorID: authorID,
        AuthorName: authorName,
        AuthorEmail: authorEmail,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(GetJwtSecret())
}

func AuthMiddleware() app.HandlerFunc {
    return func(c context.Context, ctx *app.RequestContext) {
        authHeader := string(ctx.GetHeader("Authorization"))
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing or invalid token"})
            ctx.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
            return GetJwtSecret(), nil
        })

        if err != nil || !token.Valid {
            ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
            ctx.Abort()
            return
        }

        claims, ok := token.Claims.(*Claims)
        if !ok {
            ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
            ctx.Abort()
            return
        }

        // Save claims to context for next handlers
        ctx.Set("author_id", claims.AuthorID)
        ctx.Set("author_name", claims.AuthorName)
        ctx.Set("author_email", claims.AuthorEmail)

        ctx.Next(c)
    }
}
