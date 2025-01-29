package db

import (
	"log"

	"github.com/12sub/websockets/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var data = "user:password@tcp(localhost:3306)/your_db?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	dsn := data
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	DB.AutoMigrate(&models.User{})
}