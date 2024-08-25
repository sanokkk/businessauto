package models

import (
	"autoshop/pkg/uuid_helpers"
	"github.com/google/uuid"
)

type Category struct {
	Id       uuid.UUID
	Title    string
	Products []Product `gorm:"many2many:product_category;"`
}

type Product struct {
	Id         uuid.UUID
	Title      string
	Price      float64
	ImagesIds  uuid_helpers.UUIDArray `gorm:"column:imageIds;type:jsonb"`
	Maker      string
	Categories []Category `gorm:"many2many:product_category;" json:"categories,omitempty"`
}

type ProductCategory struct {
	Id         uuid.UUID
	ProductId  uuid.UUID `sql:"productId" gorm:"column:productId"`
	CategoryId uuid.UUID `sql:"categoryId" gorm:"column:categoryId"`
}
