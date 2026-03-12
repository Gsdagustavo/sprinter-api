package repositories

import (
	"context"
	"database/sql"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
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

func (r productRepository) AddNewProduct(ctx context.Context, product *entities.Product) (int64, error) {
	const query = `
	INSERT INTO products (
					  name,
                      description,
                      image,
                      price,
                      quantity,
                      discount)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	res, err := r.conn.ExecContext(
		ctx,
		query,
		&product.Name,
		&product.Description,
		&product.ImageURL,
		&product.Price,
		&product.Stock,
	)
	if err != nil {
		return -1, derr.JoinInternalError(err, "failed to execute query")
	}

	productID, err := res.LastInsertId()
	if err != nil {
		return -1, derr.JoinInternalError(err, "failed to get last inserted ID")
	}

	return productID, nil
}

func (r productRepository) DeleteProduct(ctx context.Context, id int64) error {
	const query = `UPDATE products SET status_code = 1 WHERE id = ?`

	res, err := r.conn.ExecContext(ctx, query, id)
	if err != nil {
		return derr.JoinInternalError(err, "failed to execute query")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return derr.JoinInternalError(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return derr.NotFoundError
	}

	return nil
}

func (r productRepository) UpdateProduct(ctx context.Context, product *entities.Product) error {
	const query = `
	UPDATE products 
	SET name = ?, 
		description = ?, 
		image = ?, 
		price = ?, 
		quantity = ?, 
		discount = ? 
	WHERE id = ?
	`

	result, err := r.conn.ExecContext(ctx, query, &product.Name, &product.Description, &product.ImageURL, &product.Price, &product.Stock, product.ID)
	if err != nil {
		return derr.JoinInternalError(err, "failed to execute query")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return derr.JoinInternalError(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return derr.NotFoundError
	}

	return nil
}

func (r productRepository) GetProductByID(ctx context.Context, id int64) (*entities.Product, error) {
	const query = `
	SELECT id, 
		   name, 
		   description, 
		   image, 
		   price, 
		   quantity, 
		   discount 
	FROM products 
	WHERE id = ? AND status_code = 0
	`

	var product entities.Product
	err := r.conn.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.ImageURL,
		&product.Price,
		&product.Stock,
	)
	if err != nil {
		return nil, derr.JoinInternalError(err, "failed to query or scan")
	}

	return &product, nil
}

func (r productRepository) GetAllProducts(ctx context.Context) ([]entities.Product, error) {
	const query = `
	SELECT 
		id, 
		name, 
		description, 
		image, 
		price, 
		quantity, 
		discount 
	FROM products 
	WHERE status_code = 0
	`
	rows, err := r.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, derr.JoinInternalError(err, "failed to query")
	}
	defer rows.Close()

	var products []entities.Product
	for rows.Next() {
		var product entities.Product
		err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.ImageURL,
			&product.Price,
			&product.Stock,
		)
		if err != nil {
			return nil, derr.JoinInternalError(err, "failed to scan")
		}

		products = append(products, product)
	}

	return products, nil
}
