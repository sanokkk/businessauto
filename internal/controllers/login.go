package controllers

import (
	"autoshop/internal/service"
	"autoshop/pkg/custom_errors"
	"github.com/gofiber/fiber/v2"
)

// @BasePath		/api/users
// @Summary		Аутентификация пользователя
// @Description	login the user and returns tokens
// @Param			LoginData	body	service.LoginInput	true	"Аутентификация пользователя"
// @Tags			Users
// @Accept			json
// @Produce		json
// @Success		200	{object}	dto.TokenResponse
// @Router			/api/users/login [post]
func (h *HttpHandler) Login(c *fiber.Ctx) error {
	var input service.LoginInput

	if err := c.BodyParser(&input); err != nil {
		return RespondWithErrorFiber(c, 400, "Ошибка при вводе данных", err)
	}

	if err := h.validate.Struct(&input); err != nil {
		return RespondWithErrorFiber(c, 400, err.Error(), custom_errors.ValidationError)
	}

	tokenResponse, err := h.authService.Login(input)

	if err != nil {
		return RespondWithErrorFiber(c, 400, err.Error(), custom_errors.InternalError)
	}

	return c.JSON(tokenResponse)
}
