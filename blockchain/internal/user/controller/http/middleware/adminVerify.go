package middleware

import (
	"fmt"
	"github.com/evrone/go-clean-template/internal/user/usecase"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"
)

func AdminVerify(SecretKey string, _ usecase.UserUseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenHeader := ctx.Request.Header.Get("Authorization")
		if tokenHeader == "" {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(SecretKey), nil
		})
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if exp, ok := claims["exp"].(float64); ok {
				expirationTime := time.Unix(int64(exp), 0)
				if time.Now().After(expirationTime) {
					ctx.AbortWithStatus(http.StatusUnauthorized)
					return
				}
			} else {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		} else {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if err != nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		if role != "admin" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}

		ctx.Next()
	}
}
