package models

import "gorm.io/gorm"

// User represents a Discord user in the database
type User struct {
	gorm.Model

	ID        uint `gorm:"primary_key"`
	DiscordID uint `gorm:"unique;not null"` // Discord User ID
	SongCount int  `gorm:"default:0"`
}

func (user *User) IncrementSongCount(db *gorm.DB) error {
	user.SongCount++
	if err := db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

// Get or create a user by their DiscordID
func GetOrCreateUser(db *gorm.DB, discordID uint) (*User, error) {
	var user User
	if err := db.Where("discord_id = ?", discordID).First(&user).Error; err != nil {
		// If user doesn't exist, create one
		user = User{DiscordID: discordID}
		if err := db.Create(&user).Error; err != nil {
			return nil, err
		}
	}
	return &user, nil
}
