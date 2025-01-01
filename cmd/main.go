package main

import (
	"flag"
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
		handleHelp(s, m)
	} else if strings.HasPrefix(content, "t!chart") {
		handleChart(s, m)
	} else if strings.HasPrefix(content, "t!stats") {
		// Extract the user mention if provided
		line := strings.TrimSpace(m.Content[len("t!stats"):])
		var userId string
		if line == "" {
			userId = m.Author.ID
		} else {
			userId = strings.Trim(line, "<@>")
		}
		handleStats(s, m, userId)
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
	requestedByID := extractUserID(requestedBy)

	if songStr != "" && requestedByID != "" {
		return &ParsedSongInfo{
			Name:        songStr,
			RequestedBy: requestedByID,
		}
	}
	return nil
}

// Function to extract user ID from a string
func extractUserID(requestString string) string {
	re := regexp.MustCompile(`<@(\d+)>`)
	matches := re.FindStringSubmatch(requestString)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Format the response
	response := `**Commands:**
	t!stats [@user] - Show the top 5 songs a user has queued and how many songs they've queued. If no user is mentioned, it shows stats for yourself.
	t!chart - Show the top 10 songs requested in the server and how many songs have been queued in total.
	t!help - Show this help message.`

	// Send the response to the channel
	s.ChannelMessageSend(m.ChannelID, response)
}

// Handle the t!chart command
func handleChart(s *discordgo.Session, m *discordgo.MessageCreate) {
	guildID := m.GuildID

	// Query the top 10 songs in the guild
	topSongs, err := db.TopSongsInGuild(db.GetDB(), guildID, 10) // Replace with your actual query logic
	if err != nil {
		log.Println("Error geting top songs in guild: ", err)
		return
	}

	totalSongs, err := db.GetTotalSongsInGuild(db.GetDB(), guildID)
	if err != nil {
		log.Println("Error geting total songs in guild: ", err)
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: "Top 10 Songs in the Server",
		Color: 0x00FF00,
	}
	for _, song := range topSongs {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   song.SongName,
			Value:  fmt.Sprintf("Plays: %d", song.Count),
			Inline: false,
		})
	}

	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("Total songs played in server: %d\n", totalSongs),
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		log.Println("Error sending embed: ", err)
	}
}

func handleStats(s *discordgo.Session, m *discordgo.MessageCreate, userId string) {
	guildID := m.GuildID

	userName, err := s.User(userId)
	if err != nil {
		log.Printf("Failed to get username for id: '%s' | %v", userId, err)
	}

	userStats, err := db.GetUserStats(db.GetDB(), userId, guildID)
	if err != nil {
		log.Println("Error getting user stats: ", err)
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("Top 5 Songs %s requested", userName),
		Color: 0x00FF00,
	}

	for _, song := range userStats.TopSongs {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   song.SongName,
			Value:  fmt.Sprintf("requests: %d", song.Count),
			Inline: false,
		})
	}

	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("Total songs requested in server: %d\n", userStats.TotalSongs),
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		log.Println("Error sending embed: ", err)
	}
}
