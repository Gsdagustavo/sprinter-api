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
                      price,
                      stock,
					  image_url
                      )
	VALUES (?, ?, ?, ?, ?)
	`

	res, err := r.conn.ExecContext(
		ctx,
		query,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.ImageURL,
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
		price = ?,
		stock = ?,
		image_url = ?
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
		   price,
		   stock,
		   image_url
	FROM products
	WHERE id = ? AND status_code = 0
	`

	var product entities.Product
	err := r.conn.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.ImageURL,
	)
	if err != nil {
		return nil, derr.JoinInternalError(err, "failed to query or scan")
	}

	return &product, nil
}

func (r productRepository) GetProducts(
	ctx context.Context,
	filter entities.GeneralFilter,
) (*entities.PaginatedList[entities.Product], error) {
	query := `
	SELECT
		id,
		name,
		description,
		price,
		stock,
		image_url
	FROM products
	WHERE status_code = 0
	`

	ordination := filter.Ordination
	switch filter.OrderBy {
	case "name":
		query += " ORDER BY name " + ordination
	case "price":
		query += " ORDER BY price " + ordination
	case "stock":
		query += " ORDER BY stock " + ordination
	}

	query = datastore.GetPaginated(query, filter)

	rows, err := r.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, derr.JoinInternalError(err, "failed to query")
	}
	defer rows.Close()

	products := make([]entities.Product, 0)
	for rows.Next() {
		var product entities.Product
		err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.ImageURL,
		)
		if err != nil {
			return nil, derr.JoinInternalError(err, "failed to scan")
		}

		products = append(products, product)
	}

	countQuery := datastore.GetQueryCount(query)

	var totalCount int64
	err = r.conn.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, derr.JoinInternalError(err, "failed to query or scan count")
	}

	pages := datastore.GetTotalPages(totalCount, filter)

	return &entities.PaginatedList[entities.Product]{
		Items:          products,
		RequestedItems: filter.Limit,
		TotalCount:     totalCount,
		Pages:          pages,
	}, nil
}
