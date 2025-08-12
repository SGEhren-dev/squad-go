package parser

import (
	"regexp"
	"squad-go/events"
	"strconv"

	"github.com/charmbracelet/log"
)

type PlayerDamaged struct {
	AttackerName             string
	AttackerPlayerController string
	ChainID                  string
	Damage                   float64
	Time                     string
	Victim                   string
	Weapon                   string
}

var PlayerDamagedParser = Parser{
	Regex: regexp.MustCompile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogSquad: Player:(.+) ActualDamage=([0-9.]+) from (.+) \(Online IDs:([^|]+)\| Player Controller ID: ([^ ]+)\)caused by ([A-z_0-9-]+)_C`),
	OnMatch: func(matches []string, parser *LogParser) {
		parsedDamage, err := strconv.ParseFloat(matches[4], 64)

		if err != nil {
			log.Warn("Failed to parse damage value")

			return
		}

		payload := &PlayerDamaged{
			AttackerName:             matches[5],
			AttackerPlayerController: matches[7],
			ChainID:                  matches[2],
			Damage:                   parsedDamage,
			Time:                     matches[1],
			Victim:                   matches[3],
			Weapon:                   matches[8],
		}

		parser.Emitter.Emit(events.PLAYER_DAMAGED, payload)
	},
}
