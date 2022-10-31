package repository

import (
	"Skipper_cms_auth/pkg/models"
	"errors"
	"gorm.io/gorm"
)

type AuthPostgres struct {
	db *gorm.DB
}

func NewAuthPostgres(db *gorm.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) GetUser(login, password string) (uint, []models.Role, error) {
	var user models.User
	result := r.db.Preload("Role").Where("email=? AND password=? OR phone=? AND password=?", login, password, login, password).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user.ID, nil, gorm.ErrRecordNotFound
	}
	return user.ID, user.Role, nil
}

func (r *AuthPostgres) GetUserById(userId uint) (models.User, error) {
	var user models.User
	result := r.db.Preload("Role").First(&user, userId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, gorm.ErrRecordNotFound
	}
	return user, nil
}
