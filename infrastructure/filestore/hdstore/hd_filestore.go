package hdstore

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/VitorFranciscoDev/sprinter-api/domain/entities"
	"github.com/VitorFranciscoDev/sprinter-api/infrastructure/filestore"
)

type hdFileStorage struct {
	storageFolder string
}

func NewHDFileStorage(
	config entities.Config,
) filestore.FileStorage {
	fileStorage := hdFileStorage{
		storageFolder: config.FileStorage.StorageFolder,
	}

	// Attempt to set up storage
	err := fileStorage.Setup()
	if err != nil {
		panic(fmt.Errorf("failed to setup file storage: %v", err))
	}

	return fileStorage
}

func (h hdFileStorage) Setup() error {
	// Create the storage folder if not found
	err := os.MkdirAll(h.storageFolder, os.ModePerm)
	if err != nil {
		return errors.Join(errors.New("failed to create folder"))
	}

	return nil
}

func (h hdFileStorage) Exists(path string) (bool, error) {
	// Add a leading slash if needed
	path = formatLeadingSlash(path)

	fullPath := fmt.Sprintf("%s%s", h.storageFolder, path)
	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, errors.Join(errors.New("failed to check if file exists"))
	}
	return true, nil
}

func (h hdFileStorage) CreateAll(path string) error {
	// Add a leading slash if needed
	path = formatLeadingSlash(path)

	fullPath := fmt.Sprintf("%s%s", h.storageFolder, path)
	return os.MkdirAll(fullPath, os.ModePerm)
}

func (h hdFileStorage) ServeFile(path string) (*os.File, error) {
	// Add a leading slash if not there yet
	path = formatLeadingSlash(path)

	fullPath := fmt.Sprintf("%s%s", h.storageFolder, path)

	_, err := os.Stat(fullPath)
	if err != nil {
		return nil, os.ErrNotExist
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, errors.Join(errors.New("failed to open file"))
	}

	return file, nil
}

func (h hdFileStorage) UploadFile(
	path string,
	bytes []byte,
) error {
	// Add a leading slash if not there yet
	path = formatLeadingSlash(path)

	// Create file
	fullPath := fmt.Sprintf("%s%s", h.storageFolder, path)
	file, err := os.Create(fullPath)
	if err != nil {
		return errors.Join(errors.New("failed to create file"))
	}
	defer file.Close()

	// Write bytes to the file
	_, err = file.Write(bytes)
	if err != nil {
		return errors.Join(errors.New("failed to write to file"))
	}

	return nil
}

func (h hdFileStorage) DeleteFile(path string) error {
	// Add a leading slash if not there yet
	path = formatLeadingSlash(path)

	// Create the file path
	fullPath := fmt.Sprintf("%s%s", h.storageFolder, path)

	// Check if the file exists
	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return errors.Join(errors.New("failed to stat file"))
	}

	// Delete the file
	err = os.Remove(fullPath)
	if err != nil {
		return errors.Join(errors.New("failed to remove file"))
	}

	return nil
}

func formatLeadingSlash(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}
