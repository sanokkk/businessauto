package controllers

import (
	"autoshop/internal/service"
	"autoshop/pkg/custom_errors"
	"errors"
	"github.com/gin-gonic/gin"
)

func (r *HttpHandler) Register(c *gin.Context) {
	var input service.RegisterInput

	if err := c.BindJSON(&input); err != nil {
		RespondWithError(c, 400, "Ошибка при вводе данных", err)

		return
	}

	if err := r.validate.Struct(&input); err != nil {
		RespondWithError(c, 400, err.Error(), custom_errors.ValidationError)

		return
	}

	response, err := r.authService.Register(input)
	if err != nil {
		if errors.Is(err, custom_errors.AuthenticationError) {
			RespondWithError(c, 500, err.Error(), custom_errors.AuthenticationError)

			return
		}
		RespondWithError(c, 500, err.Error(), custom_errors.RegistrationError)

		return
	}

	c.JSON(201, response)
}
