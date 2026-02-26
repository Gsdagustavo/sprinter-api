package usecases

import (
	"context"
	"errors"

	"github.com/VitorFranciscoDev/sprinter-api/domain"
	"github.com/VitorFranciscoDev/sprinter-api/domain/entities"
	"github.com/VitorFranciscoDev/sprinter-api/domain/entities/derr"
	"github.com/VitorFranciscoDev/sprinter-api/domain/rules"
	"github.com/VitorFranciscoDev/sprinter-api/domain/util"
	"github.com/VitorFranciscoDev/sprinter-api/infrastructure/datastore"
)

func NewAuthenticationUseCase(
	r datastore.AuthRepository,
	securityKey string,
) domain.AuthenticationUseCase {
	return authenticationUseCase{
		repository:  r,
		securityKey: securityKey,
	}
}

type authenticationUseCase struct {
	repository  datastore.AuthRepository
	securityKey string
}

func (a authenticationUseCase) AttemptLogin(
	ctx context.Context,
	credentials entities.UserCredentials,
) (*entities.AuthenticationResponse, error) {
	if !rules.ValidateCredentials(credentials) {
		return nil, derr.InvalidCredentials
	}

	userByEmail, err := a.repository.GetUserByEmail(ctx, credentials.Email)
	if err != nil && !errors.Is(err, derr.NotFoundError) {
		return nil, derr.JoinInternalError(err, "failed to get user by email")
	}

	if userByEmail == nil {
		return nil, derr.InvalidCredentials
	}

	valid, err := a.repository.CheckValidPassword(ctx, userByEmail.ID, credentials.Password)
	if err != nil {
		return nil, derr.JoinInternalError(err, "login attempt failed")
	}

	if !valid {
		return nil, derr.InvalidCredentials
	}

	token, err := util.GetNewAuthToken(userByEmail.ID, a.securityKey)
	if err != nil {
		return nil, derr.JoinInternalError(err, "failed to generate token")
	}

	return &entities.AuthenticationResponse{Token: token}, nil
}

func (a authenticationUseCase) AttemptRegister(
	ctx context.Context,
	credentials entities.UserCredentials,
) (*entities.AuthenticationResponse, error) {
	if rules.ValidateName(credentials.Name) == false {
		return nil, derr.InvalidCredentials
	}

	if !rules.ValidateCredentials(credentials) {
		return nil, derr.InvalidCredentials
	}

	// Check if a user with the email already exists in the database
	userByEmail, err := a.repository.GetUserByEmail(ctx, credentials.Email)
	if err != nil && !errors.Is(err, derr.NotFoundError) {
		return nil, derr.JoinInternalError(err, "failed to get user by email")
	}

	if userByEmail != nil {
		return nil, derr.UserAlreadyExists
	}

	userID, err := a.repository.AttemptRegister(ctx, credentials)
	if err != nil {
		return nil, derr.JoinInternalError(err, "register attempt failed")
	}

	token, err := util.GetNewAuthToken(userID, a.securityKey)
	if err != nil {
		return nil, derr.JoinInternalError(err, "failed to generate token")
	}

	return &entities.AuthenticationResponse{Token: token}, nil
}

func (a authenticationUseCase) GetUserByEmail(
	ctx context.Context,
	email string,
) (*entities.User, error) {
	return a.repository.GetUserByEmail(ctx, email)
}

func (a authenticationUseCase) CheckCredentials(
	ctx context.Context,
	credentials entities.UserCredentials,
) (bool, error) {
	return a.repository.CheckUserCredentials(ctx, credentials)
}
