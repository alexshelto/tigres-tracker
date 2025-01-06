package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type BotConfig struct {
	BotToken string
}

type ClientConfig struct {
	BaseURL string
}

func LoadBotConfig() BotConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: '%v'", err)
	}

	botToken := os.Getenv("BOT_TOKEN")

	if botToken == "" {
		log.Fatal("Bot token is not set. Please provide BOT_TOKEN environment variable")
	}

	return BotConfig{
		BotToken: botToken,
	}
}

func LoadClientConfig() ClientConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: '%v'", err)
	}

	host := os.Getenv("BASE_URL")

	if host == "" {
		log.Fatal("Bot token is not set. Please provide BOT_TOKEN environment variable")
	}

	return ClientConfig{
		BaseURL: host,
	}
}
