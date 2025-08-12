package parser

import (
	"regexp"
	"squad-go/events"
)

type PlayerDisconnected struct {
	ChainID          string
	EOSID            string
	IP               string
	PlayerController string
	Time             string
}

var PlayerDisconnectedParser = Parser{
	Regex: regexp.MustCompile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogNet: UChannel::Close: Sending CloseBunch\. ChIndex == [0-9]+\. Name: \[UChannel\] ChIndex: [0-9]+, Closing: [0-9]+ \[UNetConnection\] RemoteAddr: ([\d.]+):[\d]+, Name: EOSIpNetConnection_[0-9]+, Driver: GameNetDriver EOSNetDriver_[0-9]+, IsServer: YES, PC: ([^ ]+PlayerController(?:|.+)_C_[0-9]+), Owner: [^ ]+PlayerController(?:|.+)_C_[0-9]+, UniqueId: RedpointEOS:([\d\w]+)/`),
	OnMatch: func(matches []string, parser *LogParser) {
		payload := &PlayerDisconnected{
			ChainID:          matches[2],
			EOSID:            matches[5],
			IP:               matches[3],
			PlayerController: matches[4],
			Time:             matches[1],
		}

		parser.Emitter.Emit(events.PLAYER_DISCONNECTED, payload)
	},
}
