package usecases

import (
	"context"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
)

// UserUseCase defines a use case interface with methods related to user managing
type UserUseCase interface {
	// SaveUserProfilePicture saves the profile picture of the user with the given ID and
	// returns the updated file path
	SaveUserProfilePicture(
			userID int64,
			image []byte,
	) (string, error)

	// UpdateUserInformation updates the user information.
	UpdateUserInformation(
			ctx context.Context,
			userInformation entities.UserInformation,
	) (*entities.User, error)
}
