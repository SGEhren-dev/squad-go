package parser

import (
	"regexp"
	"squad-go/events"
)

type PlayerConnected struct {
	Time             string
	PlayerController string
	IP               string
	SteamID          string
	EOSID            string
}

var PlayerConnectedParser = Parser{
	Regex: regexp.MustCompile(
		`^\[([0-9]{4}\.[0-9]{2}\.[0-9]{2}-[0-9]{2}\.[0-9]{2}\.[0-9]{2}:[0-9]{3})\].*?` +
			`(BP_PlayerController_C_[0-9]+).*?` +
			`IP:\s*([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+).*?` +
			`EOS:\s*([0-9a-f]+).*?` +
			`steam:\s*([0-9]+)`,
	),
	OnMatch: func(matches []string, parser *LogParser) {
		payload := &PlayerConnected{
			EOSID:            matches[4],
			IP:               matches[3],
			PlayerController: matches[2],
			SteamID:          matches[5],
			Time:             matches[1],
		}

		parser.Emit(events.PLAYER_CONNECTED, payload)
	},
}
