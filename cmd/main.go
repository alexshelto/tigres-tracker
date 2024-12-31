package main

import (
	"fmt"
	"log"

	"github.com/alexshelto/tigres-tracker/config"
	"github.com/alexshelto/tigres-tracker/db"
)

func main() {
	fmt.Println("Hello World")

	botConfig := config.LoadBotConfig()
	dbConfig := config.LoadDBConfig()

	log.Println("Bot token: ", botConfig.BotToken)
	log.Println("DB File: ", dbConfig.DatabaseFile)

	db.ConnectDB(dbConfig.DatabaseFile)
}
