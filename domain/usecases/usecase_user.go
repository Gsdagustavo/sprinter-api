package usecases

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/Gsdagustavo/sprinter-api/domain"
	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
	"github.com/Gsdagustavo/sprinter-api/domain/rules"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/datastore"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/filestore"
)

func NewUserUseCases(
	repository datastore.AuthRepository,
	storage filestore.FileStorage,
	storageFolder string,
) domain.UserUseCase {
	return userUseCase{
		repository:    repository,
		storage:       storage,
		storageFolder: storageFolder,
	}
}

type userUseCase struct {
	repository    datastore.AuthRepository
	storage       filestore.FileStorage
	storageFolder string
}

func (u userUseCase) EditUserProfile(
	ctx context.Context,
	user *entities.User,
	name string,
	description string,
) (*entities.User, error) {
	if user == nil {
		return nil, derr.InvalidParameterError
	}

	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)

	if name == "" && description == "" {
		return nil, derr.InvalidParameterError
	}

	updated := *user
	if name != "" {
		if !rules.ValidateName(name) {
			return nil, derr.InvalidParameterError
		}

		updated.Name = name
	}
	if description != "" && len(strings.Split(updated.Description, "")) <= 255 {
		updated.Description = description
	}

	err := u.repository.UpdateUserProfile(ctx, &updated)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}
func (u userUseCase) SaveUserProfilePicture(
	userID int64,
	image []byte,
) (string, error) {
	userFolder := u.getUserFolder(userID)
	imagePath := path.Join(userFolder, "profile.jpg")

	err := u.storage.DeleteFile(imagePath)
	if err != nil {
		return "", derr.JoinInternalError(err, "failed to delete file")
	}

	err = u.storage.UploadFile(imagePath, image)
	if err != nil {
		return "", derr.JoinInternalError(err, "failed to upload file")
	}

	return imagePath, nil
}

func (u userUseCase) getUserFolder(
	userID int64,
) string {
	storageFolder := u.storageFolder
	userFolder := path.Join(storageFolder, fmt.Sprintf("%d", userID))
	return userFolder
}
