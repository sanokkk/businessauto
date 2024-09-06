package controllers

import (
	"autoshop/internal/service"
	"autoshop/pkg/custom_errors"
	"autoshop/pkg/logging"
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"net/http"
)

// @BasePath		/api/content
// @Summary		Загрузка контента
// @Description	uploads content for contentStorage
// @Tags			Content
// @Accept			multipart/form-data
// @Param			file		formData	file	true	"File to upload"
// @Param			productId	query		string	true	"Product ID"
// @Produce		json
// @Success		200	{string}	string
// @Router			/api/content/ [post]
func (h *HttpHandler) UploadFile(c *gin.Context) {
	const op = "HttpHandler.UploadFile"
	log := logging.CreateLoggerWithOp(op)

	log.Info("Поступил запрос на загрузку файла")

	productId := c.Query("productId")
	if productId == "" {
		log.Warn("Не выбран товар для загрузки контента")
		RespondWithError(c, 400, "Не выбран товар для загрузки контента", errors.New("Не выбран товар для загрузки контента"))
		return
	}

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

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Warn("Ошибка при чтении байтов файла", slog.String("error", err.Error()))
		RespondWithError(c, 400, "Ошибка при чтении байтов файла", err)
		return
	}

	res, err := h.contentService.UploadContent(service.FileInput{Content: fileBytes, Size: fileEntity.Size})
	if err != nil {
		log.Warn("Ошибка при загрузке байтов файла", slog.String("error", err.Error()))
		RespondWithError(c, 400, "Ошибка при загрузке байтов файла", err)
		return
	}

	if err = h.productService.AddContent(productId, *res); err != nil {
		log.Warn("Ошибка при загрузке обновлении контента товара", slog.String("error", err.Error()))
		RespondWithError(c, 400, "Ошибка при загрузке обновлении контента товара", err)
		return
	}

	c.Status(200)
}

// @BasePath		/api/content
// @Summary		Скачивание контента
// @Description	downloads content from contentStorage
// @Tags			Content
// @Param			contentId	query	string	true	"Content ID"
// @Produce		octet-stream
// @Success		200	{file}	file	"File downloaded successfully"
// @Router			/api/content/ [get]
func (h *HttpHandler) DownloadFile(c *gin.Context) {
	const op = "HttpHandler.UploadFile"
	log := logging.CreateLoggerWithOp(op)

	contentId := c.Query("contentId")
	log.Info("Поступил запрос на скачивание файла", slog.String("contentId", contentId))

	contentBytes, err := h.contentService.DownloadContent(contentId)
	if err != nil {
		if errors.Is(err, custom_errors.NoFileError) {
			RespondWithError(c, 404, "Файл не существует", err)

			return
		}

		log.Warn("Ошибка при скачивании контента товара", slog.String("error", err.Error()))
		RespondWithError(c, 400, "Ошибка при скачивании контента товара", err)

		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=octet-stream.local_%s", contentId))
	c.Header("Content-Type", "application/octet-stream")
	c.DataFromReader(http.StatusOK, int64(len(contentBytes)), "application/octet-stream", bytes.NewReader(contentBytes), nil)
}
