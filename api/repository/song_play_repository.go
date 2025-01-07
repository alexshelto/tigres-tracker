package repository

import (
	"context"
	"errors"

	"github.com/alexshelto/tigres-tracker/api/dto"
	"github.com/alexshelto/tigres-tracker/api/models"
	"gorm.io/gorm"
)

type SongPlayRepository struct{}

func NewSongPlayRepository() *SongPlayRepository {
	return &SongPlayRepository{}
}

func (r *SongPlayRepository) AddOrUpdatePlay(ctx context.Context, userID uint, songID uint, guildID string) error {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return errors.New("unable to get db from context")
	}

	// Check if a play entry already exists for this user, song, and guild
	var play models.Play
	result := db.Where("user_id = ? AND song_id = ? AND guild_id = ?", userID, songID, guildID).First(&play)
	if result.Error == nil {
		// If it exists, increment the play count
		play.PlayCount++
		return db.Save(&play).Error
	}

	if result.Error == gorm.ErrRecordNotFound {
		// If no entry exists, create a new play entry
		newPlay := models.Play{
			UserID:    userID,
			SongID:    songID,
			GuildID:   guildID,
			PlayCount: 1, // Initial play count
		}
		return db.Create(&newPlay).Error
	}

	return result.Error
}

func (r *SongPlayRepository) GetTopSongsInGuild(ctx context.Context, guildID string, limit int) ([]dto.SongRequestCountDTO, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("unable to get db from context")
	}

	var results []dto.SongRequestCountDTO

	err := db.
		Table("songs").
		Joins("JOIN plays ON songs.id = plays.song_id").
		Where("plays.guild_id = ?", guildID).
		Group("songs.id").
		Select("songs.song_name, SUM(plays.play_count) as count").
		Order("count DESC").
		Limit(limit).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetTopSongsByUserInGuild gets the top songs played by a user in a guild based on play count
func (r *SongPlayRepository) GetTopSongsByUserInGuild(ctx context.Context, userID uint, guildID string, limit int) ([]dto.SongRequestCountDTO, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("unable to get db from context")
	}

	var results []dto.SongRequestCountDTO

	err := db.
		Table("songs").
		Joins("JOIN plays ON songs.id = plays.song_id").
		Where("plays.guild_id = ? AND plays.user_id = ?", guildID, userID).
		Group("songs.id").
		Select("songs.song_name, SUM(plays.play_count) as count").
		Order("count DESC").
		Limit(limit).
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *SongPlayRepository) GetTotalSongPlaysInGuild(ctx context.Context, guildID string) (int, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return 0, errors.New("unable to get db from context")
	}

	var totalPlays int
	err := db.
		Table("plays").
		Where("guild_id = ?", guildID).
		Select("SUM(play_count) as total").
		Scan(&totalPlays).Error

	if err != nil {
		return 0, err
	}

	return totalPlays, nil
}

func (r *SongPlayRepository) GetTotalUserSongPlaysInGuild(ctx context.Context, userID uint, guildID string) (int, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return 0, errors.New("unable to get db from context")
	}

	var userTotalPlays int
	err := db.
		Table("plays").
		Where("guild_id = ? AND user_id = ?", guildID, userID).
		Select("SUM(play_count) as total").
		Scan(&userTotalPlays).Error
	if err != nil {
		return 0, err
	}

	return userTotalPlays, nil
}
