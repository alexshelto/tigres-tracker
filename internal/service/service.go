package service

import (
	"fmt"
	"log"
	"strings"

	"github.com/alexshelto/tigres-tracker/commands"
	"github.com/alexshelto/tigres-tracker/internal/client"
	"github.com/alexshelto/tigres-tracker/utils"
	"github.com/bwmarrin/discordgo"
)

type MessageService struct {
	Client *client.APIClient
}

func NewMessageService(client *client.APIClient) *MessageService {
	return &MessageService{
		Client: client,
	}
}

// HandleMessage is the main method to handle messages
func (s *MessageService) HandleMessage(ses *discordgo.Session, m *discordgo.MessageCreate) {
	// Here you can use s.Client to access the client or other services.
	if m.Author.ID == ses.State.User.ID {
		return // Ignore bot's own messages
	}

	if utils.IsFromPancakeBot(m.Author.ID) {
		songsToPost := ProcessNowPlayingMessageFromPancakeBot(m)
		for _, request := range songsToPost {
			_, err := s.Client.PostSongPlay(request.RequestedBy, request.Name, request.GuildID)
			if err != nil {
				log.Fatalf("Failed to POST new song: %+v", request)
			}
		}
	}

	s.HandleCommands(ses, m)
}

func (s *MessageService) HandleCommands(ses *discordgo.Session, m *discordgo.MessageCreate) {
	content := strings.ToLower(m.Content)

	switch {
	case strings.HasPrefix(content, "t!help"):
		commands.HandleHelp(ses, m)
	case strings.HasPrefix(content, "t!chart"):
		s.HandleChart(ses, m)
		/*
			case strings.HasPrefix(content, "t!stats"):
				HandleStatsCommand(s, m)
		*/
	}
}

func (s *MessageService) HandleChart(ses *discordgo.Session, m *discordgo.MessageCreate) {
	guildID := m.GuildID
	topSongs, err := s.Client.GetTopSongsInGuild(guildID, 10)

	if err != nil {
		log.Println("Error geting top songs in guild: ", err)
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

	_, err = ses.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		log.Println("Error sending embed: ", err)
	}

}

func ProcessNowPlayingMessageFromPancakeBot(m *discordgo.MessageCreate) []*utils.ParsedSongInfo {
	var songInfos []*utils.ParsedSongInfo

	if len(m.Embeds) > 0 {
		for _, embed := range m.Embeds {
			parsedInfo := handleEmbeddedNowPlaying(embed, m.GuildID)
			if parsedInfo != nil {
				songInfos = append(songInfos, parsedInfo)
			}
		}
	}
	return songInfos
}

func handleEmbeddedNowPlaying(embed *discordgo.MessageEmbed, guildID string) *utils.ParsedSongInfo {
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
		return &utils.ParsedSongInfo{
			Name:        songStr,
			RequestedBy: requestedByID,
			GuildID:     guildID,
		}
	}
	return nil
}
