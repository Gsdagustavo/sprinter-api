package repositories

import (
	"context"
	"database/sql"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/datastore"
)

func NewUserRepository(
	settings datastore.RepositorySettings,
) datastore.UserRepository {
	return &userRepository{
		conn:     settings.Connection(),
		settings: settings,
	}
}

type userRepository struct {
	conn     *sql.DB
	settings datastore.RepositorySettings
}

func (ur *userRepository) EditUserProfile(
	ctx context.Context,
	editIt entities.EditUserProfileDTO,
) error {
	return nil
}
