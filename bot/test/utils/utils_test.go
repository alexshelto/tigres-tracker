package utils

import (
	"github.com/alexshelto/tigres-tracker/utils"
	"testing"
)

func TestExtractUserID(t *testing.T) {
	parameters := []struct {
		input, expected string
	}{
		{"Requested By: <@1>", "1"},
		{"Requested By: <@1234567899>", "1234567899"},
		{"Requested By: <@4353474563399578719>", "4353474563399578719"},
		{"Requested By: <@4304215243331631818>", "4304215243331631818"},
		{"Requested By: <@999999999999999999999999999999>", "999999999999999999999999999999"},
	}

	for i := range parameters {
		actual := utils.ExtractUserID(parameters[i].input)
		if actual != parameters[i].expected {
			t.Logf("expected: '%s' , actual: '%s'", parameters[i].expected, actual)
			t.Fail()
		}
	}
}
