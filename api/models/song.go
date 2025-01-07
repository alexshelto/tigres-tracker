package models

import "gorm.io/gorm"

type Song struct {
	gorm.Model

	SongName string `gorm:"not null;index:song_guild,unique"`
	GuildID  string `gorm:"not null;index:song_guild,unique"` // index for faster lookups
}
