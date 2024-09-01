package controllers

import (
	"autoshop/pkg/logging"
	"fmt"
	"github.com/gin-gonic/gin"
)

// @BasePath		/api/products
// @Summary		Получение категорий
// @Description	gets categories
// @Tags			Categories
// @Accept			json
// @Produce		json
// @Success		200	{object}	dto.GetCategoriesDto
// @Router			/api/categories [get]
func (r *HttpHandler) GetCategories(c *gin.Context) {
	const op = "HttpHandler.GetCategories"
	log := logging.CreateLoggerWithOp(op)

	log.Info("Получил запрос на получение всех категорий")

	categoriesResponse, err := r.productService.GetCategories()
	if err != nil {
		log.Warn(fmt.Sprintf("Ошибка получения категорий: %w", err))

		RespondWithError(c, 400, "Ошибка получения категорий, попробуйте еще раз", err)
		return
	}

	c.JSON(200, categoriesResponse)
}
