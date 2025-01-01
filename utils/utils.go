package utils

import (
	"fmt"
	"regexp"
)

const PancakeBotID = "239631525350604801"

func IsFromPancakeBot(authorId string) bool {
	return authorId == PancakeBotID
}

// extract user ID from a string
// ex string: `Requested by: <@2142792696840213736>`
func ExtractUserID(requestString string) string {
	fmt.Printf("Received string: '%s'\n", requestString)
	re := regexp.MustCompile(`<@(\d+)>`)
	matches := re.FindStringSubmatch(requestString)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
