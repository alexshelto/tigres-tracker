package utils

import (
	"fmt"
	"regexp"
)

// Function to extract user ID from a string
// ex string: `Requested by: @joedale`
func ExtractUserID(requestString string) string {
	fmt.Printf("Received string: '%s'\n", requestString)
	re := regexp.MustCompile(`<@(\d+)>`)
	matches := re.FindStringSubmatch(requestString)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
