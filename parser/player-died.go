package parser

import (
	"regexp"
	"squad-go/events"
)

type PlayerDied struct {
	Attacker  string
	ChainID   string
	Damage    string
	Time      string
	Victim    string
	Weapon    string
	WoundTime string
}

var PlayerDiedParser = Parser{
	Regex: regexp.MustCompile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogSquadTrace: \[DedicatedServer](?:ASQSoldier::)?Die\(\): Player:(.+) KillingDamage=(?:-)*([0-9.]+) from ([A-z_0-9]+) \(Online IDs:([^)|]+)\| Contoller ID: ([\w\d]+)\) caused by ([A-z_0-9-]+)_C`),
	OnMatch: func(matches []string, parser *LogParser) {
		payload := &PlayerDied{
			Attacker:  matches[5],
			ChainID:   matches[2],
			Damage:    matches[4],
			Time:      matches[1],
			Victim:    matches[3],
			Weapon:    matches[8],
			WoundTime: matches[1],
		}

		parser.Emitter.Emit(events.PLAYER_DIED, payload)
	},
}
