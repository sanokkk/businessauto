package middleware

import (
	"autoshop/internal/config"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func CheckFeatureFlag() fiber.Handler {
	return func(c *fiber.Ctx) error {
		contentStorageCfg := config.MustLoadConfig().ContentConfig

		if !contentStorageCfg.UseContentStorage {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Действие с контентом еще недоступно")})
		}

		return c.Next()
	}
}
