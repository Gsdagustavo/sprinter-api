package repositories

import (
	"context"
	"database/sql"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
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

func (r *userRepository) UpdateUserProfile(
	ctx context.Context,
	accountInformation entities.AccountInformation,
) error {
	const query = `
	UPDATE users
	SET username = ?,
		biography = ?,
	WHERE id = ?
	`

	result, err := r.conn.ExecContext(ctx, query, &accountInformation.Username, &accountInformation.Biography, &accountInformation.ID)
	if err != nil {
		return derr.JoinError("failed to execute query", err)
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return derr.NotFoundError
	}

	return nil
}
func (r *userRepository) GetUserById(ctx context.Context, id int64) (*entities.User, error) {
	const query = `
	SELECT 
    	id, 
    	name, 
    	email,  
    	carbo_coins, 
    	carbon, 
    	traveled_distance 
	FROM users WHERE id = ?
	`

	var user entities.User
	row := r.conn.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CarboCoins,
		&user.Carbon,
		&user.TraveledDistance,
	)

	if err == sql.ErrNoRows {
		return nil, derr.JoinError("user does not exists", err)
	}

	if err != nil {
		return nil, derr.JoinError("failed to scan the rows", err)
	}

	return &user, nil
}
