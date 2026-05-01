package repositories

import (
	"context"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
)

// UserRepository define methods for managing user data
type UserRepository interface {
	// UpdateUserInformation attempts to edit an user profile with the given information
	UpdateUserInformation(
			ctx context.Context,
			accountInformation entities.UserInformation,
	) error

	// GetUserById attempts to get the user from the given user id
	GetUserById(ctx context.Context, id int64) (*entities.User, error)
}