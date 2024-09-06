package service

import (
	"autoshop/internal/domain/models"
	"autoshop/internal/service/dto"
	"autoshop/internal/storage"
	"autoshop/internal/storage/filters"
	"autoshop/pkg/custom_errors"
	"autoshop/pkg/logging"
	"fmt"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"log/slog"
)

type ProductsService interface {
	GetProducts(filter *filters.FilterBody) (dto.GetProductsDto, error)
	GetCategories() (*dto.GetCategoriesDto, error)
	AddContent(productId string, uuid uuid.UUID) error
}

type ProductService struct {
	productsStorage storage.ProductStorage
}

func NewProductService(productsStorage storage.ProductStorage) *ProductService {
	return &ProductService{productsStorage: productsStorage}
}

func (s *ProductService) GetProducts(filter *filters.FilterBody) (dto.GetProductsDto, error) {
	const op = "ProductService.GetProducts"
	log := logging.CreateLoggerWithOp(op)

	if filter.Filter == nil {
		return s.processWithoutFilter(filter, log)
	}

	return s.processWithFilter(filter, log)
}

func (s *ProductService) processWithFilter(filter *filters.FilterBody, log *slog.Logger) (dto.GetProductsDto, error) {
	var productFilter filters.ProductFilter
	if err := mapstructure.Decode(filter.Filter, &productFilter); err != nil {
		log.Warn("Ошибка конвертации базового фильтра к продуктовому")

		return dto.GetProductsDto{}, custom_errors.ConvertationError
	}

	//todo find out why i cant parse struct from interface
	/*productFilter, canConvert := filter.Filter.(filters.ProductFilter)
	if !canConvert {
		log.Warn("Ошибка конвертации базового фильтра к продуктовому")

		return nil, custom_errors.ConvertationError
	}*/

	products, err := s.productsStorage.GetWithFilter(productFilter, *filter.Skip, *filter.Take, filter.Order)
	if err != nil {
		log.Warn(fmt.Sprintf("Ошибка получения товаров с фильтрами: %w", err))

		return dto.GetProductsDto{}, err
	}

	return dto.GetProductsDto{Products: products}, nil
}

func (s *ProductService) processWithoutFilter(filter *filters.FilterBody, log *slog.Logger) (dto.GetProductsDto, error) {
	log.Info("Запрашиваю товары без фильтров")

	products, err := s.productsStorage.Get(*filter.Skip, *filter.Take, filter.Order)
	if err != nil {
		log.Warn(fmt.Sprintf("Ошибка получения товаров без фильтров: %w", err))

		return dto.GetProductsDto{}, err
	}

	return dto.GetProductsDto{Products: products}, nil
}

func (s *ProductService) GetCategories() (*dto.GetCategoriesDto, error) {
	const op = "ProductService.GetCategories"
	log := logging.CreateLoggerWithOp(op)

	log.Info("Начинаю получение категорий")

	result, err := s.productsStorage.GetCategories()
	if err != nil {
		log.Warn(fmt.Sprintf("Ошибка получения категорий: %s", err.Error()))

		return nil, err
	}

	return &dto.GetCategoriesDto{Categories: result}, nil
}

func (s *ProductService) AddContent(productId string, contentId uuid.UUID) error {
	const op = "ProductService.AddContent"
	log := logging.CreateLoggerWithOp(op)

	log.Info("Начинаю обновление контента по товару", slog.String("prodId", productId))

	prod, err := s.productsStorage.GetById(productId)
	if err != nil {
		log.Warn(fmt.Sprintf("Ошибка получения товара: %s", err.Error()))

		return err
	}

	prod.ImagesIds = append(prod.ImagesIds, contentId)
	if err = s.productsStorage.UpdateProduct(
		prod.Id.String(),
		func(product *models.Product) {
			product.ImagesIds = prod.ImagesIds
		}); err != nil {
		log.Warn(fmt.Sprintf("Ошибка получения товара: %s", err.Error()))

		return err
	}

	return nil
}
