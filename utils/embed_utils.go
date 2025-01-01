package utils

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type ParsedSongInfo struct {
	Name        string
	RequestedBy string
}

// process embed message from Pancake Bot and extract song information
func ProcessEmbedDataForNowPlaying(embed *discordgo.MessageEmbed) *ParsedSongInfo {
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
	requestedByID := ExtractUserID(requestedBy)

	if songStr != "" && requestedByID != "" {
		return &ParsedSongInfo{
			Name:        songStr,
			RequestedBy: requestedByID,
		}
	}
	return nil
}
