package usecases

import (
	"context"

	"github.com/Gsdagustavo/sprinter-api/domain"
	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/datastore"
)

func NewProductUseCases(r datastore.ProductRepository) domain.ProductUseCase {
	return productUseCases{repository: r}
}

type productUseCases struct {
	repository datastore.ProductRepository
}

func (p productUseCases) AddNewProduct(ctx context.Context, product *entities.Product) error {
	return nil
}

func (p productUseCases) DeleteProduct(ctx context.Context, id int64) error {
	return nil
}

func (p productUseCases) UpdateProduct(ctx context.Context, product *entities.Product) error {
	return nil
}

func (p productUseCases) GetAllAvailableProducts(ctx context.Context) ([]*entities.Product, error) {
	return nil, nil
}
