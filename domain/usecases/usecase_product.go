package usecases

import (
	"context"

	"github.com/Gsdagustavo/sprinter-api/domain"
	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/rules"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/datastore"
)

type productUseCases struct {
	repository datastore.ProductRepository
}

func NewProductUseCases(r datastore.ProductRepository) domain.ProductUseCase {
	return productUseCases{repository: r}
}

func (p productUseCases) AddNewProduct(ctx context.Context, product *entities.Product) (int64, error) {
	rules.ValidateProduct(product)

	return p.repository.AddNewProduct(ctx, product)
}

func (p productUseCases) DeleteProduct(ctx context.Context, id int64) error {
	return p.repository.DeleteProduct(ctx, id)
}

func (p productUseCases) UpdateProduct(ctx context.Context, product *entities.Product) error {
	rules.ValidateProduct(product)

	return p.repository.UpdateProduct(ctx, product)
}

func (p productUseCases) GetSingleProduct(ctx context.Context, id int64) (*entities.Product, error) {
	return p.repository.GetSingleProduct(ctx, id)
}

func (p productUseCases) GetAllAvailableProducts(ctx context.Context) ([]*entities.Product, error) {
	return p.repository.GetAllAvailableProducts(ctx)
}
