package commands

/*
import (
	"fmt"
	"github.com/alexshelto/tigres-tracker/db"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"log"
)

func HandleStats(s *discordgo.Session, m *discordgo.MessageCreate, userId string, database *gorm.DB) {
	guildID := m.GuildID

	userName, err := s.User(userId)
	if err != nil {
		log.Printf("Failed to get username for id: '%s' | %v", userId, err)
	}

	userStats, err := db.GetUserStats(database, userId, guildID)
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
*/
