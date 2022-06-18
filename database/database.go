package database

import (
	"backend/database/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var database *gorm.DB = nil

func Connect() {
	db, err := gorm.Open(sqlite.Open("database.db"))
	if err != nil {
		panic("Database not working")
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.Order{})
	db.AutoMigrate(&models.ProductOrder{}, &models.Order{}, &models.Product{})

	database = db
}

func GetDatabase() *gorm.DB {
	return database
}
