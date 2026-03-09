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
	// Add is trying to add a new product
	Add(ctx context.Context, product *entities.Product) error

	// Delete is trying to delete a product
	Delete(ctx context.Context, id int64) error

	// Update is trying to update a product
	Update(ctx context.Context, product *entities.Product) error

	// Get is trying to get one single product
	Get(ctx context.Context, id int64) (*entities.Product, error)

	// GetAll is trying to get all products
	GetAll(ctx context.Context) ([]*entities.Product, error)
}

// ActivityRepository defines methods for managing activity data.
type ActivityRepository interface {

	// StartActivity creates a new activity with the given type and returns its unique identifier or an error if it fails.
	StartActivity(ctx context.Context, activityType entities.ActivityType) (int64, error)

	// EndActivity terminates an active activity by its unique identifier and updates its status in the repository.
	EndActivity(ctx context.Context, activityID int64) error
}
