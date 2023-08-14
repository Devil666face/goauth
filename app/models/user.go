package models

import (
	"auth/app/database"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null" form:"username"`
	Password string `gorm:"not null" form:"password"`
	Admin    bool   `gorm:"default:false" form:"admin"`
}

type UserForm struct {
	Username        string `form:"username"`
	Password        string `form:"password"`
	PasswordConfirm string `form:"password_confirm"`
	Admin           string `form:"admin"`
}

func CreateUser(user *User) *gorm.DB {
	return database.DB.Create(user)
}

func GetUser(dest *User, user string) *gorm.DB {
	return database.DB.First(dest, user)
}

func GetUserByUsername(dest *User, username string) *gorm.DB {
	return database.DB.Where("username= ?", username).Take(&dest)
}
