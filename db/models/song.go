package models

import "gorm.io/gorm"

type Song struct {
	gorm.Model

	SongName    string
	RequestedBy uint // Foreign key to User table (User id)
	GuildId     uint // Discord Server (Guild) ID where song was requested
}

func (song *Song) CreateSong(db *gorm.DB) error {
	if err := db.Create(&song).Error; err != nil {
		return err
	}
	return nil
}
