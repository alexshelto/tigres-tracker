package repository

import (
	"context"
	"testing"

	"github.com/alexshelto/tigres-tracker/api/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSongRepository_GetOrCreateSong(t *testing.T) {

	db, sqlDB := testutils.SetupTestDB(t)
	defer sqlDB.Close()

	ctx := context.WithValue(context.Background(), "db", db)

	songRepo := &SongRepository{}

	guildID := "1234"
	songName := "Rooster"

	song, err := songRepo.GetOrCreateSong(ctx, guildID, songName)

	require.NoError(t, err)
	assert.NotNil(t, song)
	assert.Equal(t, guildID, song.GuildID)
	assert.Equal(t, songName, song.SongName)

	// Song Now already exists, should retrieve existing Song
	song, err = songRepo.GetOrCreateSong(ctx, guildID, songName)

	require.NoError(t, err)
	assert.NotNil(t, song)
	assert.Equal(t, guildID, song.GuildID)
	assert.Equal(t, songName, song.SongName)
}

func TestSongRepository_GetOrCreateSong_SameSongDifferentGuild(t *testing.T) {

	db, sqlDB := testutils.SetupTestDB(t)
	defer sqlDB.Close()

	ctx := context.WithValue(context.Background(), "db", db)

	songRepo := &SongRepository{}

	guildID_1 := "1234"
	guildID_2 := "1234"
	songName := "Rooster"

	song, err := songRepo.GetOrCreateSong(ctx, guildID_1, songName)

	require.NoError(t, err)
	assert.NotNil(t, song)
	assert.Equal(t, guildID_1, song.GuildID)
	assert.Equal(t, songName, song.SongName)

	// Now add song entry for different guild, should allow
	song, err = songRepo.GetOrCreateSong(ctx, guildID_2, songName)

	require.NoError(t, err)
	assert.NotNil(t, song)
	assert.Equal(t, guildID_2, song.GuildID)
	assert.Equal(t, songName, song.SongName)
}

func TestSongRepository_GetOrCreateSong_NoDBInContext(t *testing.T) {
	_, sqlDB := testutils.SetupTestDB(t)
	defer sqlDB.Close()

	// Test Case: No DB in context, should return an error
	ctx := context.Background() // Context without DB

	songRepo := &SongRepository{}

	guildID := "1234"
	songName := "Rooster"

	// Test Case: No DB in context, should return an error
	song, err := songRepo.GetOrCreateSong(ctx, guildID, songName)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, song)
}
