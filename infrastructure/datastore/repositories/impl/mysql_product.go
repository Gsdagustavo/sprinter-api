package repositories

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/datastore"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/datastore/repositories"
)

func NewProductRepository(settings repositories.SettingsRepository) repositories.ProductRepository {
	return &productRepository{
		conn:     settings.Connection(),
		settings: settings,
	}
}

type productRepository struct {
	conn     *sql.DB
	settings repositories.SettingsRepository
}

//go:embed _query/product/addNewProduct.sql
var addNewProductQuery string

//go:embed _query/product/deleteProduct.sql
var deleteProductQuery string

//go:embed _query/product/updateProduct.sql
var updateProductQuery string

//go:embed _query/product/getProductById.sql
var getProductByIdQuery string

//go:embed _query/product/getProducts.sql
var getProductsQuery string

func (r productRepository) AddNewProduct(ctx context.Context, product *entities.Product) (int64, error) {
	res, err := r.conn.ExecContext(
		ctx,
		addNewProductQuery,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.ImageURL,
	)
	if err != nil {
		return -1, derr.JoinError("failed to execute query", err)
	}

	productID, err := res.LastInsertId()
	if err != nil {
		return -1, derr.JoinError("failed to get last inserted ID", err)
	}

	return productID, nil
}

func (r productRepository) DeleteProduct(ctx context.Context, id int64) error {
	res, err := r.conn.ExecContext(ctx, deleteProductQuery, id)
	if err != nil {
		return derr.JoinError("failed to execute query", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return derr.JoinError("failed to get rows affected", err)
	}

	if rowsAffected == 0 {
		return derr.NotFoundError
	}

	return nil
}

func (r productRepository) UpdateProduct(ctx context.Context, product *entities.Product) error {
	result, err := r.conn.ExecContext(ctx, updateProductQuery, &product.Name, &product.Description, &product.ImageURL, &product.Price, &product.Stock, product.ID)
	if err != nil {
		return derr.JoinError("failed to execute query", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return derr.JoinError("failed to get rows affected", err)
	}

	if rowsAffected == 0 {
		return derr.NotFoundError
	}

	return nil
}

func (r productRepository) GetProductByID(ctx context.Context, id int64) (*entities.Product, error) {
	var product entities.Product
	err := r.conn.QueryRowContext(ctx, getProductByIdQuery, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.ImageURL,
	)
	if err != nil {
		return nil, derr.JoinError("failed to query or scan", err)
	}

	return &product, nil
}

func (r productRepository) GetProducts(
	ctx context.Context,
	filter entities.GeneralFilter,
) (*entities.PaginatedList[entities.Product], error) {
	ordination := filter.Ordination
	switch filter.OrderBy {
	case "name":
		getProductsQuery += " ORDER BY name " + ordination
	case "price":
		getProductsQuery += " ORDER BY price " + ordination
	case "stock":
		getProductsQuery += " ORDER BY stock " + ordination
	}

	getProductsQuery = datastore.GetPaginated(getProductsQuery, filter)

	rows, err := r.conn.QueryContext(ctx, getProductsQuery)
	if err != nil {
		return nil, derr.JoinError("failed to query", err)
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
			return nil, derr.JoinError("failed to scan", err)
		}

		products = append(products, product)
	}

	countQuery := datastore.GetQueryCount(getProductsQuery)

	var totalCount int64
	err = r.conn.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, derr.JoinError("failed to query or scan count", err)
	}

	pages := datastore.GetTotalPages(totalCount, filter)

	return &entities.PaginatedList[entities.Product]{
		Items:          products,
		RequestedItems: filter.Limit,
		TotalCount:     totalCount,
		Pages:          pages,
	}, nil
}
