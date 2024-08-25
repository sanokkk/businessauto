package controllers

import (
	"autoshop/pkg/logging"
	"github.com/gin-gonic/gin"
)

func (r *HttpHandler) Reauth(c *gin.Context) {
	const op = "HttpHandler.Reauth"
	log := logging.CreateLoggerWithOp(op)

	log.Info("Поступил запрос на регенерацию токена")

	refresh := c.GetString("refresh")
	token, err := r.authService.Reauth(refresh)
	if err != nil {
		RespondWithError(c, 401, err.Error(), err)

		return
	}

	c.JSON(200, struct {
		Token string `json:"token"`
	}{Token: token})
}
