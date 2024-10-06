package controllers

import (
	"autoshop/internal/service"
	"autoshop/pkg/custom_errors"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// @BasePath		/api/users
// @Summary		Регистрация нового пользователя
// @Description	register the user and returns tokens
// @Param			RegisterData	body	service.RegisterInput	true	"Регистрация нового пользователя"
// @Tags			Users
// @Accept			json
// @Produce		json
// @Success		201	{object}	dto.TokenResponse
// @Router			/api/users/register [post]
func (r *HttpHandler) Register(c *fiber.Ctx) error {
	var input service.RegisterInput

	if err := c.BodyParser(&input); err != nil {
		return RespondWithErrorFiber(c, 400, "Ошибка при вводе данных", err)
	}

	if err := r.validate.Struct(&input); err != nil {
		return RespondWithErrorFiber(c, 400, err.Error(), custom_errors.ValidationError)
	}

	response, err := r.authService.Register(input)
	if err != nil {
		if errors.Is(err, custom_errors.AuthenticationError) {
			return RespondWithErrorFiber(c, 500, err.Error(), custom_errors.AuthenticationError)
		}
		return RespondWithErrorFiber(c, 500, err.Error(), custom_errors.RegistrationError)
	}

	return c.Status(http.StatusCreated).JSON(response)
}
