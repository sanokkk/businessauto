package models

import "github.com/google/uuid"

type User struct {
	Id           uuid.UUID
	Email        string
	FullName     string `gorm:"column:fullname"`
	PasswordHash string `gorm:"column:passwordhash"`
	Role         string
}
