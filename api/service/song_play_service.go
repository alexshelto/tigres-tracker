package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/alexshelto/tigres-tracker/api/dto"
	"github.com/alexshelto/tigres-tracker/api/models"
	"github.com/alexshelto/tigres-tracker/api/repository"
)

type SongPlayService struct {
	userRepo     repository.IUserRepository
	songRepo     repository.ISongRepository
	songPlayRepo repository.ISongPlayRepository
}

func NewSongPlayService(
	userRepo repository.IUserRepository,
	songRepo repository.ISongRepository,
	songPlayRepo repository.ISongPlayRepository,
) *SongPlayService {
	return &SongPlayService{userRepo, songRepo, songPlayRepo}
}

func (s *SongPlayService) AddSongPlay(ctx context.Context, request models.SongPlayRequest) error {
	user, err := s.userRepo.GetOrCreateUser(ctx, request.UserID)
	if err != nil {
		return errors.New(fmt.Sprintf("could not fetch user %s | %v", request.UserID, err))
	}

	song, err := s.songRepo.GetOrCreateSong(ctx, request.GuildID, request.SongName)
	if err != nil {
		return errors.New(fmt.Sprintf("could not create song %s | %v", request.SongName, err))
	}

	err = s.songPlayRepo.AddOrUpdatePlay(ctx, user.ID, song.ID, request.GuildID)
	if err != nil {
		return errors.New("could not create play")
	}

	return nil
}

func (s *SongPlayService) GetTopSongsInGuild(ctx context.Context, guildID string, count int) ([]dto.SongRequestCountDTO, error) {
	stats, err := s.songPlayRepo.GetTopSongsInGuild(ctx, guildID, count)

	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (s *SongPlayService) GetTopSongsByUserInGuild(ctx context.Context, discordID string, guildID string, count int) ([]dto.SongRequestCountDTO, error) {
	user, err := s.userRepo.GetOrCreateUser(ctx, discordID)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not fetch user %s | %v", discordID, err))
	}

	stats, err := s.songPlayRepo.GetTopSongsByUserInGuild(ctx, user.ID, guildID, count)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (s *SongPlayService) GetTotalSongPlaysInGuild(ctx context.Context, guildID string) (int, error) {
	total, err := s.songPlayRepo.GetTotalSongPlaysInGuild(ctx, guildID)
	if err != nil {
		return 0, errors.New("failed to query total songs in guild")
	}

	return total, nil
}

func (s *SongPlayService) GetTotalUserSongPlaysInGuild(ctx context.Context, discordID string, guildID string) (int, error) {
	user, err := s.userRepo.GetOrCreateUser(ctx, discordID)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("could not fetch user %s | %v", discordID, err))
	}

	total, err := s.songPlayRepo.GetTotalUserSongPlaysInGuild(ctx, user.ID, guildID)
	if err != nil {
		return 0, errors.New("failed to fetch users song total")
	}

	return total, nil
}
