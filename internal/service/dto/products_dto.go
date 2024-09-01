package dto

import "autoshop/internal/domain/models"

type GetProductsDto struct {
	Products []models.Product `json:"products"`
}

type GetCategoriesDto struct {
	Categories []models.Category `json:"categories"`
}
