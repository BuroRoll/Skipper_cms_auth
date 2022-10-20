package repository

import (
	"Skipper_cms_auth/pkg/models"
	"Skipper_cms_auth/pkg/models/forms/inputForms"
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

func (r *AuthPostgres) CreateUser(userRegister inputForms.SignUpUserForm) (uint, error) {
	var user models.User
	user = models.User{
		Phone:      userRegister.Phone,
		FirstName:  userRegister.FirstName,
		SecondName: userRegister.SecondName,
		Password:   userRegister.Password,
	}
	result := r.db.Create(&user)
	if result.Error != nil {
		return user.ID, result.Error
	}
	return user.ID, nil
}
