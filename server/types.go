package server

import "github.com/SquadGO/squad-rcon-go/v2/rconTypes"

type Teamkill struct {
	AttackerName string
	TeamID       string
	VictimName   string
	Weapon       string
}

type PlayerDamaged struct {
	Attacker rconTypes.Player
	Victim   rconTypes.Player
}
