package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
	"github.com/Gsdagustavo/sprinter-api/domain/util"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/datastore"
)

type authenticationRepository struct {
	conn     *sql.DB
	settings datastore.RepositorySettings
}

func NewAuthenticationRepository(settings datastore.RepositorySettings) datastore.AuthRepository {
	return &authenticationRepository{
		conn:     settings.Connection(),
		settings: settings,
	}
}

func (r authenticationRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*entities.User, error) {
	const query = `
	SELECT 
    	id, 
    	name, 
    	email,  
    	carbo_coins, 
    	carbon, 
    	traveled_distance 
	FROM users WHERE email = ?
	`

	var user entities.User
	row := r.conn.QueryRowContext(ctx, query, email)
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

		return nil, derr.JoinInternalError(err, "failed to scan")
	}

	return &user, nil
}

func (r authenticationRepository) AttemptRegister(
	ctx context.Context,
	credentials entities.UserCredentials,
) (int64, error) {
	const query = `INSERT INTO users (name, email, password) VALUES (?, ?, ?)`

	result, err := r.conn.ExecContext(
		ctx,
		query,
		credentials.Name,
		credentials.Email,
		credentials.Password,
	)
	if err != nil {
		return -1, derr.JoinInternalError(err, "failed to execute query")
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return -1, derr.JoinInternalError(err, "failed to get last inserted ID")
	}

	return userID, nil
}

func (r authenticationRepository) CheckUserCredentials(
	ctx context.Context,
	credentials entities.UserCredentials,
) (bool, error) {
	query := `
	SELECT password 
	FROM users
	WHERE email = ?
	`

	var password string
	err := r.conn.QueryRowContext(ctx, query, credentials.Email).Scan(&password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, derr.NotFoundError
		}

		return false, derr.JoinInternalError(err, "failed to query or scan")
	}

	valid := util.CheckValidPassword(credentials.Password, password)
	return valid, nil
}

func (r authenticationRepository) GetUserByID(
	ctx context.Context,
	userID int64,
) (*entities.User, error) {
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
	row := r.conn.QueryRowContext(ctx, query, userID)
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

		return nil, derr.JoinInternalError(err, "failed to scan")
	}

	return &user, nil
}
