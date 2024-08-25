package service

import (
	"autoshop/internal/domain/models"
	"autoshop/internal/storage"
	"autoshop/internal/storage/filters"
	"autoshop/pkg/custom_errors"
	"autoshop/pkg/logging"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"log/slog"
)

type ProductsService interface {
	GetProducts(filter *filters.FilterBody) ([]models.Product, error)
}

type ProductService struct {
	productsStorage storage.ProductStorage
}

func NewProductService(productsStorage storage.ProductStorage) *ProductService {
	return &ProductService{productsStorage: productsStorage}
}

func (s *ProductService) GetProducts(filter *filters.FilterBody) ([]models.Product, error) {
	const op = "ProductService.GetProducts"
	log := logging.CreateLoggerWithOp(op)

	if filter.Filter == nil {
		return s.processWithoutFilter(filter, log)
	}

	return s.processWithFilter(filter, log)
}

func (s *ProductService) processWithFilter(filter *filters.FilterBody, log *slog.Logger) ([]models.Product, error) {
	var productFilter filters.ProductFilter
	if err := mapstructure.Decode(filter.Filter, &productFilter); err != nil {
		log.Warn("Ошибка конвертации базового фильтра к продуктовому")

		return nil, custom_errors.ConvertationError
	}

	//todo find out why i cant parse struct from interface
	/*productFilter, canConvert := filter.Filter.(filters.ProductFilter)
	if !canConvert {
		log.Warn("Ошибка конвертации базового фильтра к продуктовому")

		return nil, custom_errors.ConvertationError
	}*/

	products, err := s.productsStorage.GetWithFilter(productFilter, filter.Skip, filter.Take, filter.Order)
	if err != nil {
		log.Warn(fmt.Sprintf("Ошибка получения товаров с фильтрами: %w", err))

		return nil, err
	}

	return products, nil
}

func (s *ProductService) processWithoutFilter(filter *filters.FilterBody, log *slog.Logger) ([]models.Product, error) {
	log.Info("Запрашиваю товары без фильтров")

	products, err := s.productsStorage.Get(filter.Skip, filter.Take, filter.Order)
	if err != nil {
		log.Warn(fmt.Sprintf("Ошибка получения товаров без фильтров: %w", err))

		return nil, err
	}

	return products, nil
}
