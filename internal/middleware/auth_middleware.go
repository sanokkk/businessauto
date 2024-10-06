package middleware

import (
	"autoshop/pkg/custom_errors"
	"autoshop/pkg/jwt_helper"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
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
		c.Set("role", claims.Role)
		c.Next()
	}
}

func AuthenticateFiber() fiber.Handler {
	return func(c *fiber.Ctx) error {
		headers := c.GetReqHeaders()
		tokens, exists := headers["Authorization"]
		if !exists || tokens[0] == "" {
			return c.Status(http.StatusUnauthorized).JSON(map[string]string{"error": fmt.Sprintf("Нет токена авторизации")})
		}

		token := tokens[0]
		claims, isTokenValid, err := jwt_helper.ValidateToken(token)
		if !isTokenValid {
			if err != nil {
				if errors.Is(err, custom_errors.TokenExpiredError) {
					return c.Status(http.StatusUnauthorized).JSON(map[string]string{"error": fmt.Sprintf("Токен протух")})
				}
			}

			return c.Status(http.StatusUnauthorized).JSON(map[string]string{"error": fmt.Sprintf("Ошибка валидации токена")})
		}

		c.Locals("uid", claims.UserId)
		c.Locals("token", token)
		c.Locals("role", claims.Role)
		return c.Next()
	}
}

func CheckForRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")

		claims, _, err := jwt_helper.ValidateToken(token)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": fmt.Sprintf("Ошибка валидации токена")})
		}

		if claims.Role != role {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": fmt.Sprintf("Нет прав на действие")})
		}

		return c.Next()
	}
}
