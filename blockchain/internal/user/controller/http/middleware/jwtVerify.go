package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func JwtVerify(SecretKey string) gin.HandlerFunc {
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
		var claims jwt.MapClaims
		var ok bool
		if claims, ok = token.Claims.(jwt.MapClaims); ok && token.Valid {
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
		id := fmt.Sprintf("%v", claims["user_id"])
		ctx.Set("user_id", id)

		ctx.Next()
	}
}
