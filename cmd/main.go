package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"

	"log"

	"github.com/alexshelto/tigres-tracker/config"
	"github.com/alexshelto/tigres-tracker/db"
)

func main() {
	botConfig := config.LoadBotConfig()
	dbConfig := config.LoadDBConfig()

	db.ConnectDB(dbConfig.DatabaseFile)

	// Create new Discord Session
	dg, err := discordgo.New("Bot " + botConfig.BotToken)
	if err != nil {
		log.Fatal("Error creating Discord session: ", err)
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
	RequestedBy uint
}

// messageCreate is called whenever a new message is created
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if len(m.Embeds) > 0 {
		for _, embed := range m.Embeds {
			songInfo := processEmbedDataForNowPlaying(embed)
			if songInfo != nil {
				fmt.Printf("Parsed song info: %+v\n", songInfo)
				fmt.Println("Guild id: ", m.GuildID)
				fmt.Println("id: ", m.ID)
			}
		}
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
	requestedByID := extractUserID(requestedBy)

	if songStr != "" && requestedByID != 0 {
		return &ParsedSongInfo{
			Name:        songStr,
			RequestedBy: requestedByID,
		}
	}
	return nil
}

// Function to extract user ID from a string
func extractUserID(requestString string) uint {
	re := regexp.MustCompile(`<@(\d+)>`)
	matches := re.FindStringSubmatch(requestString)
	if len(matches) > 1 {
		// Convert to uint
		var userID uint
		fmt.Sscanf(matches[1], "%d", &userID)
		return userID
	}
	return 0
}
