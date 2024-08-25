package controllers

import (
	"autoshop/pkg/logging"
	"github.com/gin-gonic/gin"
	"log/slog"
)

func (h *HttpHandler) UploadFile(c *gin.Context) {
	const op = "HttpHandler.UploadFile"
	log := logging.CreateLoggerWithOp(op)

	log.Info("Поступил запрос на загрузку файла")

	fileEntity, err := c.FormFile("file")
	if err != nil {
		log.Warn("Ошибка при загрузке файла", slog.String("error", err.Error()))
		RespondWithError(c, 400, "Ошибка при загрузке файла", err)
		return
	}

	file, err := fileEntity.Open()
	if err != nil {
		log.Warn("Ошибка при обработке файла", slog.String("error", err.Error()))
		RespondWithError(c, 400, "Ошибка при обработке файла", err)
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Warn("Ошибка при закрытии файла", slog.String("error", err.Error()))
			RespondWithError(c, 400, "Ошибка при закрытии файла", err)
			return
		}
	}()

	var fileBytes []byte
	_, err = file.Read(fileBytes)
	if err != nil {
		log.Warn("Ошибка при чтении байтов файла", slog.String("error", err.Error()))
		RespondWithError(c, 400, "Ошибка при чтении байтов файла", err)
		return
	}

	res, err := h.contentService.UploadContent(fileBytes)
	if err != nil {
		log.Warn("Ошибка при загрузке байтов файла", slog.String("error", err.Error()))
		RespondWithError(c, 400, "Ошибка при загрузке байтов файла", err)
		return
	}

	c.JSON(200, res)
}
