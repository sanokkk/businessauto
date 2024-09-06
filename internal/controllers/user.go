package controllers

import (
	"github.com/gin-gonic/gin"
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
func (r *HttpHandler) GetMyUser(c *gin.Context) {
	userId := c.GetString("uid")

	parsedId, err := uuid.Parse(userId)
	if err != nil {
		RespondWithError(c, 401, "Ошибка парсинга JWT-токена", err)

		return
	}

	user, err := r.authService.GetUser(parsedId)
	if err != nil {
		RespondWithError(c, 401, "Ошибка получения информации о пользователе", err)

		return
	}

	c.JSON(200, user)
}
