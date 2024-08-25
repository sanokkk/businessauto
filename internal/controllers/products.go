package controllers

import (
	"autoshop/internal/service/dto"
	"autoshop/pkg/logging"
	"fmt"
	"github.com/gin-gonic/gin"
)

// @BasePath		/api/products
// @Summary		Получение товаров
// @Description	gets products with pagination and filters
// @Param			GetProductsRequest	body	dto.Request	true	"Получение товаров"
// @Tags			Products
// @Accept			json
// @Produce		json
// @Success		200	{object}	dto.GetProductsDto
// @Router			/api/products/get [post]
func (r *HttpHandler) GetProducts(c *gin.Context) {
	const op = "HttpHandler.GetProducts"
	log := logging.CreateLoggerWithOp(op)

	log.Info("Получил запрос на получение товаров")

	var request dto.Request
	if err := c.BindJSON(&request); err != nil {
		log.Warn(fmt.Sprintf("Ошибка получения запроса: %w", err))

		RespondWithError(c, 400, fmt.Sprintf("Ошибка получения запроса: %w", err), err)
		return
	}

	productsResponse, err := r.productService.GetProducts(&request.Body)
	if err != nil {
		log.Warn(fmt.Sprintf("Ошибка получения товаров: %w", err))

		RespondWithError(c, 400, "Ошибка получения товаров, попробуйте еще раз", err)
		return
	}

	c.JSON(200, productsResponse)
}
