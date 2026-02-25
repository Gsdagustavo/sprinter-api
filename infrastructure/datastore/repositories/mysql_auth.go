package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/VitorFranciscoDev/sprinter-api/domain/entities"
	"github.com/VitorFranciscoDev/sprinter-api/domain/entities/derr"
	"github.com/VitorFranciscoDev/sprinter-api/infrastructure/datastore"
	"golang.org/x/crypto/bcrypt"
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

func (r authenticationRepository) CheckValidPassword(
	ctx context.Context,
	userID int64,
	password string,
) (bool, error) {
	const query = `SELECT id, password FROM users WHERE id = ?`

	var passwordHash string
	err := r.conn.QueryRowContext(ctx, query, userID).Scan(&userID, &passwordHash)
	if err != nil {
		return false, derr.NewInternalError("failed to scan password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err != nil, nil
}

func (r authenticationRepository) AttemptRegister(
	ctx context.Context,
	credentials entities.UserCredentials,
) (int64, error) {
	const query = `INSERT INTO users (name, email, password) VALUES (?, ?, ?, ?)`

	result, err := r.conn.ExecContext(
		ctx,
		query,
		credentials.Name,
		credentials.Email,
		credentials.Password,
	)
	if err != nil {
		return -1, derr.NewInternalError("failed to execute query")
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return -1, derr.NewInternalError("failed to get last inserted ID")
	}

	return userID, nil
}
