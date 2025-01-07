package repository

import (
	"context"
	"errors"

	"github.com/alexshelto/tigres-tracker/api/models"
	"gorm.io/gorm"
)

type SongRepository struct{}

func NewSongRepository() *SongRepository {
	return &SongRepository{}
}

// Implementing ISongRepository Interface methods
func (r *SongRepository) GetOrCreateSong(ctx context.Context, guildID, songName string) (*models.Song, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("unable to get db from context")
	}

	var song models.Song

	// Check if the song already exists for this guild_id and song_name
	result := db.Where("guild_id = ? AND song_name = ?", guildID, songName).First(&song)
	if result.Error == nil {
		return &song, nil
	}

	// If not found, create a new song
	if result.Error == gorm.ErrRecordNotFound {
		song = models.Song{
			GuildID:  guildID,
			SongName: songName,
		}

		if err := db.Create(&song).Error; err != nil {
			return nil, err
		}

		return &song, nil
	}

	return nil, result.Error
}
