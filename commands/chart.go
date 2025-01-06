package commands

/*
import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// Handle the t!chart command
func HandleChart(s *discordgo.Session, m *discordgo.MessageCreate) {
	guildID := m.GuildID

	// Query the top 10 songs in the guild
	topSongs, err :=

	topSongs, err := db.TopSongsInGuild(database, guildID, 10) // Replace with your actual query logic
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
*/
