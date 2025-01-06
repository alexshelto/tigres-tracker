package service

import (
	"log"
	"strings"

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
