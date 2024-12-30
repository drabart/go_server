package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Nickname     string `gorm:"uniqueIndex"`
	PasswordHash string
	IsAdmin      bool
}
