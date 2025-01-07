package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"

	"log"

	"github.com/alexshelto/tigres-tracker/config"
	"github.com/alexshelto/tigres-tracker/internal/client"
	"github.com/alexshelto/tigres-tracker/internal/handler"
)

type Flags struct {
	ChannelID string
}

func main() {
	flags := Flags{}
	flag.StringVar(&flags.ChannelID, "channel", "", "Channel ID to hydrate messages from")
	flag.Parse()

	botConfig := config.LoadBotConfig()
	clientConfig := config.LoadClientConfig()

	// Create new Discord Session
	dg, err := discordgo.New("Bot " + botConfig.BotToken)
	if err != nil {
		log.Fatal("Error creating Discord session: ", err)
	}

	client := client.NewClient(clientConfig)
	handler.InitHandlers(dg, client)

	if flags.ChannelID != "" {
		handler.HydrateHistory(dg, flags.ChannelID)
		return
	}

	dg.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = dg.Open()
	if err != nil {
		log.Fatal("error opening connection: ", err)
	}

	// Wait till ctrl-c or other term signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc

	defer dg.Close()
}
