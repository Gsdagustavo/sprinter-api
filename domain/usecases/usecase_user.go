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

func (u userUseCase) EditUserProfile(ctx context.Context, editedUser entities.EditUserProfileDTO) (string, error) {
	if err := rules.ValidateName(editedUser.Username); err != nil {
		return "", err
	}
	if err := rules.ValidateBiography(editedUser.Biography); err != nil {
		return "", err
	}
	err := u.userRepo.EditUserProfile(ctx, editedUser)
	if err != nil {
		return "", err
	}
	return "Edited user with Succes!", nil
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
