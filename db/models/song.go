package models

import "gorm.io/gorm"

type Song struct {
	gorm.Model

	SongName    string
	RequestedBy uint   // Foreign key to User table (User id)
	MessageID   string `gorm:"unique;not null"` // Unique message ID for deduplication
	GuildId     uint   // Discord Server (Guild) ID where song was requested
}

func (song *Song) CreateSong(db *gorm.DB) error {
	if err := db.Create(&song).Error; err != nil {
		return err
	}
	return nil
}
