package storage

import (
	"autoshop/internal/domain/models"
	"autoshop/internal/storage/filters"
	"github.com/google/uuid"
)

type UserStorage interface {
	SaveUser(user *models.User) error
	GetUser(id uuid.UUID) (*models.User, error)
	CheckCredentials(email string, password string) (*models.User, error)
}

type ProductStorage interface {
	//InsertOne(product *models.Product, categories []models.Category) error
	GetById(id string) (*models.Product, error)
	UpdateProduct(id string, updateFunc func(product *models.Product)) error
	GetWithFilter(filter filters.ProductFilter, skip int, take int, order []filters.OrderBy) ([]models.Product, error)
	Get(skip int, take int, order []filters.OrderBy) ([]models.Product, error)
	GetCategories() ([]models.Category, error)
	CreateCategory(category models.Category, productIds []uuid.UUID) (*models.Category, error)
}
