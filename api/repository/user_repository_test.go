package repository

import (
	"context"
	"testing"

	"github.com/alexshelto/tigres-tracker/api/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_GetOrCreateUser(t *testing.T) {
	db, sqlDB := testutils.SetupTestDB(t)
	defer sqlDB.Close()

	ctx := context.WithValue(context.Background(), "db", db)

	userRepo := &UserRepository{}

	discordID := "1234"

	user, err := userRepo.GetOrCreateUser(ctx, discordID)

	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, discordID, user.DiscordID)

	// User Now already exists, should retrieve existing user
	user, err = userRepo.GetOrCreateUser(ctx, discordID)

	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, discordID, user.DiscordID)
}

func TestUserRepository_GetOrCreateUser_NoDBInContext(t *testing.T) {
	_, sqlDB := testutils.SetupTestDB(t)
	defer sqlDB.Close()

	// Create a new repository instance
	repo := &UserRepository{}

	// Define the Discord ID for testing
	discordID := "123456"

	// Test Case: No DB in context, should return an error
	ctx := context.Background() // Context without DB
	user, err := repo.GetOrCreateUser(ctx, discordID)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, user)
}
