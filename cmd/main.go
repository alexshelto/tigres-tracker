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
		//hydrateMessageHistory(dg, flags.ChannelID)
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

/*
func hydrateMessageHistory(s *discordgo.Session, channelID string) {
	channel, err := s.Channel(channelID)
	if err != nil {
		log.Fatalf("error fetching channel: %v", err)
		return
	}

	var lastMessageID string

	for {
		messages, err := s.ChannelMessages(channelID, 100, lastMessageID, "", "")
		if err != nil {
			log.Fatalf("error fetching messages: %v", err)
			return
		}

		if len(messages) == 0 {
			log.Println("No more messages to process.")
			break
		}

		// Process each message
		for _, m := range messages {
			if len(m.Embeds) > 0 {
				for _, embed := range m.Embeds {
					handleEmbed(embed, channel.GuildID, m.ID)
				}
			}
		}
		lastMessageID = messages[len(messages)-1].ID
	}
}
*/
