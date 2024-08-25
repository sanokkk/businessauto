package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
