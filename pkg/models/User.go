package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	//Base information
	Email      string `gorm:"index:unique"`
	Phone      string `gorm:"index:unique"`
	Password   string `json:"-"`
	FirstName  string
	SecondName string
	Patronymic string
	Role       []Role `gorm:"many2many:user_roles;"`
}

type Role struct {
	Name string `gorm:"primaryKey"`
}

type UserRoles struct {
	UserID   uint   `gorm:"primaryKey"`
	RoleName string `gorm:"primaryKey"`
}
