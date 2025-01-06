package handlers

/*
import (
	"log"
	"strings"

	"github.com/alexshelto/tigres-tracker/commands"
	"github.com/bwmarrin/discordgo"
)

func HandleCommands(s *discordgo.Session, m *discordgo.MessageCreate) {

	log.Printf("Inside of handler commands")

	content := strings.ToLower(m.Content)

	switch {
	case strings.HasPrefix(content, "t!help"):
		commands.HandleHelp(s, m)
	case strings.HasPrefix(content, "t!chart"):
		commands.HandleChart(s, m, db.GetDB())
	case strings.HasPrefix(content, "t!stats"):
		HandleStatsCommand(s, m)
	}
}

func HandleStatsCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	line := strings.TrimSpace(m.Content[len("t!stats"):])
	var userId string
	if line == "" {
		userId = m.Author.ID
	} else {
		userId = strings.Trim(line, "<@>")
	}
	commands.HandleStats(s, m, userId, db.GetDB())
}

*/
