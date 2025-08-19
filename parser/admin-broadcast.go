package parser

import (
	"regexp"
	"squad-go/events"
)

type AdminBroadcast struct {
	From    string
	Message string
	Time    string
}

var AdminBroadcastParser = Parser{
	Regex: regexp.MustCompile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogSquad: ADMIN COMMAND: Message broadcasted <(.+)> from (.+)`),
	OnMatch: func(matches []string, parser *LogParser) {
		payload := &AdminBroadcast{
			From:    matches[4],
			Message: matches[3],
			Time:    matches[1],
		}

		parser.Emit(events.ADMIN_BROADCAST, payload)
	},
}
