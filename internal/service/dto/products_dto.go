package dto

import (
	"autoshop/internal/domain/models"
	"github.com/google/uuid"
)

type GetProductsDto struct {
	Products []models.Product `json:"products"`
}

type GetCategoriesDto struct {
	Categories []models.Category `json:"categories"`
}

type CreateCategoryDto struct {
	Title      string `form:"title" validate:"required"`
	ProductIds []uuid.UUID
	ImageId    uuid.UUID `json:"-"`
}
