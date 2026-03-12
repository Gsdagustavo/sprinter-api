package domain

import (
	"context"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
)

// AuthenticationUseCase defines a use case interface with methods related to user authentication
type AuthenticationUseCase interface {
	// AttemptLogin attempts to log in with the given credentials
	AttemptLogin(
		ctx context.Context,
		credentials entities.UserCredentials,
	) (*entities.AuthenticationResponse, error)

	// AttemptRegister attempts to register a new user in the database
	AttemptRegister(
		ctx context.Context,
		credentials entities.UserCredentials,
	) (*entities.AuthenticationResponse, error)

	// GetUserByEmail returns the user with the given email
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)

	// CheckCredentials validates the user credentials
	CheckCredentials(
		ctx context.Context,
		credentials entities.UserCredentials,
	) (bool, error)

	// GetUserByToken returns the user with the given authentication token
	GetUserByToken(
		ctx context.Context,
		token string,
	) (*entities.User, error)
}

// UserUseCase defines a use case interface with methods related to user managing
type UserUseCase interface {
	// SaveUserProfilePicture saves the profile picture of the user with the given ID and
	// returns the updated file path
	SaveUserProfilePicture(
		userID int64,
		image []byte,
	) (string, error)
}

// ProductUseCase defines a use case interface with methods related to product managing
type ProductUseCase interface {
	// AddNewProduct adds a new product to the database
	AddNewProduct(
		ctx context.Context,
		product *entities.Product,
	) (int64, error)

	// DeleteProduct deletes a product from the database
	DeleteProduct(ctx context.Context, id int64) error

	// UpdateProduct updates a product in the database
	UpdateProduct(ctx context.Context, product *entities.Product) error

	// GetProductByID retrieves a single product by its ID from the database. Returns the product or an error if not found.
	GetProductByID(ctx context.Context, id int64) (*entities.Product, error)

	// GetAllProducts returns all available products from the database
	GetAllProducts(ctx context.Context) ([]entities.Product, error)
}
