package storage

import (
	"autoshop/internal/config"
	"autoshop/internal/domain/models"
	"autoshop/internal/storage/filters"
	"autoshop/pkg/logging"
	"errors"
	"fmt"
	"github.com/google/uuid"
	errorswrap "github.com/hashicorp/errwrap"
	"github.com/iancoleman/strcase"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log/slog"
	"reflect"
	"strings"
)

type ProductStore struct {
	db *gorm.DB
}

func NewProductStore(cfg *config.DbConfig) *ProductStore {
	db, err := gorm.Open(sqlite.Open(cfg.DbConnectionString), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic(err)
	}
	return &ProductStore{db: db}
}

func (s *ProductStore) AddProduct(product models.Product, category string) (*models.Product, error) {
	const op = "ProductStore.AddProduct"
	log := logging.CreateLoggerWithOp(op)

	var dbCategory models.Category

	res := s.db.Table("categories").Where("title = ?", category).First(&dbCategory)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Error("Не нашел категорий")

			return nil, ErrNotFound
		}

		err := errorswrap.Wrap(errors.New("Ошибка получения категорий"), res.Error)
		log.Error(err.Error())

		return nil, err
	}

	res = s.db.Table("products").Create(&product)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Error("Не нашел категорий")

			return nil, ErrNotFound
		}

		err := errorswrap.Wrap(errors.New("Ошибка получения категорий"), res.Error)
		log.Error(err.Error())

		return nil, err
	}

	var prodCategory models.ProductCategory = models.ProductCategory{
		Id:         uuid.New(),
		ProductId:  product.Id,
		CategoryId: dbCategory.Id}
	res = s.db.Table("product_category").Create(&prodCategory)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Error("Не нашел категорий")

			return nil, ErrNotFound
		}

		err := errorswrap.Wrap(errors.New("Ошибка получения категорий"), res.Error)
		log.Error(err.Error())

		return nil, err
	}

	return &product, nil
}

func (s *ProductStore) GetWithFilter(filter filters.ProductFilter, skip int, take int, order []filters.OrderBy) ([]models.Product, error) {
	const op = "ProductStore.GetWithFilter"
	log := logging.CreateLoggerWithOp(op)

	tx := s.db.Table("products")

	if filter.CategoryFilter != nil {
		tx = s.getTxByCategories(filter, log, tx)
	}

	if filter.PriceFilter != nil {
		tx = tx.Where("products.price >= ? and products.price <= ?", filter.PriceFilter.Min, filter.PriceFilter.Max)
	}

	if filter.MakerFilter != nil {
		tx = tx.Where("products.maker IN (?)", filter.MakerFilter.Makers)
	}

	if filter.TitleFilter != nil {
		tx = tx.Where("lower(products.title) like ?", "%"+strings.ToLower(filter.TitleFilter.Title)+"%")
	}

	orderExpr := getOrderByExpression(order)
	log.Debug(fmt.Sprintf("Сформировал OrderBy %s", orderExpr))

	var products []models.Product

	res := tx.
		Offset(skip).
		Limit(take).
		Order(orderExpr).
		Find(&products)

	if res.Error != nil {
		err := errorswrap.Wrap(errors.New("Ошибка получения продуктов"), res.Error)
		log.Warn(err.Error())

		return nil, res.Error
	}

	return products, nil
}

func (s *ProductStore) getTxByCategories(filter filters.ProductFilter, log *slog.Logger, tx *gorm.DB) *gorm.DB {
	var categoriesIds []string

	catRes := s.db.
		Table("categories").
		Where("title IN (?)", filter.CategoryFilter.Categories).
		Select("id").
		Find(&categoriesIds)

	if catRes.Error != nil {
		err := errorswrap.Wrap(errors.New("Ошибка получения категорий для товаров"), catRes.Error)
		log.Warn(err.Error())

		return tx
	}

	return tx.
		Distinct().
		InnerJoins("JOIN product_category ON products.id = product_category.productId").
		Where("product_category.categoryId in (?)", categoriesIds)
}

func (s *ProductStore) Get(skip int, take int, order []filters.OrderBy) ([]models.Product, error) {
	const op = "ProductStore.Get"
	log := logging.CreateLoggerWithOp(op)

	orderExpr := getOrderByExpression(order)
	log.Debug("Сформировал OrderBy ", orderExpr)

	var products []models.Product

	res := s.db.Offset(skip).Limit(take).Order(orderExpr).Find(&products)
	if res.Error != nil {
		err := errorswrap.Wrap(errors.New("Ошибка получения продуктов"), res.Error)
		log.Warn(err.Error())

		return nil, res.Error
	}

	return products, nil
}

func (s *ProductStore) GetCategories() ([]models.Category, error) {
	const op = "ProductStore.GetCategories"
	log := logging.CreateLoggerWithOp(op)

	var result []models.Category

	tx := s.db.Table("categories").Find(&result)
	if tx.Error != nil {
		err := errorswrap.Wrap(errors.New("Ошибка получения категорий"), tx.Error)
		log.Warn(err.Error())

		return nil, tx.Error
	}

	return result, nil
}

func getOrderByExpression(orders []filters.OrderBy) string {
	var sb strings.Builder

	prodType := reflect.TypeOf(models.Product{})

	for _, orderBy := range orders {
		if _, hasField := prodType.FieldByName(orderBy.Field); !hasField {
			continue
		}

		switch orderBy.Desc {
		case false:
			sb.WriteString(strcase.ToLowerCamel(orderBy.Field) + ",")
			break
		case true:
			sb.WriteString(fmt.Sprintf("%s desc,", strcase.ToLowerCamel(orderBy.Field)))
			break
		}
	}

	return strings.TrimRight(sb.String(), ",")
}
