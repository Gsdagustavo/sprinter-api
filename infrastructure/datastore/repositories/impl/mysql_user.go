package repositories

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/datastore/repositories"
)

func NewUserRepository(
	settings repositories.SettingsRepository,
) repositories.UserRepository {
	return &userRepository{
		conn:     settings.Connection(),
		settings: settings,
	}
}

type userRepository struct {
	conn     *sql.DB
	settings repositories.SettingsRepository
}

//go:embed _query/user/updateUserInformation.sql
var updateUserInformationQuery string

//go:embed _query/user/getUserById.sql
var getUserByIDQuery string

func (r *userRepository) UpdateUserInformation(
	ctx context.Context,
	accountInformation entities.UserInformation,
) error {
	_, err := r.conn.ExecContext(
		ctx,
		updateUserInformationQuery,
		&accountInformation.Username,
		&accountInformation.Biography,
		&accountInformation.ID,
	)
	if err != nil {
		return derr.JoinError("failed to execute query", err)
	}

	return nil
}

func (r *userRepository) GetUserById(ctx context.Context, id int64) (*entities.User, error) {
	var user entities.User
	row := r.conn.QueryRowContext(ctx, getUserByIDQuery, id)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CarboCoins,
		&user.Carbon,
		&user.TraveledDistance,
	)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, derr.NotFoundError
		}

		return nil, derr.JoinError("failed to scan the rows", err)
	}

	return &user, nil
}
