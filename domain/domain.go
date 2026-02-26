package domain

import (
	"context"

	"github.com/VitorFranciscoDev/sprinter-api/domain/entities"
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
	GetUserByEmail(
		ctx context.Context,
		email string,
	) (*entities.User, error)

	// CheckCredentials validates the user credentials
	CheckCredentials(
		ctx context.Context,
		credentials entities.UserCredentials,
	) (bool, error)
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
