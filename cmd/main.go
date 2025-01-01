package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"

	"log"

	"github.com/alexshelto/tigres-tracker/commands"
	"github.com/alexshelto/tigres-tracker/config"
	"github.com/alexshelto/tigres-tracker/db"
	"github.com/alexshelto/tigres-tracker/utils"
)

type Flags struct {
	ChannelID string
}

func main() {

	flags := Flags{}
	flag.StringVar(&flags.ChannelID, "channel", "", "Channel ID to hydrate messages from")
	flag.Parse()

	botConfig := config.LoadBotConfig()
	dbConfig := config.LoadDBConfig()

	db.ConnectDB(dbConfig.DatabaseFile)

	// Create new Discord Session
	dg, err := discordgo.New("Bot " + botConfig.BotToken)
	if err != nil {
		log.Fatal("Error creating Discord session: ", err)
	}

	if flags.ChannelID != "" {
		hydrateMessageHistory(dg, flags.ChannelID)
		return
	}

	// Register the messageCreate func as a fallback for MessageCreate events
	dg.AddHandler(messageCreate)

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

type ParsedSongInfo struct {
	Name        string
	RequestedBy string
}

// messageCreate is called whenever a new message is created
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	content := strings.ToLower(m.Content)

	if strings.HasPrefix(content, "t!help") {
		commands.HandleHelp(s, m)
	} else if strings.HasPrefix(content, "t!chart") {
		commands.HandleChart(s, m, db.GetDB())
	} else if strings.HasPrefix(content, "t!stats") {
		// Extract the user mention if provided
		line := strings.TrimSpace(m.Content[len("t!stats"):])
		var userId string
		if line == "" {
			userId = m.Author.ID
		} else {
			userId = strings.Trim(line, "<@>")
		}
		commands.HandleStats(s, m, userId, db.GetDB())
	}

	if len(m.Embeds) > 0 {
		for _, embed := range m.Embeds {
			handleEmbed(embed, m.GuildID, m.ID)
		}
	}
}

func handleEmbed(embed *discordgo.MessageEmbed, guildId string, messageId string) {
	songInfo := processEmbedDataForNowPlaying(embed)
	if songInfo != nil {
		log.Printf("Parsed song info: %+v\n", songInfo)

		err := db.AddSongAndIncrementUser(db.GetDB(), songInfo.Name, guildId, songInfo.RequestedBy, messageId)

		if err != nil {
			log.Printf("Error saving song info to DB: %v\n", err)
		} else {
			log.Printf(
				"saved song '%s' requested by user id '%s' in guild: '%s' with message ID '%s'",
				songInfo.Name, songInfo.RequestedBy, guildId, messageId,
			)
		}
	}

}

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

// Function to process embed data and extract song information
func processEmbedDataForNowPlaying(embed *discordgo.MessageEmbed) *ParsedSongInfo {
	if embed.Title == "" || !strings.EqualFold(strings.TrimSpace(embed.Title), "Now Playing") {
		return nil
	}

	description := embed.Description
	if description == "" {
		return nil
	}

	lines := strings.Split(description, "\n")
	if len(lines) < 2 {
		return nil
	}

	// Extract song name and requester info
	songStr := lines[0]
	requestedBy := lines[len(lines)-1]
	requestedByID := utils.ExtractUserID(requestedBy)

	if songStr != "" && requestedByID != "" {
		return &ParsedSongInfo{
			Name:        songStr,
			RequestedBy: requestedByID,
		}
	}
	return nil
}
