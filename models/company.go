package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Company struct {
	gorm.Model
	Name      string `gorm:"size:255,not null" json:"companyName"`
	Companies []CompanyUsers
}

func (c *Company) CreateCompany() (*Company, error) {
	var err error
	err = DB.Create(&c).Error

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Company) CreateCompanyWithUser(u User) (*Company, User, error) {
	err := DB.Transaction(func(tx *gorm.DB) error {
		_, err := c.CreateCompany()
		if err != nil {
			return err
		}

		fmt.Println("Create company with user", u.Password)
		_, err = u.SaveUser("Create Company With User")

		if err != nil {
			return err
		}

		companyUser := CompanyUsers{
			UserID:    u.ID,
			CompanyID: c.ID,
			User:      u,
			Company:   *c,
		}

		_, err = companyUser.SaveUserToCompany()

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, User{}, err
	}

	return c, u, nil
}
