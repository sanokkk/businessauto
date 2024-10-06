package controllers

import (
	"autoshop/internal/service"
	"autoshop/pkg/custom_errors"
	"autoshop/pkg/logging"
	"bytes"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"log/slog"
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
func (h *HttpHandler) UploadFile(c *fiber.Ctx) error {
	const op = "HttpHandler.UploadFile"
	log := logging.CreateLoggerWithOp(op)

	log.Info("Поступил запрос на загрузку файла")

	productId := c.Query("productId")
	if productId == "" {
		log.Warn("Не выбран товар для загрузки контента")
		return RespondWithErrorFiber(c, 400, "Не выбран товар для загрузки контента", errors.New("Не выбран товар для загрузки контента"))
	}

	fileEntity, err := c.FormFile("file")
	if err != nil {
		log.Warn("Ошибка при загрузке файла", slog.String("error", err.Error()))
		return RespondWithErrorFiber(c, 400, "Ошибка при загрузке файла", err)
	}

	file, err := fileEntity.Open()
	if err != nil {
		log.Warn("Ошибка при обработке файла", slog.String("error", err.Error()))
		return RespondWithErrorFiber(c, 400, "Ошибка при обработке файла", err)
	}

	defer func() error {
		if err := file.Close(); err != nil {
			log.Warn("Ошибка при закрытии файла", slog.String("error", err.Error()))
			return RespondWithErrorFiber(c, 400, "Ошибка при закрытии файла", err)
		}
		return nil
	}()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Warn("Ошибка при чтении байтов файла", slog.String("error", err.Error()))
		return RespondWithErrorFiber(c, 400, "Ошибка при чтении байтов файла", err)
	}

	res, err := h.contentService.UploadContent(service.FileInput{Content: fileBytes, Size: fileEntity.Size})
	if err != nil {
		log.Warn("Ошибка при загрузке байтов файла", slog.String("error", err.Error()))
		return RespondWithErrorFiber(c, 400, "Ошибка при загрузке байтов файла", err)
	}

	if err = h.productService.AddContent(productId, *res); err != nil {
		log.Warn("Ошибка при загрузке обновлении контента товара", slog.String("error", err.Error()))
		return RespondWithErrorFiber(c, 400, "Ошибка при загрузке обновлении контента товара", err)
	}

	c.Status(fiber.StatusOK)
	return nil
}

// @BasePath		/api/content
// @Summary		Скачивание контента
// @Description	downloads content from contentStorage
// @Tags			Content
// @Param			contentId	query	string	true	"Content ID"
// @Produce		octet-stream
// @Success		200	{file}	file	"File downloaded successfully"
// @Router			/api/content/ [get]
func (h *HttpHandler) DownloadFile(c *fiber.Ctx) error {
	const op = "HttpHandler.UploadFile"
	log := logging.CreateLoggerWithOp(op)

	contentId := c.Query("contentId")
	log.Info("Поступил запрос на скачивание файла", slog.String("contentId", contentId))

	contentBytes, err := h.contentService.DownloadContent(contentId)
	if err != nil {
		if errors.Is(err, custom_errors.NoFileError) {
			return RespondWithErrorFiber(c, 404, "Файл не существует", err)
		}

		log.Warn("Ошибка при скачивании контента товара", slog.String("error", err.Error()))
		return RespondWithErrorFiber(c, 400, "Ошибка при скачивании контента товара", err)
	}

	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=octet-stream.local_%s", contentId))
	c.Set("Content-Type", "application/octet-stream")
	return c.SendStream(bytes.NewReader(contentBytes), len(contentBytes))
}
