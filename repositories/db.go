package repositories

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"go.mod/models"
	"log"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := "root:admin@tcp(127.0.0.1:3306)/go_db?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established.")

	autoMigrate()
}

func autoMigrate() {
	err := DB.AutoMigrate(&models.Hotel{}, &models.Room{}, &models.Guest{}, &models.Booking{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
	log.Println("Database migration completed.")
}
