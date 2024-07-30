package uuidgenerator

import "github.com/google/uuid"

type Uuid struct {}

func New() *Uuid {
	return &Uuid{}
}

func (u *Uuid) GenerateUUID() string {
	return uuid.New().String()
}