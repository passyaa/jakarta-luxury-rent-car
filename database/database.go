package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Menggunakan path project yang benar

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))

	// Disable prepared statement caching
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Menonaktifkan penggunaan prepared statement
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Menambahkan logger untuk debugging
	})
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
