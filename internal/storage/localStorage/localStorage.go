package localstorage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

type Storage struct{
	StoragePath string
	StorageType string
}

func New(storagePath string) (*Storage, error) {
	err := os.MkdirAll(storagePath, os.ModePerm)

	if err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &Storage{
		StoragePath: storagePath,
		StorageType: "local",
	}, nil
}

func (s *Storage) GetStoragePath() string {
	return s.StoragePath
}

func (s *Storage) GetStorageType() string {
	return s.StorageType
}

func (s *Storage) createFilePath(name string) string {
	return fmt.Sprintf("%s/%s", s.StoragePath, name)
}

func (s *Storage) SaveFile(file multipart.File, name string) error {
	dst, err := os.Create(s.createFilePath(name))

	if err != nil {
		return err
	}

	if _, err := io.Copy(dst, file); err != nil {
		return err
	}

	defer dst.Close()

	return nil
}

func (s *Storage) GetFile(name string) ([]byte, error) {
	buf, err := os.ReadFile(s.createFilePath(name))

	if err != nil {
		return nil, err
	}

	return buf, nil
}
