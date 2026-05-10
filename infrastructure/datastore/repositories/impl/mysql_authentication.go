package repositories

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
	"github.com/Gsdagustavo/sprinter-api/domain/util"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/datastore/repositories"
)

func NewAuthenticationRepository(settings repositories.SettingsRepository) repositories.AuthRepository {
	return &authenticationRepository{
		conn:     settings.Connection(),
		settings: settings,
	}
}

type authenticationRepository struct {
	conn     *sql.DB
	settings repositories.SettingsRepository
}

//go:embed _query/auth/get_user_by_email.sql
var getUserByEmail string

//go:embed _query/auth/attempt_register.sql
var attemptRegister string

//go:embed _query/auth/check_user_credential.sql
var checkUserCredential string

//go:embed _query/auth/get_user_by_id.sql
var getUserById string

//go:embed _query/auth/attempt_complete_register.sql
var attemptCompleteRegister string

func (r authenticationRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*entities.User, error) {
	var user entities.User
	row := r.conn.QueryRowContext(ctx, getUserByEmail, email)
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

		return nil, derr.JoinError("failed to scan", err)
	}

	return &user, nil
}

func (r authenticationRepository) AttemptRegister(
	ctx context.Context,
	credentials entities.UserCredentials,
) (int64, error) {
	result, err := r.conn.ExecContext(
		ctx,
		attemptRegister,
		credentials.Name,
		credentials.Email,
		credentials.Password,
	)
	if err != nil {
		return -1, derr.JoinError("failed to execute query", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return -1, derr.JoinError("failed to get last inserted ID", err)
	}

	return userID, nil
}

func (r authenticationRepository) CheckUserCredentials(
	ctx context.Context,
	credentials entities.UserCredentials,
) (bool, error) {
	var password string
	err := r.conn.QueryRowContext(ctx, checkUserCredential, credentials.Email).Scan(&password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, derr.NotFoundError
		}

		return false, derr.JoinError("failed to query or scan", err)
	}

	valid := util.CheckValidPassword(credentials.Password, password)
	return valid, nil
}

func (r authenticationRepository) GetUserByID(
	ctx context.Context,
	userID int64,
) (*entities.User, error) {
	var user entities.User
	row := r.conn.QueryRowContext(ctx, getUserById, userID)
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

		return nil, derr.JoinError("failed to scan", err)
	}

	return &user, nil
}

func (r authenticationRepository) AttemptCompleteRegistration(
	ctx context.Context,
	information entities.UserInformation,
) error {
	_, err := r.conn.ExecContext(
		ctx,
		attemptCompleteRegister,
		information.Username,
		information.Biography,
	)
	if err != nil {
		return derr.JoinError("failed to execute query", err)
	}

	return nil
}
