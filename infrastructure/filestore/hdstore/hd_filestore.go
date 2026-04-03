package hdstore

import (
	"fmt"
	"os"
	"strings"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/filestore"
)

type hdFileStorage struct {
	storageFolder string
}

func NewHDFileStorage(
	config entities.Settings,
) filestore.FileStorage {
	fileStorage := hdFileStorage{
		storageFolder: config.StorageSettings.StorageFolder,
	}

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
		return derr.JoinError("failed to create folder", err)
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

		return false, derr.JoinError("failed to check if file exists", err)
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
		return nil, derr.JoinError("failed to open file", err)
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
		return derr.JoinError("failed to create file", err)
	}
	defer file.Close()

	// Write bytes to the file
	_, err = file.Write(bytes)
	if err != nil {
		return derr.JoinError("failed to write to file", err)
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

		return derr.JoinError("failed to stat file", err)
	}

	// Delete the file
	err = os.Remove(fullPath)
	if err != nil {
		return derr.JoinError("failed to remove file", err)
	}

	return nil
}

func formatLeadingSlash(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}
