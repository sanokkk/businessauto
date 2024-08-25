package middleware

import (
	"autoshop/internal/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckFeatureFlag() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentStorageCfg := config.MustLoadConfig().ContentConfig

		if !contentStorageCfg.UseContentStorage {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Действие с контентом еще недоступно")})

			return
		}

		c.Next()
	}

}
