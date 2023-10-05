package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"html"
	"strings"
	"time"
)

type User struct {
	gorm.Model
	Email    string `gorm:"size:255;not null;unique" json:"email"`
	Password string `gorm:"size:255;not null;" json:"password"`
}

func GetUserByEmail(email string) (User, error) {
	var user User

	err := DB.Where("email = ?", email).First(&user).Error

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func GetUserById(id uint) (User, error) {
	var user User

	if err := DB.First(&user, id).Error; err != nil {
		return user, errors.New("user not found")
	}

	user.PrepareToSend()

	return user, nil
}

func (u *User) PrepareToSend() {
	u.Password = ""
}

func (u *User) SaveUser() (*User, error) {
	var err error
	err = DB.Create(&u).Error

	if err != nil {
		return &User{}, err
	}

	return u, nil
}

func (u *User) IsUserDeleted() bool {
	return u.DeletedAt != nil
}

func (u User) HasBeenDataModified(date time.Time) bool {
	return u.UpdatedAt.Unix() > date.Unix()
}

func (u *User) BeforeSave() error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.Password = string(hashPassword)
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))

	return nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
