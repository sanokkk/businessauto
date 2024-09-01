package storage

import (
	"autoshop/internal/domain/models"
	"autoshop/internal/storage/filters"
)

type UserStorage interface {
	SaveUser(user *models.User) error
	GetUser(id string) (*models.User, error)
	CheckCredentials(email string, password string) (*models.User, error)
}

type ProductStorage interface {
	//InsertOne(product *models.Product, categories []models.Category) error
	//GetById(id string) (*models.Product, error)
	GetWithFilter(filter filters.ProductFilter, skip int, take int, order []filters.OrderBy) ([]models.Product, error)
	Get(skip int, take int, order []filters.OrderBy) ([]models.Product, error)
	GetCategories() ([]models.Category, error)
}
