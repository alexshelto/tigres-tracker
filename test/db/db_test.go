package db

import (
	"strconv"
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
	err := db.AddSongAndIncrementUser(db.GetDB(), songName, guildId, requestedByID, "1")
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
	err = db.AddSongAndIncrementUser(db.GetDB(), songName, guildId, requestedByID, "2")
	assert.Nil(t, err)

	db.GetDB().First(&user, "discord_id = ?", requestedByID)
	// Check that the song count is incremented
	assert.Equal(t, 2, user.SongCount)
}

func TestTopSongsByCount(t *testing.T) {
	// Setup a test database (in-memory or temporary file)
	db.ConnectDB(":memory:") // For in-memory DB

	// Ensure tables are created in the test DB
	db.GetDB().AutoMigrate(&models.User{}, &models.Song{})

	// Add some songs and increment song counts
	err := db.AddSongAndIncrementUser(db.GetDB(), "Song1", 999, 12345, "1")
	assert.NoError(t, err)
	err = db.AddSongAndIncrementUser(db.GetDB(), "Song1", 999, 12345, "2")
	assert.NoError(t, err)
	err = db.AddSongAndIncrementUser(db.GetDB(), "Song2", 999, 123456, "3")
	assert.NoError(t, err)

	// Fetch top songs by song count
	songs, err := db.TopSongsByCount(db.GetDB(), 5)
	assert.NoError(t, err)
	assert.Len(t, songs, 2, "Expected two songs to be returned")
	assert.Equal(t, "Song1", songs[0].SongName, "Expected Song1 to be the most requested")
	assert.Equal(t, "Song2", songs[1].SongName, "Expected Song2 to be the second most requested")
}

func TestTopSongsByUser(t *testing.T) {
	db.ConnectDB(":memory:") // For in-memory DB

	// Ensure tables are created in the test DB
	db.GetDB().AutoMigrate(&models.User{}, &models.Song{})

	// Add some songs and increment song counts
	for i := 0; i < 10; i++ {
		id := strconv.Itoa(i)
		err := db.AddSongAndIncrementUser(db.GetDB(), "SongA", 999, 1, id)
		assert.NoError(t, err)
	}
	for i := 0; i < 15; i++ {
		id := strconv.Itoa(10 + i)
		err := db.AddSongAndIncrementUser(db.GetDB(), "SongB", 999, 1, id)
		assert.NoError(t, err)
	}
	for i := 0; i < 5; i++ {
		id := strconv.Itoa(100 + i)
		err := db.AddSongAndIncrementUser(db.GetDB(), "SongC", 999, 1, id)
		assert.NoError(t, err)
	}

	// Different user
	for i := 0; i < 20; i++ {
		id := strconv.Itoa(1000 + i)
		err := db.AddSongAndIncrementUser(db.GetDB(), "Song1", 999, 2, id)
		assert.NoError(t, err)
	}

	songs, err := db.TopSongsByUser(db.GetDB(), 1, 4)
	assert.NoError(t, err)
	assert.Len(t, songs, 3, "Expected two songs to be returned")
	assert.Equal(t, "SongB", songs[0].SongName, "Expected songB to be most requested")
	assert.Equal(t, "SongA", songs[1].SongName, "Expected songA to be second most requested")
	assert.Equal(t, "SongC", songs[2].SongName, "Expected songC to be second most requested")

}
