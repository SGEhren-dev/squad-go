package parser

import (
	"regexp"
	"squad-go/events"
)

type DeployableDamaged struct {
	Causer     string
	Damage     string
	DamageType string
	Deployable string
	Health     string
	Thread     string
	Time       string
	Weapon     string
}

var DeployableDamagedParser = Parser{
	Regex: regexp.MustCompile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogSquadTrace: \[DedicatedServer](?:ASQDeployable::)?TakeDamage\(\): ([A-Za-z0-9_]+)_C_[0-9]+: ([0-9.]+) damage attempt by causer ([A-Za-z0-9_]+)_C_[0-9]+ instigator (.+) with damage type ([A-Za-z0-9_]+)_C health remaining ([0-9.]+)`),
	OnMatch: func(matches []string, parser *LogParser) {
		payload := &DeployableDamaged{
			Causer:     matches[5],
			Damage:     matches[4],
			DamageType: matches[7],
			Deployable: matches[3],
			Health:     matches[8],
			Thread:     matches[2],
			Time:       matches[1],
			Weapon:     matches[7],
		}

		parser.Emit(events.DEPLOYABLE_DAMAGED, payload)
	},
}
