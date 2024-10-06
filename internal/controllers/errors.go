package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Code             uint16 `json:"code"`
	ErrorText        string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func RespondWithError(c *gin.Context, code int, errorDescription string, err error) {
	errorResp := ErrorResponse{Code: uint16(code), ErrorDescription: errorDescription, ErrorText: err.Error()}

	c.JSON(code, errorResp)
}

func RespondWithErrorFiber(c *fiber.Ctx, code int, errorDescription string, err error) error {
	errorResp := ErrorResponse{Code: uint16(code), ErrorDescription: errorDescription, ErrorText: err.Error()}

	return c.Status(code).JSON(errorResp)
}
