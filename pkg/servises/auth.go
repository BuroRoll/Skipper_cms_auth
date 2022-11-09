package service

import (
	"Skipper_cms_auth/pkg/models"
	"Skipper_cms_auth/pkg/repository"
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"gopkg.in/gomail.v2"
	"html/template"
	"log"
	"path/filepath"
	"runtime"
	"time"
)

const (
	salt                  = "14hjqrhj1231qw124617ajfha1123ssfqa3ssjs190"
	signingKey            = "qrkjk#4#%35FSFJlja#4353KSFjH"
	signingRefreshKey     = "qrkjk#sdfioh12bkj@nkk3k1axv["
	tokenTTL              = time.Second * 10
	refreshTokenTTL       = time.Hour * 12 * 365
	resetPasswordTokenTTL = time.Hour
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

type resetPasswordTokenClaims struct {
	jwt.StandardClaims
	UserId uint `json:"user_id"`
}

type refreshTokenClaims struct {
	jwt.StandardClaims
	UserId uint `json:"user_id"`
}

func (s *AuthService) GetUser(login, password string) (uint, []models.Role, error) {
	return s.repo.GetUser(login, generatePasswordHash(password))
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

func (s *AuthService) SendEmailToResetPassword(email string) error {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return errors.New("Пользователь не найден ")
	}
	err = s.SendVerifyEmail(user.ID, user.Email)
	return err
}

func (s *AuthService) SendVerifyEmail(userId uint, email string) error {
	token, err := GenerateTokenForResetPassword(userId)
	err = SendEmailToVerify(email, token)
	if err != nil {
		return errors.New("Не удалось отправить сообщение ")
	}
	return nil
}

func GenerateTokenForResetPassword(userId uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &resetPasswordTokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(resetPasswordTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userId,
	})
	resetPasswordToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", errors.New("Ошибка генерации токена сброса пароля ")
	}
	return resetPasswordToken, nil
}

func SendEmailToVerify(email string, token string) error {
	type data struct {
		Token string
		Link  string
	}
	userData := data{
		Token: token,
		Link:  "https://skipper.gq/reset-password?",
	}
	_, b, _, _ := runtime.Caller(0)
	Root := filepath.Join(filepath.Dir(b), "../..")
	t := template.New("resetPassword.html")
	var err error
	t, err = t.ParseFiles(Root + "/resetPassword.html")
	if err != nil {
		log.Println(err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, userData); err != nil {
		log.Println(err)
	}
	result := tpl.String()
	m := gomail.NewMessage()
	m.SetHeader("From", "buroroll@ya.ru")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Сброс пароля")
	m.SetBody("text/html", result)
	d := gomail.NewDialer("smtp.yandex.ru", 465, "buroroll@ya.ru", "zubsglsuzxxmqyht")
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
