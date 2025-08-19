package parser

import (
	"regexp"
	"squad-go/events"
)

type NewGame struct {
	ChainID   string
	DLC       string
	LayerName string
	MapName   string
	Time      string
}

var NewGameParser = Parser{
	Regex: regexp.MustCompile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogWorld: Bringing World \/([A-z0-9]+)\/(?:Maps\/)?([A-z0-9-]+)\/(?:.+\/)?([A-z0-9-]+)(?:\.[A-z0-9-]+)`),
	OnMatch: func(matches []string, parser *LogParser) {
		payload := &NewGame{
			ChainID:   matches[2],
			DLC:       matches[3],
			LayerName: matches[5],
			MapName:   matches[4],
			Time:      matches[1],
		}

		parser.Emit(events.NEW_GAME, payload)
	},
}
