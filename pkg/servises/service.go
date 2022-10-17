package service

import (
	"Skipper_cms_auth/pkg/models"
	"Skipper_cms_auth/pkg/models/forms"
	"Skipper_cms_auth/pkg/repository"
)

type Authorization interface {
	CreateUser(user forms.SignUpUserForm) (uint, error)
	GetUser(login, password string) (uint, []models.Role, error)
	GenerateToken(login, password string) (string, string, error)
	GenerateTokenByID(userId uint, roles []models.Role) (string, string, error)
	ParseToken(token string) (uint, []models.Role, error)
	ParseRefreshToken(token string) (uint, error)
	GetUserData(userId uint) (models.User, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
