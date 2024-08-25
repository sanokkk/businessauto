package dto

import "autoshop/internal/domain/models"

type GetProductsDto struct {
	Products []models.Product `json:"products"`
}
