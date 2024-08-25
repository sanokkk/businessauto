package middleware

import (
	"autoshop/pkg/custom_errors"
	"autoshop/pkg/jwt_helper"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Нет токена авторизации")})
			c.Abort()
			return
		}

		claims, isTokenValid, err := jwt_helper.ValidateToken(token)
		if !isTokenValid {
			if err != nil {
				if errors.Is(err, custom_errors.TokenExpiredError) {
					c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Токен протух")})
					c.Abort()
					return
				}
			}

			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Ошибка валидации токена")})
			c.Abort()
			return
		}

		c.Set("uid", claims.UserId)
		c.Set("token", token)
		c.Next()
	}
}
