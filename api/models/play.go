package models

import "gorm.io/gorm"

type Play struct {
	gorm.Model

	UserID    uint   `gorm:"not null;index:user_song_guild,unique;foreignKey:DiscordID;constraint:OnDelete:CASCADE"` // Foreign Key to User
	SongID    uint   `gorm:"not null;index:user_song_guild,unique;foreignKey:MessageID;constraint:OnDelete:CASCADE"` // Foreign Key to Song
	GuildID   string `json:"guild_id" gorm:"not null;index:user_song_guild,unique;"`                                 // Indexed for querying by guild
	PlayCount int    `json:"play_count" gorm:"not null;default:1"`
}
