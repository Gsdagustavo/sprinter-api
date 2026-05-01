package repositories

import (
	"context"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
)

// ProductRepository defines methods for managing product data.
type ProductRepository interface {
	// AddNewProduct is trying to add a new product
	AddNewProduct(ctx context.Context, product *entities.Product) (int64, error)

	// DeleteProduct is trying to delete a product
	DeleteProduct(ctx context.Context, id int64) error

	// UpdateProduct is trying to update a product
	UpdateProduct(ctx context.Context, product *entities.Product) error

	// GetProductByID is trying to get one single product
	GetProductByID(ctx context.Context, id int64) (*entities.Product, error)

	// GetProducts returns a paginated list of products that match the given filter.
	GetProducts(
			ctx context.Context,
			filter entities.GeneralFilter,
	) (*entities.PaginatedList[entities.Product], error)
}