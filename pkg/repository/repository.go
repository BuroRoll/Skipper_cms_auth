package repository

import (
	"Skipper_cms_auth/pkg/models"
	"gorm.io/gorm"
)

type Authorization interface {
	GetUser(email, password string) (uint, []models.Role, error)
	GetUserById(userId uint) (models.User, error)
}

type UserData interface {
	GetUserById(userId uint) (models.User, error)
}

type Repository struct {
	Authorization
	UserData
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
