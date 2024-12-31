package db

import (
	"fmt"
	"log"

	"github.com/alexshelto/tigres-tracker/db/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(databaseFile string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(databaseFile), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	err = DB.AutoMigrate(&models.User{}, &models.Song{})
	if err != nil {
		log.Fatal("Error during auto-migration: ", err)
	}

	fmt.Println("Database connected and tables migrated.")
}

func GetDB() *gorm.DB {
	return DB
}

func AddSongAndIncrementUser(db *gorm.DB, songName string, guildId uint, requestedByID uint) error {
	song := models.Song{
		SongName:    songName,
		RequestedBy: requestedByID,
		GuildId:     guildId,
	}

	// Create song in DB
	if err := song.CreateSong(db); err != nil {
		return fmt.Errorf("error creating song: %v", err)
	}

	user, err := models.GetOrCreateUser(db, requestedByID)
	if err != nil {
		return fmt.Errorf("error getting or creating user: %v", err)
	}

	// Increment the user's song count
	if err := user.IncrementSongCount(db); err != nil {
		return fmt.Errorf("error incrementing song count: %v", err)
	}

	log.Printf("Song '%s' added by user %d\n", songName, requestedByID)

	return nil
}
