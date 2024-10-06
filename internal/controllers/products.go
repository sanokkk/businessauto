package controllers

import (
	"autoshop/internal/storage/filters"
	"autoshop/pkg/logging"
	"fmt"
	"github.com/gofiber/fiber/v2"
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
func (r *HttpHandler) GetProducts(c *fiber.Ctx) error {
	const op = "HttpHandler.GetProducts"
	log := logging.CreateLoggerWithOp(op)

	log.Info("Получил запрос на получение товаров")

	var request filters.FilterBody
	if err := c.BodyParser(&request); err != nil {
		log.Warn(fmt.Sprintf("Ошибка получения запроса: %s", err))

		return RespondWithErrorFiber(c, 400, fmt.Sprintf("Ошибка получения запроса: %s", err.Error()), err)
	}

	productsResponse, err := r.productService.GetProducts(&request)
	if err != nil {
		log.Warn(fmt.Sprintf("Ошибка получения товаров: %w", err))

		return RespondWithErrorFiber(c, 400, "Ошибка получения товаров, попробуйте еще раз", err)
	}

	return c.JSON(productsResponse)
}
