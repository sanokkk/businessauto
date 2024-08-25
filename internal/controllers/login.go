package controllers

import (
	"autoshop/internal/service"
	"autoshop/pkg/custom_errors"
	"github.com/gin-gonic/gin"
)

func (h *HttpHandler) Login(c *gin.Context) {
	var input service.LoginInput

	if err := c.BindJSON(&input); err != nil {
		RespondWithError(c, 400, "Ошибка при вводе данных", err)

		return
	}

	if err := h.validate.Struct(&input); err != nil {
		RespondWithError(c, 400, err.Error(), custom_errors.ValidationError)

		return
	}

	tokenResponse, err := h.authService.Login(input)

	if err != nil {
		RespondWithError(c, 400, err.Error(), custom_errors.InternalError)

		return
	}

	c.JSON(200, tokenResponse)
}
