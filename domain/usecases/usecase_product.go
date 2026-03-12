package usecases

import (
	"context"

	"github.com/Gsdagustavo/sprinter-api/domain"
	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
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
	err := rules.ValidateProduct(product)
	if err != nil {
		return 0, err
	}

	return p.repository.AddNewProduct(ctx, product)
}

func (p productUseCases) DeleteProduct(ctx context.Context, id int64) error {

	return p.repository.DeleteProduct(ctx, id)
}

func (p productUseCases) UpdateProduct(ctx context.Context, product *entities.Product) error {
	err := rules.ValidateProduct(product)
	if err != nil {
		return err
	}

	return p.repository.UpdateProduct(ctx, product)
}

func (p productUseCases) GetProductByID(ctx context.Context, id int64) (*entities.Product, error) {
	return p.repository.GetProductByID(ctx, id)
}

func (p productUseCases) GetProducts(
	ctx context.Context,
	filter entities.GeneralFilter,
) (*entities.PaginatedList[entities.Product], error) {
	products, err := p.repository.GetProducts(ctx, filter)
	if err != nil {
		return nil, derr.JoinInternalError(err, "failed to get products")
	}

	return products, nil
}
