package repositories

import (
	"context"
	"database/sql"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/datastore"
)

func NewProductRepository(settings datastore.RepositorySettings) datastore.ProductRepository {
	return &productRepository{
		conn:     settings.Connection(),
		settings: settings,
	}
}

type productRepository struct {
	conn     *sql.DB
	settings datastore.RepositorySettings
}

func (r productRepository) AddNewProduct(ctx context.Context, product *entities.Product) error {
	const query = `
		
	`
}
