package repository

import (
	"context"
	"testing"

	"github.com/alexshelto/tigres-tracker/api/models"
	"github.com/alexshelto/tigres-tracker/api/testutils"
	"github.com/stretchr/testify/assert"
)

var (
	songName1 = "Rooster"
	songName2 = "Higher"
	songName3 = "Highwayman"

	guildID_1 = "9000"
	guildID_2 = "9001"

	userID_1 = "1000"
	userID_2 = "1001"
	userID_3 = "1002"
	userID_4 = "1003"

	user_1 = models.User{DiscordID: userID_1}
	user_2 = models.User{DiscordID: userID_2}

	song_1_guild_1 = models.Song{SongName: songName1, GuildID: guildID_1}
	song_2_guild_1 = models.Song{SongName: songName2, GuildID: guildID_1}
	song_3_guild_1 = models.Song{SongName: songName3, GuildID: guildID_1}

	song_1_guild_2 = models.Song{SongName: songName1, GuildID: guildID_2}
	song_2_guild_2 = models.Song{SongName: songName2, GuildID: guildID_2}
)

func TestSongPlayRepository_GetTopSongsInGuild(t *testing.T) {
	db, sqlDB := testutils.SetupTestDB(t)
	defer sqlDB.Close()

	ctx := context.WithValue(context.Background(), "db", db)
	repo := &SongPlayRepository{}

	// User 1 Guild 1
	testutils.PopulateDBWithUserAndSong(t, db, &user_1, &song_1_guild_1, 100)
	testutils.PopulateDBWithUserAndSong(t, db, &user_1, &song_2_guild_1, 200)
	testutils.PopulateDBWithUserAndSong(t, db, &user_1, &song_3_guild_1, 300)

	// User 2 Guild 1
	testutils.PopulateDBWithUserAndSong(t, db, &user_2, &song_1_guild_1, 1000)

	// User 1 Guild 2
	testutils.PopulateDBWithUserAndSong(t, db, &user_1, &song_1_guild_2, 100)

	stats, err := repo.GetTopSongsInGuild(ctx, guildID_1, 2)
	assert.Nil(t, err)
	assert.NotNil(t, stats)

	assert.Equal(t, 2, len(stats))

	assert.Equal(t, 1100, stats[0].Count)
	assert.Equal(t, songName1, stats[0].SongName)

	assert.Equal(t, 300, stats[1].Count)
	assert.Equal(t, songName3, stats[1].SongName)
}

func TestSongPlayRepository_GetTopSongsInGuild_requestMoreThanExists(t *testing.T) {
	db, sqlDB := testutils.SetupTestDB(t)
	defer sqlDB.Close()

	ctx := context.WithValue(context.Background(), "db", db)
	repo := &SongPlayRepository{}

	// User 1 Guild 1
	testutils.PopulateDBWithUserAndSong(t, db, &user_1, &song_2_guild_1, 200)
	testutils.PopulateDBWithUserAndSong(t, db, &user_1, &song_3_guild_1, 300)

	stats, err := repo.GetTopSongsInGuild(ctx, guildID_1, 5)
	assert.Nil(t, err)
	assert.NotNil(t, stats)

	assert.Equal(t, 2, len(stats))

	assert.Equal(t, 300, stats[0].Count)
	assert.Equal(t, songName3, stats[0].SongName)

	assert.Equal(t, 200, stats[1].Count)
	assert.Equal(t, songName2, stats[1].SongName)
}

func TestSongPlayRepository_GetTopSongsByUserInGuild(t *testing.T) {
	db, sqlDB := testutils.SetupTestDB(t)
	defer sqlDB.Close()

	ctx := context.WithValue(context.Background(), "db", db)
	repo := &SongPlayRepository{}
	userRepo := &UserRepository{}

	// User 1 Guild 1
	testutils.PopulateDBWithUserAndSong(t, db, &user_1, &song_1_guild_1, 5)
	testutils.PopulateDBWithUserAndSong(t, db, &user_1, &song_2_guild_1, 10)

	// User 2 Guild 1
	testutils.PopulateDBWithUserAndSong(t, db, &user_2, &song_1_guild_1, 1000)
	testutils.PopulateDBWithUserAndSong(t, db, &user_2, &song_2_guild_1, 8000)

	user1_id, _ := userRepo.GetOrCreateUser(ctx, user_1.DiscordID)

	stats, err := repo.GetTopSongsByUserInGuild(ctx, user1_id.ID, guildID_1, 2)

	assert.Nil(t, err)
	assert.NotNil(t, stats)

	assert.Equal(t, 2, len(stats))

	assert.Equal(t, 10, stats[0].Count)
	assert.Equal(t, songName2, stats[0].SongName)

	assert.Equal(t, 5, stats[1].Count)
	assert.Equal(t, songName1, stats[1].SongName)
}

func TestSongPlayRepository_GetTotalSongsInGuild(t *testing.T) {
	db, sqlDB := testutils.SetupTestDB(t)
	defer sqlDB.Close()

	ctx := context.WithValue(context.Background(), "db", db)
	repo := &SongPlayRepository{}

	// User 1 Guild 1
	testutils.PopulateDBWithUserAndSong(t, db, &user_1, &song_1_guild_1, 5)
	testutils.PopulateDBWithUserAndSong(t, db, &user_1, &song_2_guild_1, 10)

	// User 2 Guild 1
	testutils.PopulateDBWithUserAndSong(t, db, &user_2, &song_1_guild_1, 1000)
	testutils.PopulateDBWithUserAndSong(t, db, &user_2, &song_2_guild_1, 2000)

	// User 2 Guild 2
	testutils.PopulateDBWithUserAndSong(t, db, &user_2, &song_1_guild_2, 1000)
	testutils.PopulateDBWithUserAndSong(t, db, &user_2, &song_2_guild_2, 2000)

	count, err := repo.GetTotalSongPlaysInGuild(ctx, guildID_1)

	assert.Nil(t, err)
	assert.Equal(t, count, 3015)
}

func TestSongPlayRepository_GetTotalUserSongPlaysInGuild(t *testing.T) {
	db, sqlDB := testutils.SetupTestDB(t)
	defer sqlDB.Close()

	ctx := context.WithValue(context.Background(), "db", db)
	repo := &SongPlayRepository{}
	userRepo := &UserRepository{}

	// User 1 Guild 1
	testutils.PopulateDBWithUserAndSong(t, db, &user_1, &song_1_guild_1, 5)
	testutils.PopulateDBWithUserAndSong(t, db, &user_1, &song_2_guild_1, 10)

	// User 2 Guild 1
	testutils.PopulateDBWithUserAndSong(t, db, &user_2, &song_1_guild_1, 1000)
	testutils.PopulateDBWithUserAndSong(t, db, &user_2, &song_2_guild_1, 8000)

	user1_id, _ := userRepo.GetOrCreateUser(ctx, user_1.DiscordID)

	count, err := repo.GetTotalUserSongPlaysInGuild(ctx, user1_id.ID, guildID_1)

	assert.Nil(t, err)
	assert.Equal(t, count, 15)
}
