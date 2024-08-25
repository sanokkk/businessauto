package models

import "github.com/google/uuid"

type User struct {
	Id           uuid.UUID
	Email        string
	FullName     string
	PasswordHash string
	Role         string
}
