package utils

import (
	"testing"

	"github.com/alexshelto/tigres-tracker/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestProcessEmbedDataForNowPlaying(t *testing.T) {
	title := "Now Playing"
	description1 := "Rooster (2022 Remaster) - Alice In Chains\n`[0:00 / 06:14]`\n\nRequested by: <@123>"
	description2 := "The Cup - Dave Blunts\n`[0:00 / 02:26]`\n\nRequested by: <@456>"

	expectedSongName1 := "Rooster (2022 Remaster) - Alice In Chains"
	expectedSongName2 := "The Cup - Dave Blunts"

	parameters := []struct {
		input    *discordgo.MessageEmbed
		expected *utils.ParsedSongInfo
	}{
		{
			&discordgo.MessageEmbed{Title: title, Description: description1},
			&utils.ParsedSongInfo{Name: expectedSongName1, RequestedBy: "123"},
		},
		{
			&discordgo.MessageEmbed{Title: title, Description: description2},
			&utils.ParsedSongInfo{Name: expectedSongName2, RequestedBy: "456"},
		},
	}

	for i := range parameters {
		actual := utils.ProcessEmbedDataForNowPlaying(parameters[i].input)
		assert.Equal(t, parameters[i].expected, actual)
	}
}
