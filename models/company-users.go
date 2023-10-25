package models

import "github.com/jinzhu/gorm"

type CompanyUsers struct {
	gorm.Model
	CompanyID uint `gorm:"not null"`
	UserID    uint `gorm:"not null"`
	User      User
	Company   Company
}

//func (c *CompanyUsers) UserBelongsToCompany(company *Company) (bool, error) {
//
//
//}

func (c *CompanyUsers) SaveUserToCompany() (*CompanyUsers, error) {
	err := DB.Create(&c).Error

	if err != nil {
		return nil, err
	}

	return c, nil
}
