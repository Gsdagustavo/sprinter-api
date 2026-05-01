package usecases

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
	) (string, error)

	// AttemptRegister attempts to register a new user in the database
	AttemptRegister(
			ctx context.Context,
			credentials entities.UserCredentials,
	) (string, error)

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

	// AttemptCompleteRegistration tries to complete the user registration
	AttemptCompleteRegistration(
			ctx context.Context,
			information entities.UserInformation,
	) error

	// UploadProfileImage uploads a profile image for the user with the given ID
	UploadProfileImage(
			ctx context.Context,
			userID int64,
			image []byte,
	) error
}