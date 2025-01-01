package commands

import (
	"github.com/bwmarrin/discordgo"
)

func HandleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Format the response
	// TODO: Embedded message
	response := `**Commands:**
	t!stats [@user] - Show the top 5 songs a user has queued and how many songs they've queued. If no user is mentioned, it shows stats for yourself.
	t!chart - Show the top 10 songs requested in the server and how many songs have been queued in total.
	t!help - Show this help message.`

	// Send the response to the channel
	s.ChannelMessageSend(m.ChannelID, response)
}
