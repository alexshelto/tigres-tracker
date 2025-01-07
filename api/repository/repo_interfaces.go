package repository

import (
	"context"

	"github.com/alexshelto/tigres-tracker/api/dto"
	"github.com/alexshelto/tigres-tracker/api/models"
)

// IUserRepository defines the interface for user-related operations
type IUserRepository interface {
	GetOrCreateUser(ctx context.Context, discordID string) (*models.User, error)
}

// ISongRepository defines the interface for song-related operations
type ISongRepository interface {
	GetOrCreateSong(ctx context.Context, guildID, songName string) (*models.Song, error)
}

// ISongPlayRepository defines the interface for song-play-related operations
type ISongPlayRepository interface {
	AddOrUpdatePlay(ctx context.Context, userId uint, songID uint, guildID string) error
	GetTopSongsInGuild(ctx context.Context, guildID string, limit int) ([]dto.SongRequestCountDTO, error)
	GetTopSongsByUserInGuild(ctx context.Context, userID uint, guildID string, limit int) ([]dto.SongRequestCountDTO, error)
	GetTotalSongPlaysInGuild(ctx context.Context, guildID string) (int, error)
	GetTotalUserSongPlaysInGuild(ctx context.Context, userID uint, guildID string) (int, error)
}
