package localstorage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

const fileStoragePath = "./storage/files"

type Storage struct{}

func New() (*Storage, error) {
	err := os.MkdirAll(fileStoragePath, os.ModePerm)

	if err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &Storage{}, nil
}

func (s *Storage) SaveFile(file multipart.File, name string) error {
	dst, err := os.Create(fmt.Sprintf("%s/%s", fileStoragePath, name))

	if err != nil {
		return err
	}

	if _, err := io.Copy(dst, file); err != nil {
		return err
	}

	defer dst.Close()

	return nil
}
