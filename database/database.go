package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// Menggunakan path project yang benar
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Auto Migrate the database
	// err = DB.AutoMigrate(
	// 	&models.User{},
	// 	&models.Car{},
	// 	&models.Driver{},
	// 	&models.EventPackage{},
	// 	&models.RentalHistory{},
	// 	&models.CallAssistance{},
	// 	&models.Membership{},
	// )

	// if err != nil {
	// 	log.Fatal("Failed to auto migrate: ", err)
	// }
}
