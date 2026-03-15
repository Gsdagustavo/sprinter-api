package datastore

import (
	"context"
	"database/sql"
	"time"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
)

// RepositorySettings creates and manages the access for the database
type RepositorySettings interface {
	// Connection returns a database connection
	Connection() *sql.DB

	// Dismount closes all connections with the database
	Dismount() error

	// ServerTime returns the current time on the server
	ServerTime(ctx context.Context) (*time.Time, error)
}

// AuthRepository defines methods for user authentication and registration.
type AuthRepository interface {
	// GetUserByEmail returns the user with the given email
	GetUserByEmail(
		ctx context.Context,
		email string,
	) (*entities.User, error)

	// GetUserByID returns the user with the given ID
	GetUserByID(
		ctx context.Context,
		userID int64,
	) (*entities.User, error)

	// AttemptRegister tries to register a new user with the given credentials
	AttemptRegister(
		ctx context.Context,
		credentials entities.UserCredentials,
	) (int64, error)

	// CheckUserCredentials validates the user credentials
	CheckUserCredentials(
		ctx context.Context,
		credentials entities.UserCredentials,
	) (bool, error)
}

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

// ActivityRepository defines methods for managing activity data.
type ActivityRepository interface {
	// Save is trying to save the activity
	Save(ctx context.Context, activity *entities.Activity) error

	// Get is trying to get an activity information
	Get(ctx context.Context, id int64) (*entities.Activity, error)

	// GetAll is trying to get all activities
	GetAll(ctx context.Context) ([]*entities.Activity, error)
}
