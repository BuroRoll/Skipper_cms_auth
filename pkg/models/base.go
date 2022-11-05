package models

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var cmsDb *gorm.DB

func init() {
	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable port=%s password=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_PASSWORD"))
	conn, err := gorm.Open(postgres.Open(dbUri), &gorm.Config{})
	if err != nil {
		log.Fatalf("error %s", err)
	}
	cmsDb = conn

	if err = cmsDb.AutoMigrate(
		&User{},
		&Role{},
		&UserRoles{},
	); err == nil && cmsDb.Migrator().HasTable(&Role{}) {
		if err := cmsDb.First(&Role{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			var roles = []Role{{Name: "super_admin"}, {Name: "admin"}, {Name: "support"}, {Name: "editor"}}
			cmsDb.Create(&roles)
		}
		if err := cmsDb.First(&User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			var firstUser = User{
				FirstName:  "admin",
				SecondName: "admin",
				Email:      "admin@admin.ru",
				Password:   GeneratePasswordHash("admin"),
			}
			cmsDb.Create(&firstUser)
			cmsDb.Model(&firstUser).Association("Role").Append(&Role{Name: "super_admin"})
		}
	}
	err = cmsDb.SetupJoinTable(&User{}, "Role", &UserRoles{})
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func GetDB() *gorm.DB {
	return cmsDb
}

func GeneratePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte("14hjqrhj1231qw124617ajfha1123ssfqa3ssjs190")))
}
