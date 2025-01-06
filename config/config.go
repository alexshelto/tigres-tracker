package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type BotConfig struct {
	BotToken string
}

type DBConfig struct {
	DatabaseFile string
}

func LoadDBConfig() *DBConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: '%v'", err)
	}

	databaseFile := os.Getenv("DATABASE_FILE")

	if databaseFile == "" {
		log.Fatal("Database file is not set. Please provide DATABASE_FILE environment variable")
	}

	return &DBConfig{
		DatabaseFile: databaseFile,
	}
}

func LoadBotConfig() *BotConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: '%v'", err)
	}

	botToken := os.Getenv("BOT_TOKEN")

	if botToken == "" {
		log.Fatal("Bot token is not set. Please provide BOT_TOKEN environment variable")
	}

	return &BotConfig{
		BotToken: botToken,
	}

}
