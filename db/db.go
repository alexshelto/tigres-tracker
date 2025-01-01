package db

import (
	"fmt"
	"log"

	"github.com/alexshelto/tigres-tracker/db/models"
	"github.com/alexshelto/tigres-tracker/dto"
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

	log.Println("Database connected and tables migrated.")
}

func GetDB() *gorm.DB {
	return DB
}

func AddSongAndIncrementUser(db *gorm.DB, songName string, guildId string, requestedByID string, messageId string) error {
	song := models.Song{
		SongName:    songName,
		RequestedBy: requestedByID,
		MessageID:   messageId,
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

	return nil
}

// TopSongsByCount retrieves the top songs by the number of requests
func TopSongsByCount(db *gorm.DB, limit int) ([]models.Song, error) {
	var songs []models.Song

	// Query to get top songs by song count, ordered by the number of requests
	err := db.Model(&models.Song{}).
		Select("song_name, COUNT(*) as count").
		Group("song_name").
		Order("count DESC").
		Limit(limit).
		Find(&songs).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching top songs: %v", err)
	}

	return songs, nil
}

func TopSongsByUser(db *gorm.DB, userId uint, limit int) ([]models.Song, error) {
	var songs []models.Song

	err := db.Model(&models.Song{}).
		Select("song_name, COUNT(*) as count").
		Where("requested_by = ?", userId).
		Group("song_name").
		Order("count DESC").
		Limit(limit).
		Find(&songs).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching top songs by user: %v", err)
	}

	return songs, nil
}

func TopSongsInGuild(db *gorm.DB, guildId string, limit int) ([]dto.SongRequestCountDTO, error) {
	var songCounts []dto.SongRequestCountDTO

	err := db.Model(&models.Song{}).
		Select("song_name, COUNT(*) as count").
		Where("guild_id = ?", guildId).
		Group("song_name").
		Order("count DESC").
		Limit(limit).
		Find(&songCounts).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching top songs by guild: %v", err)
	}

	return songCounts, nil
}

func GetTotalSongsInGuild(db *gorm.DB, guildID string) (int64, error) {
	var count int64
	err := db.Table("songs").Where("guild_id = ?", guildID).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("could not count songs in guild %v", err)
	}
	return count, nil
}
