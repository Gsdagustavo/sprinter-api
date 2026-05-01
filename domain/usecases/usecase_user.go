package usecases

import (
	"context"
	"fmt"
	"path"

	"github.com/Gsdagustavo/sprinter-api/domain"
	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
	"github.com/Gsdagustavo/sprinter-api/domain/rules"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/datastore"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/filestore"
)

func NewUserUseCases(
	storage filestore.FileStorage,
	storageConfig entities.FileStorageSettings,
	userRepo datastore.UserRepository,
) domain.UserUseCase {
	return userUseCase{
		storage:       storage,
		storageConfig: storageConfig,
		userRepo:      userRepo,
	}
}

type userUseCase struct {
	storage       filestore.FileStorage
	storageConfig entities.FileStorageSettings
	userRepo      datastore.UserRepository
}

func (u userUseCase) UpdateUserInformation(
	ctx context.Context,
	userInformation entities.UserInformation,
) (*entities.User, error) {
	err := rules.ValidateUserInformation(userInformation)
	if err != nil {
		return nil, err
	}

	err = u.userRepo.UpdateUserInformation(ctx, userInformation)
	if err != nil {
		return nil, derr.JoinError("failed to update the user information", err)
	}

	user, err := u.userRepo.GetUserById(ctx, userInformation.ID)
	if err != nil {
		return nil, derr.JoinError("failed to get the updated user ", err)
	}

	return user, nil
}

func (u userUseCase) SaveUserProfilePicture(
	userID int64,
	image []byte,
) (string, error) {
	userFolder := u.getUserFolder(userID)
	imagePath := path.Join(userFolder, "profile.jpg")

	err := u.storage.DeleteFile(imagePath)
	if err != nil {
		return "", derr.JoinError("failed to delete file", err)
	}

	err = u.storage.UploadFile(imagePath, image)
	if err != nil {
		return "", derr.JoinError("failed to upload file", err)
	}

	return imagePath, nil
}

func (u userUseCase) getUserFolder(
	userID int64,
) string {
	storageFolder := u.storageConfig.StorageFolder
	userFolder := path.Join(storageFolder, fmt.Sprintf("%d", userID))
	return userFolder
}
