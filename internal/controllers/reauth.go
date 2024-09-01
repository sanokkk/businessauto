package controllers

import (
	"autoshop/pkg/logging"
	"github.com/gin-gonic/gin"
)

// @BasePath		/api/users
// @Summary		Обновление токена
// @Description	login the user and returns tokens
// @Param Authorization header string true "Рефреш"
// @Tags			Users
// @Accept			json
// @Produce		json
// @Success		200	{object} dto.ReauthResponse
// @Router			/api/users/reauth [get]
func (r *HttpHandler) Reauth(c *gin.Context) {
	const op = "HttpHandler.Reauth"
	log := logging.CreateLoggerWithOp(op)

	log.Info("Поступил запрос на регенерацию токена")

	refresh := c.GetHeader("Authorization")
	tokenResponse, err := r.authService.Reauth(refresh)
	if err != nil {
		RespondWithError(c, 401, err.Error(), err)

		return
	}

	c.JSON(200, tokenResponse)
}
