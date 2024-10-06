package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// @BasePath		/api/users
// @Summary		Получение информации о пользователе
// @Description	gets user information
// @Param			Authorization	header	string	true	"JWT"
// @Tags			Users
// @Accept			json
// @Produce		json
// @Success		200	{object}	dto.GetUserResponse
// @Router			/api/users/ [get]
func (r *HttpHandler) GetMyUser(c *fiber.Ctx) error {
	userId, exists := c.Locals("uid").(string)
	if !exists {

	}

	parsedId, err := uuid.Parse(userId)
	if err != nil {
		RespondWithErrorFiber(c, 401, "Ошибка парсинга JWT-токена", err)

		return err
	}

	user, err := r.authService.GetUser(parsedId)
	if err != nil {
		RespondWithErrorFiber(c, 401, "Ошибка получения информации о пользователе", err)

		return err
	}

	_ = c.JSON(user)

	return nil
}
