package parser

import (
	"regexp"
	"squad-go/events"
)

type RoundWinner struct {
	ChainID string
	Layer   string
	Time    string
	Winner  string
}

var RoundWinnerParser = Parser{
	Regex: regexp.MustCompile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogSquadTrace: \[DedicatedServer](?:ASQGameMode::)?DetermineMatchWinner\(\): (.+) won on (.+)`),
	OnMatch: func(matches []string, parser *LogParser) {
		payload := &RoundWinner{
			ChainID: matches[2],
			Layer:   matches[4],
			Time:    matches[1],
			Winner:  matches[3],
		}

		parser.Emitter.Emit(events.ROUND_WINNER, payload)
	},
}
