package controllers

import (
	"autoshop/internal/service"
	"autoshop/internal/service/dto"
	"autoshop/pkg/custom_errors"
	"autoshop/pkg/helpers"
	"autoshop/pkg/logging"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"mime/multipart"
	"strings"
)

// @BasePath		/api/products
// @Summary		Получение категорий
// @Description	gets categories
// @Tags			Categories
// @Accept			json
// @Produce		json
// @Success		200	{object}	dto.GetCategoriesDto
// @Router			/api/categories [get]
func (r *HttpHandler) GetCategories(c *fiber.Ctx) error {
	const op = "HttpHandler.GetCategories"
	log := logging.CreateLoggerWithOp(op)

	log.Info("Получил запрос на получение всех категорий")

	categoriesResponse, err := r.productService.GetCategories()
	if err != nil {
		log.Warn(fmt.Sprintf("Ошибка получения категорий: %w", err))

		return RespondWithErrorFiber(c, 400, "Ошибка получения категорий, попробуйте еще раз", err)
	}

	return c.JSON(categoriesResponse)
}

// @BasePath		/api/categories
// @Summary		Создание категории
// @Description	creates categories
// @Tags			Categories
// @Accept			json
// @Accept			multipart/form-data
// @Param			Authorization	header		string		true	"Authorization header"
// @Param			file			formData	file		true	"File to upload"
// @Param			title			formData	string		true	"category title"
// @Param			productIds		formData	[]string	false	"category productIds"
// @Produce		json
// @Success		200	{object}	models.Category
// @Router			/api/categories [post]
func (r *HttpHandler) HandleAddCategory(c *fiber.Ctx) error {
	const op = "HttpHandler.HandleAddCategory"
	log := logging.CreateLoggerWithOp(op)

	log.Info("Получил запрос на добавление категории")

	log.Debug("Получаю файл из формы")
	formFile, err := c.FormFile("file")
	if err != nil {
		log.Warn("Ошибка получения изображения для категории", slog.String("error", err.Error()))

		return RespondWithErrorFiber(c, 400, "Ошибка получения изображения для категории", err)
	}

	imageId, err := r.uploadImage(formFile, log)
	if err != nil {
		return RespondWithErrorFiber(c, 500, "Внутренняя ошибка сервера", err)
	}

	productIdsStr := strings.Split(c.FormValue("productIds"), ",")
	productIds, err := helpers.ConvertStringArray(productIdsStr)
	if err != nil {
		log.Warn("Ошибка при получении запроса", slog.String("error", err.Error()))

		return RespondWithErrorFiber(c, 400, "Ошибка при получении запроса", err)
	}

	var request dto.CreateCategoryDto
	request.Title = c.FormValue("title")

	if err := r.validate.Struct(&request); err != nil {
		return RespondWithErrorFiber(c, 400, err.Error(), custom_errors.ValidationError)
	}

	request.ProductIds = productIds
	request.ImageId = *imageId
	category, err := r.productService.AddCategory(request)

	if err != nil {
		log.Warn("Ошибка при создании категории", slog.String("error", err.Error()))

		if err != nil {
			return RespondWithErrorFiber(c, 500, "Внутренняя ошибка сервера", err)
		}
	}

	return c.JSON(category)
}

func (r *HttpHandler) uploadImage(formFile *multipart.FileHeader, log *slog.Logger) (*uuid.UUID, error) {
	file, err := formFile.Open()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Warn("Ошибка при закрытии файла", slog.String("error", err.Error()))
		}
	}()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Warn("Ошибка получения байтов файла", slog.String("error", err.Error()))

		return nil, err
	}

	fileInput := service.FileInput{Content: fileBytes, Size: formFile.Size}

	imageId, err := r.contentService.UploadContent(fileInput)
	if err != nil {
		log.Warn("Ошибка загрузки файла", slog.String("error", err.Error()))

		return nil, err
	}

	return imageId, nil
}
