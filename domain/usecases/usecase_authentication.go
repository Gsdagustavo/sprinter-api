package usecases

import (
	"context"
	"errors"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Gsdagustavo/sprinter-api/domain"
	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
	"github.com/Gsdagustavo/sprinter-api/domain/rules"
	"github.com/Gsdagustavo/sprinter-api/domain/util"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/datastore"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/filestore"
)

func NewAuthenticationUseCase(
	repository datastore.AuthRepository,
	securityKey string,
	storage filestore.FileStorage,
) domain.AuthenticationUseCase {
	return authenticationUseCase{
		repository:  repository,
		securityKey: securityKey,
		storage:     storage,
	}
}

type authenticationUseCase struct {
	repository  datastore.AuthRepository
	securityKey string
	storage     filestore.FileStorage
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

	valid, err := a.repository.CheckUserCredentials(ctx, credentials)
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

	hashedPassword, err := util.Hash(credentials.Password)
	if err != nil {
		return nil, derr.JoinInternalError(err, "failed to hash password")
	}

	credentials.Password = hashedPassword

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

func (a authenticationUseCase) GetUserByToken(
	ctx context.Context,
	token string,
) (*entities.User, error) {
	id, expired, err := util.GetUserIDFromToken(token, a.securityKey)
	if err != nil {
		return nil, derr.JoinInternalError(err, "failed to get user ID from token")
	}

	if expired {
		return nil, derr.NewUnauthorizedError("unauthorized")
	}

	return a.repository.GetUserByID(ctx, int64(id))
}

func (a authenticationUseCase) AttemptCompleteRegistration(
	ctx context.Context,
	information *entities.AccountInformation,
) (int64, error) {
	information.Username = strings.TrimSpace(information.Username)
	information.Biography = strings.TrimSpace(information.Biography)

	if !rules.ValidateName(information.Username) {
		return 0, derr.InvalidUsername
	}

	if !rules.ValidateBiography(information.Biography) {
		return 0, derr.InvalidBiography
	}

	return a.repository.AttemptCompleteRegistration(ctx, *information)
}

func (a authenticationUseCase) RegisterProfileImage(
	ctx context.Context,
	userID int64,
	image []byte,
) error {
	_, err := a.repository.GetUserByID(ctx, userID)
	if err != nil {
		return derr.JoinInternalError(err, "failed to get user by id")
	}

	err = a.storage.CreateAll("/user/profile")
	if err != nil {
		return derr.JoinInternalError(err, "failed to create user profile folder")
	}

	path := filepath.Join("/user/profile", strconv.FormatInt(userID, 10))
	err = a.storage.UploadFile(path, image)
	if err != nil {
		return derr.JoinInternalError(err, "failed to upload image")
	}

	return nil
}
