package db

import (
	"testing"

	"github.com/alexshelto/tigres-tracker/db"
	"github.com/alexshelto/tigres-tracker/db/models"
	"github.com/stretchr/testify/assert"
)

func TestAddSongAndIncrementUser(t *testing.T) {
	// Setup a test database (in-memory or temporary file)
	db.ConnectDB(":memory:") // For in-memory DB

	// Ensure tables are created in the test DB
	db.GetDB().AutoMigrate(&models.User{}, &models.Song{})

	// Test data
	songName := "Test Song"
	guildId := uint(999)
	requestedByID := uint(1)

	// Call AddSongAndIncrementUser function
	err := db.AddSongAndIncrementUser(db.GetDB(), songName, guildId, requestedByID)
	assert.Nil(t, err)

	var user models.User
	db.GetDB().First(&user, "discord_id = ?", requestedByID)

	// Check that the song count is incremented
	assert.Equal(t, 1, user.SongCount)

	// Fetch the song from the DB
	var song models.Song
	db.GetDB().First(&song, "song_name = ?", songName)

	// Check that the song was added with the correct data
	assert.Equal(t, songName, song.SongName)
	assert.Equal(t, requestedByID, song.RequestedBy)

	// Add another song, check incremented
	// Call AddSongAndIncrementUser function
	err = db.AddSongAndIncrementUser(db.GetDB(), songName, guildId, requestedByID)
	assert.Nil(t, err)

	db.GetDB().First(&user, "discord_id = ?", requestedByID)
	// Check that the song count is incremented
	assert.Equal(t, 2, user.SongCount)
}
