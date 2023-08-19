package models

import (
	"fmt"

	"auth/app/database"
	"auth/app/utils"

	"gorm.io/gorm"
)

const (
	USER string = "User"
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

func (form *UserForm) IsAdmin() bool {
	if form.Admin != "" {
		return true
	}
	return false
}

func (form *UserForm) IsEmptyUsername() (string, bool) {
	if form.Username == "" {
		return "Username is required.", true
	}
	return "", false
}

func (form *UserForm) IsPasswordsMatch() (string, bool) {
	if form.Password != form.PasswordConfirm {
		return "The passwords don't match.", true
	}
	return "", false
}

func (form *UserForm) IsPasswordsEmpty() (string, bool) {
	if form.Password == "" || form.PasswordConfirm == "" {
		return "Password is required.", true
	}
	return "", false
}

func (form *UserForm) IsPasswordsShort() (string, bool) {
	if len([]rune(form.Password)) < utils.PASSWORDLEN {
		return fmt.Sprintf("The minimum len of password is %d", utils.PASSWORDLEN), true
	}
	return "", false
}

func (form *UserForm) CheckPasswordForCreate() (string, bool) {
	if message, ok := form.IsPasswordsMatch(); ok {
		return message, ok
	}
	if message, ok := form.IsPasswordsEmpty(); ok {
		return message, ok
	}
	if message, ok := form.IsPasswordsShort(); ok {
		return message, ok
	}
	return "", false
}

func (user *User) Set(username string, password string, admin bool) {
	user.Username = username
	user.Password = password
	user.Admin = admin
}

func CreateUser(user *User) *gorm.DB {
	return database.DB.Create(user)
}

func UpdateUser(user *User) *gorm.DB {
	return database.DB.Save(user)
}

func DeleteUser(user *User) *gorm.DB {
	return database.DB.Unscoped().Delete(user)
}

func GetUser(dest *User, id string) *gorm.DB {
	return database.DB.First(dest, id)
}

func GetUserByUsername(dest *User, username string) *gorm.DB {
	return database.DB.Where("username= ?", username).Take(&dest)
}

func GetAllUsers() []User {
	var users []User
	database.DB.Find(&users)
	return users
}
