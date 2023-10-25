package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"html"
	"strings"
	"time"
)

type User struct {
	gorm.Model
	Email     string `gorm:"size:255;not null;unique" json:"email"`
	Password  string `gorm:"size:255;not null;" json:"password"`
	Companies []CompanyUsers
}

func GetUserByEmail(email string) (User, error) {
	var user User

	err := DB.Where("email = ?", email).Preload("Companies.Company").First(&user).Error
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

func (u *User) SaveUser(form string) (*User, error) {
	var err error

	err = u.beforeUserSaved()
	if err != nil {
		return nil, err
	}

	err = DB.Create(&u).Error
	fmt.Println("launched from: \t", form)
	if err != nil {
		return &User{}, err
	}

	return u, nil
}

func (u *User) IsUserDeleted() bool {
	return u.DeletedAt != nil
}

func (u *User) HasBeenDataModified(date time.Time) bool {
	return u.UpdatedAt.Unix() > date.Unix()
}

func (u *User) beforeUserSaved() error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	fmt.Println("before save - before assigning \t ", u.Password)
	u.Password = string(hashPassword)

	fmt.Println("before save - after assigning \t", u.Password)
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))

	return nil
}

func VerifyPassword(password, hashedPassword string) error {
	hpass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	fmt.Println(password, hashedPassword, string(hpass))
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
