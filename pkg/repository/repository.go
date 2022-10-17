package repository

import (
	"Skipper_cms_auth/pkg/models"
	"Skipper_cms_auth/pkg/models/forms"
	"gorm.io/gorm"
)

type Authorization interface {
	CreateUser(user forms.SignUpUserForm) (uint, error)
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
