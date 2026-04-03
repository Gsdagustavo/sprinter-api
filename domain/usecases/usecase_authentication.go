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
) (string, error) {
	if err := rules.ValidateCredentials(credentials); err != nil {
		return "", derr.InvalidCredentials
	}

	userByEmail, err := a.repository.GetUserByEmail(ctx, credentials.Email)
	if err != nil && !errors.Is(err, derr.NotFoundError) {
		return "", derr.JoinError("failed to get user by email", err)
	}

	if userByEmail == nil {
		return "", derr.InvalidCredentials
	}

	valid, err := a.repository.CheckUserCredentials(ctx, credentials)
	if err != nil {
		return "", derr.JoinError("login attempt failed", err)
	}

	if !valid {
		return "", derr.InvalidCredentials
	}

	token, err := util.GetNewAuthToken(userByEmail.ID, a.securityKey)
	if err != nil {
		return "", derr.JoinError("failed to generate token", err)
	}

	return token, nil
}

func (a authenticationUseCase) AttemptRegister(
	ctx context.Context,
	credentials entities.UserCredentials,
) (string, error) {
	if err := rules.ValidateCredentials(credentials); err != nil {
		return "", err
	}

	// Check if a user with the email already exists in the database
	userByEmail, err := a.repository.GetUserByEmail(ctx, credentials.Email)
	if err != nil && !errors.Is(err, derr.NotFoundError) {
		return "", derr.JoinError("failed to get user by email", err)
	}

	if userByEmail != nil {
		return "", derr.UserAlreadyExists
	}

	hashedPassword, err := util.Hash(credentials.Password)
	if err != nil {
		return "", derr.JoinError("failed to hash password", err)
	}

	credentials.Password = hashedPassword

	userID, err := a.repository.AttemptRegister(ctx, credentials)
	if err != nil {
		return "", derr.JoinError("register attempt failed", err)
	}

	token, err := util.GetNewAuthToken(userID, a.securityKey)
	if err != nil {
		return "", derr.JoinError("failed to generate token", err)
	}

	return token, nil
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
		return nil, derr.JoinError("failed to get user ID from token", err)
	}

	if expired {
		return nil, derr.NewUnauthorizedError("unauthorized")
	}

	return a.repository.GetUserByID(ctx, int64(id))
}

func (a authenticationUseCase) AttemptCompleteRegistration(
	ctx context.Context,
	information entities.AccountInformation,
) error {
	information.Username = strings.TrimSpace(information.Username)
	information.Biography = strings.TrimSpace(information.Biography)

	var err error
	if err = rules.ValidateName(information.Username); err != nil {
		return err
	}

	if err = rules.ValidateBiography(information.Biography); err != nil {
		return err
	}

	err = a.repository.AttemptCompleteRegistration(ctx, information)
	if err != nil {
		return derr.JoinError("failed to complete registration", err)
	}

	return nil
}

func (a authenticationUseCase) UploadProfileImage(
	_ context.Context,
	userID int64,
	image []byte,
) error {
	const profileFolderPath = "/user/profile"

	err := a.storage.CreateAll(profileFolderPath)
	if err != nil {
		return derr.JoinError("failed to create user profile folder", err)
	}

	path := filepath.Join(profileFolderPath, strconv.FormatInt(userID, 10))
	err = a.storage.UploadFile(path, image)
	if err != nil {
		return derr.JoinError("failed to upload image", err)
	}

	return nil
}
