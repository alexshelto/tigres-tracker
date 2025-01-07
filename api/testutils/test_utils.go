package testutils

import (
	"database/sql"
	"log"
	"testing"

	"github.com/alexshelto/tigres-tracker/api/models"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func PopulateDBWithUserAndSong(t *testing.T, db *gorm.DB, user *models.User, song *models.Song, plays int) {
	err := db.FirstOrCreate(user, models.User{DiscordID: user.DiscordID}).Error
	require.NoError(t, err)
	require.NotNil(t, user.ID)
	log.Printf("After adding User: %+v", user)

	err = db.FirstOrCreate(song, models.Song{SongName: song.SongName, GuildID: song.GuildID}).Error
	require.NoError(t, err)

	play := &models.Play{
		UserID:    user.ID,
		SongID:    song.ID,
		GuildID:   song.GuildID,
		PlayCount: plays,
	}

	err = db.Create(play).Error
	require.NoError(t, err)
}

func SetupTestDB(t *testing.T) (*gorm.DB, *sql.DB) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	sqlDB, err := db.DB()

	err = db.AutoMigrate(&models.User{}, &models.Song{}, &models.Play{})
	require.NoError(t, err)

	return db, sqlDB
}
