package usecases

import (
	"context"

	"github.com/VitorFranciscoDev/sprinter-api/domain/entities"
	"github.com/VitorFranciscoDev/sprinter-api/domain/entities/derr"
	"github.com/VitorFranciscoDev/sprinter-api/domain/rules"
	"github.com/VitorFranciscoDev/sprinter-api/domain/util"
	"github.com/VitorFranciscoDev/sprinter-api/infrastructure/datastore"
)

type AuthenticationUseCase struct {
	repository  datastore.AuthRepository
	securityKey string
}

// NewAuthenticationUseCase instantiates a new Auth Use Case
func NewAuthenticationUseCase(
	r datastore.AuthRepository,
	securityKey string,
) *AuthenticationUseCase {
	return &AuthenticationUseCase{
		repository:  r,
		securityKey: securityKey,
	}
}

// AttemptLogin attempt to log in a user with the given credentials
func (a *AuthenticationUseCase) AttemptLogin(
	ctx context.Context,
	credentials entities.UserCredentials,
) (*entities.AuthenticationResponse, error) {
	if !rules.ValidateCredentials(credentials) {
		return nil, derr.InvalidCredentials
	}

	userID, didLogin, err := a.repository.AttemptLogin(ctx, credentials)
	if err != nil {
		return nil, derr.JoinInternalError(err, "login attempt failed")
	}

	if !didLogin {
		return nil, derr.InvalidCredentials
	}

	token, err := util.GetNewAuthToken(userID, a.securityKey)
	if err != nil {
		return nil, derr.JoinInternalError(err, "failed to generate token")
	}

	return &entities.AuthenticationResponse{Token: token}, nil
}

// AttemptRegister attempt to register a new user with the given credentials
func (a *AuthenticationUseCase) AttemptRegister(
	ctx context.Context,
	credentials entities.UserCredentials,
) (*entities.AuthenticationResponse, error) {
	if rules.ValidateName(credentials.Name) == false {
		return nil, derr.InvalidCredentials
	}

	if !rules.ValidateCredentials(credentials) {
		return nil, derr.InvalidCredentials
	}

	userID, didRegister, err := a.repository.AttemptLogin(ctx, credentials)
	if err != nil {
		return nil, derr.JoinInternalError(err, "register attempt failed")
	}

	if !didRegister {
		return nil, derr.InvalidCredentials
	}

	token, err := util.GetNewAuthToken(userID, a.securityKey)
	if err != nil {
		return nil, derr.JoinInternalError(err, "failed to generate token")
	}

	return &entities.AuthenticationResponse{Token: token}, nil
}
