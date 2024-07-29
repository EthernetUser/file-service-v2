package database

import (
	"errors"
	"time"
)

var (
	ErrorNotFound = errors.New("not found")
	ErrorAlreadyExists = errors.New("already exists")
)

type FileToSave struct {
	OriginalName string
	Name         string
	Path         string
	Size         int64
}

type File struct {
	Id        int64
	OriginalName string
	Name         string
	Path         string
	Size         int
	Timestamp    time.Time
}
