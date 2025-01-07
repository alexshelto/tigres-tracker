package repository

import (
	"context"
	"errors"

	"github.com/alexshelto/tigres-tracker/api/models"
	"gorm.io/gorm"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// GetOrCreateUser retrieves a user by Discord ID or creates a new one if not found
func (r *UserRepository) GetOrCreateUser(ctx context.Context, discordID string) (*models.User, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("unable to get db from context")
	}

	var user models.User
	// Check if the user exists
	err := db.Where("discord_id = ?", discordID).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		// If not found, create a new user
		user = models.User{DiscordID: discordID}
		if err := db.Create(&user).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}
