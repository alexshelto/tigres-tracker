package database

import (
	"fmt"
	"log"
	"os"

	"github.com/alexshelto/tigres-tracker/api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//var Db *gorm.DB

func ConnectDb() *gorm.DB {
	dsn := fmt.Sprintf("host=db user=%s password=%s dbname=%s port=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Error connecting to db: %v", err)
		os.Exit(2)
	}

	log.Println("Connected to db.")
	//db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("running migration")

	db.AutoMigrate(&models.Song{}, &models.User{}, &models.Play{})

	return db
}
