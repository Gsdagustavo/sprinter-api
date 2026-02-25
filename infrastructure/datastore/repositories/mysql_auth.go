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

func (r authenticationRepository) AttemptLogin(
	ctx context.Context,
	credentials entities.UserCredentials,
) (int64, bool, error) {
	const query = `
		SELECT id, password FROM users WHERE email = ?
	`

	var userID int64
	var passwordHash string
	err := r.conn.QueryRowContext(ctx, query, credentials.Email).Scan(&userID, &passwordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, false, derr.NewBadRequestError("invalid credentials")
		}

		return -1, false, derr.NewInternalError("failed to scan password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(credentials.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return -1, false, derr.NewBadRequestError("invalid credentials")
		}

		return -1, false, derr.NewInternalError("failed to compare password with hash")
	}

	return userID, true, nil
}

func (r authenticationRepository) AttemptRegister(
	ctx context.Context,
	credentials entities.UserCredentials,
) (int64, bool, error) {
	const query = `INSERT INTO users (name, email, password) VALUES (?, ?, ?, ?)`

	result, err := r.conn.ExecContext(
		ctx,
		query,
		credentials.Name,
		credentials.Email,
		credentials.Password,
	)
	if err != nil {
		return -1, false, derr.NewInternalError("failed to execute query")
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return -1, false, derr.NewInternalError("failed to get last inserted ID")
	}

	return userID, true, nil
}
