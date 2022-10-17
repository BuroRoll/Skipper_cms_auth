package service

import (
	"Skipper_cms_auth/pkg/models"
	"Skipper_cms_auth/pkg/models/forms"
	"Skipper_cms_auth/pkg/repository"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

const (
	salt              = "14hjqrhj1231qw124617ajfha1123ssfqa3ssjs190"
	signingKey        = "qrkjk#4#%35FSFJlja#4353KSFjH"
	signingRefreshKey = "qrkjk#sdfioh12bkj@nkk3k1axv["
	tokenTTL          = time.Hour * 12
	refreshTokenTTL   = time.Hour * 12 * 365
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId uint          `json:"user_id"`
	Roles  []models.Role `json:"roles"`
}

type refreshTokenClaims struct {
	jwt.StandardClaims
	UserId uint `json:"user_id"`
}

func (s *AuthService) CreateUser(user forms.SignUpUserForm) (uint, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GetUser(login, password string) (uint, []models.Role, error) {
	return s.repo.GetUser(login, generatePasswordHash(password))
}

func generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) GenerateToken(login, password string) (string, string, error) {
	userId, roles, err := s.GetUser(login, password)
	if err != nil {
		return "", "", err
	}
	return s.GenerateTokenByID(userId, roles)
}

func (s *AuthService) GenerateTokenByID(userId uint, roles []models.Role) (string, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userId,
		roles,
	})
	t, err := token.SignedString([]byte(signingKey))
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &refreshTokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refreshTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userId,
	})
	rt, err := refreshToken.SignedString([]byte(signingRefreshKey))

	if err != nil {
		return "", "", err
	}

	return t, rt, err
}

func (s *AuthService) ParseRefreshToken(accessToken string) (uint, error) {
	token, err := jwt.ParseWithClaims(accessToken, &refreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingRefreshKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*refreshTokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func (s *AuthService) ParseToken(accessToken string) (uint, []models.Role, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, nil, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, nil, errors.New("token claims are not of type *tokenClaims")
	}
	return claims.UserId, claims.Roles, nil
}

func (s *AuthService) GetUserData(userId uint) (models.User, error) {
	return s.repo.GetUserById(userId)
}
