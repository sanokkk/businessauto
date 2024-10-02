package controllers

import (
	"autoshop/internal/service"
	"autoshop/internal/service/dto"
	"autoshop/pkg/custom_errors"
	"autoshop/pkg/logging"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
)

// @BasePath		/api/products
// @Summary		Получение категорий
// @Description	gets categories
// @Tags			Categories
// @Accept			json
// @Produce		json
// @Success		200	{object}	dto.GetCategoriesDto
// @Router			/api/categories [get]
func (r *HttpHandler) GetCategories(c *gin.Context) {
	const op = "HttpHandler.GetCategories"
	log := logging.CreateLoggerWithOp(op)

	log.Info("Получил запрос на получение всех категорий")

	categoriesResponse, err := r.productService.GetCategories()
	if err != nil {
		log.Warn(fmt.Sprintf("Ошибка получения категорий: %w", err))

		RespondWithError(c, 400, "Ошибка получения категорий, попробуйте еще раз", err)
		return
	}

	c.JSON(200, categoriesResponse)
}

// @BasePath		/api/categories
// @Summary		Создание категории
// @Description	creates categories
// @Tags			Categories
// @Accept			json
// @Accept			multipart/form-data
// @Param			Authorization	header	string	true	"Authorization header"
// @Param			file		formData	file	true	"File to upload"
// @Param			title		formData	string	true	"category title"
// @Param			productIds		formData	[]string	true	"category productIds"
// @Produce		json
// @Success		200	{object}	models.Category
// @Router			/api/categories [post]
func (r *HttpHandler) HandleAddCategory(c *gin.Context) {
	const op = "HttpHandler.HandleAddCategory"
	log := logging.CreateLoggerWithOp(op)

	log.Info("Получил запрос на добавление категории")

	log.Debug("Получаю файл из формы")
	formFile, err := c.FormFile("file")
	if err != nil {
		log.Warn("Ошибка получения изображения для категории", slog.String("error", err.Error()))

		RespondWithError(c, 400, "Ошибка получения изображения для категории", err)
		return
	}

	imageId, err := r.uploadImage(c, formFile, log)
	if err != nil {
		RespondWithError(c, 500, "Внутренняя ошибка сервера", err)
		return
	}

	form, err := c.MultipartForm()
	// todo конвертить это в массив строк
	productIdsStr := form.Value["productId"][0]
	var productIds []uuid.UUID
	if err := json.Unmarshal([]byte(productIdsStr), &productIds); err != nil {
		log.Warn("Ошибка при получении запроса", slog.String("error", err.Error()))

		RespondWithError(c, 400, "Ошибка при получении запроса", err)
		return
	}

	var request dto.CreateCategoryDto
	if err := c.ShouldBind(&request); err != nil {
		log.Warn("Ошибка при получении запроса", slog.String("error", err.Error()))

		RespondWithError(c, 400, "Ошибка при получении запроса", err)
		return
	}

	if err := r.validate.Struct(&request); err != nil {
		RespondWithError(c, 400, err.Error(), custom_errors.ValidationError)

		return
	}

	request.ImageId = *imageId
	category, err := r.productService.AddCategory(request)

	if err != nil {
		log.Warn("Ошибка при создании категории", slog.String("error", err.Error()))

		if err != nil {
			RespondWithError(c, 500, "Внутренняя ошибка сервера", err)
			return
		}
	}

	c.JSON(http.StatusCreated, category)

}

func (r *HttpHandler) uploadImage(c *gin.Context, formFile *multipart.FileHeader, log *slog.Logger) (*uuid.UUID, error) {
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

		RespondWithError(c, 500, "Внутренняя ошибка сервера", err)
		return nil, err
	}

	fileInput := service.FileInput{Content: fileBytes, Size: formFile.Size}

	imageId, err := r.contentService.UploadContent(fileInput)
	if err != nil {
		log.Warn("Ошибка загрузки файла", slog.String("error", err.Error()))

		RespondWithError(c, 500, "Внутренняя ошибка сервера", err)
		return nil, err
	}

	return imageId, nil
}
