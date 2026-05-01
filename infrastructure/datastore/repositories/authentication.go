package repositories

import (
	"context"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
)

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

	// AttemptCompleteRegistration tries to complete the user registration
	AttemptCompleteRegistration(
			ctx context.Context,
			information entities.UserInformation,
	) error
}