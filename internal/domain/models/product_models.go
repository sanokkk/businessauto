package models

import (
	"autoshop/pkg/uuid_helpers"
	"github.com/google/uuid"
)

type Category struct {
	Id       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Products []Product `gorm:"many2many:product_category;"`
}

type Product struct {
	Id         uuid.UUID              `json:"id"`
	Title      string                 `json:"title"`
	Price      float64                `json:"price"`
	ImagesIds  uuid_helpers.UUIDArray `gorm:"column:imageIds;type:jsonb" json:"imagesIds"`
	Maker      string                 `json:"maker"`
	Categories []Category             `gorm:"many2many:product_category;" json:"-,omitempty"`
}

type ProductCategory struct {
	Id         uuid.UUID
	ProductId  uuid.UUID `sql:"productId" gorm:"column:productId"`
	CategoryId uuid.UUID `sql:"categoryId" gorm:"column:categoryId"`
}
