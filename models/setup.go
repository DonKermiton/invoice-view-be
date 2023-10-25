package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"os"
)

var DB *gorm.DB

func ConnectDatabase() {

	err := godotenv.Load(".env")

	if err != nil {
		fmt.Errorf("error loading .env file")
	}

	Dbdriver := os.Getenv("DB_DRIVER")
	DbHost := os.Getenv("DB_HOST")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbPort := os.Getenv("DB_PORT")

	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
	DB, err = gorm.Open(Dbdriver, DBURL)

	if err != nil {
		fmt.Printf("Cannot connect to database %s; \n %s", err, DBURL)
	} else {
		fmt.Printf("Connected")
	}

	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&CompanyUsers{})
	DB.AutoMigrate(&Company{})
}
