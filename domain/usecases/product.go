package usecases

import (
	"context"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
)

// ProductUseCase defines a use case interface with methods related to product managing
type ProductUseCase interface {
	// AddNewProduct adds a new product to the database
	AddNewProduct(ctx context.Context, product *entities.Product) (int64, error)

	// DeleteProduct deletes a product from the database
	DeleteProduct(ctx context.Context, id int64) error

	// UpdateProduct updates a product in the database
	UpdateProduct(ctx context.Context, product *entities.Product) error

	// GetProductByID retrieves a single product by its ID from the database. Returns the product or an error if not found.
	GetProductByID(ctx context.Context, id int64) (*entities.Product, error)

	// GetProducts returns a paginated list of products that match the given filter.
	GetProducts(
			ctx context.Context,
			filter entities.GeneralFilter,
	) (*entities.PaginatedList[entities.Product], error)
}