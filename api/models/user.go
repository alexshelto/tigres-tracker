package models

import "gorm.io/gorm"

// User represents a Discord user in the database
type User struct {
	gorm.Model

	DiscordID string `json:"discord_id" gorm:"unique;not null"` // Discord User ID
}
