package parser

import (
	"regexp"
	"squad-go/events"
)

type PlayerRevived struct {
	ChainID string
	Reviver string
	SteamID string
	Time    string
	Victim  string
}

var PlayerRevivedParser = Parser{
	Regex: regexp.MustCompile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogSquad: (.+) \(Online IDs:([^)]+)\) has revived (.+) \(Online IDs:([^)]+)\)\.`),
	OnMatch: func(matches []string, parser *LogParser) {
		payload := &PlayerRevived{
			ChainID: matches[2],
			Reviver: matches[3],
			SteamID: matches[4],
			Time:    matches[1],
			Victim:  matches[5],
		}

		parser.Emit(events.PLAYER_REVIVED, payload)
	},
}
