package usecases

import (
	"fmt"
	"path"

	"github.com/Gsdagustavo/sprinter-api/domain"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/filestore"
)

func NewUserUseCases(
	storage filestore.FileStorage,
	storageFolder string,
) domain.UserUseCase {
	return userUseCase{
		storage:       storage,
		storageFolder: storageFolder,
	}
}

type userUseCase struct {
	storage       filestore.FileStorage
	storageFolder string
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
