package parser

import (
	"regexp"
	"squad-go/events"
	"strconv"

	"github.com/charmbracelet/log"
)

type ServerTickRate struct {
	ChainID  string
	TickRate float64
	Time     string
}

var ServerTickRateParser = Parser{
	Regex: regexp.MustCompile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogSquad: USQGameState: Server Tick Rate: ([0-9.]+)`),
	OnMatch: func(matches []string, parser *LogParser) {
		tickRate, err := strconv.ParseFloat(matches[3], 64)

		if err != nil {
			log.Error("Failed to parse tick rate: " + err.Error())

			return
		}

		payload := &ServerTickRate{
			ChainID:  matches[2],
			TickRate: tickRate,
			Time:     matches[1],
		}

		parser.Emitter.Emit(events.SERVER_TICK_RATE, payload)
	},
}
