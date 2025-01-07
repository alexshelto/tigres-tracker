package handler

import (
	"github.com/alexshelto/tigres-tracker/internal/client"
	"github.com/alexshelto/tigres-tracker/internal/service"
	"github.com/bwmarrin/discordgo"
)

var messageService *service.MessageService

func InitHandlers(sess *discordgo.Session, client *client.APIClient) {
	messageService = service.NewMessageService(client)

	// Register the messageCreate handler
	sess.AddHandler(messageCreate)
}

// This is the messageCreate handler that will delegate the logic to the service
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Call the service to handle the message
	messageService.HandleMessage(s, m)
}

func HydrateHistory(s *discordgo.Session, channelID string) {
	messageService.HydratePancakeNowPlayingHistoryFromChannelID(s, channelID)
}
