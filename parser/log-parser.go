package parser

import (
	"io"
	"regexp"

	"github.com/charmbracelet/log"
	"github.com/iamalone98/eventEmitter"
	"github.com/papertrail/go-tail/follower"
)

type Parser struct {
	Regex   *regexp.Regexp
	OnMatch func([]string, *LogParser)
}

type LogParser struct {
	eventEmitter.EventEmitter
	TailHandle *follower.Follower
	Rules      []*Parser
}

func NewLogParser() *LogParser {
	newParser := LogParser{
		EventEmitter: eventEmitter.NewEventEmitter(),
		Rules: []*Parser{
			&AdminBroadcastParser,
			&DeployableDamagedParser,
			&NewGameParser,
			&PlayerConnectedParser,
			&PlayerDamagedParser,
			&PlayerDiedParser,
			&PlayerDisconnectedParser,
			&PlayerRevivedParser,
			&RoundWinnerParser,
			&ServerTickRateParser,
		},
	}

	return &newParser
}

func (parser *LogParser) ParseLogFile(filepath string) {
	log.Info("Starting log file parser...")

	tail, err := follower.New(filepath, follower.Config{
		Whence: io.SeekEnd,
		Offset: 0,
		Reopen: true,
	})

	parser.TailHandle = tail

	if err != nil {
		log.Error("Error opening log file: " + err.Error())

		return
	}

	for line := range tail.Lines() {
		for _, rule := range parser.Rules {
			if rule.Regex.MatchString(line.String()) {
				matches := rule.Regex.FindStringSubmatch(line.String())

				rule.OnMatch(matches, parser)
			}
		}
	}

	if tail.Err() != nil {
		log.WithPrefix("[LogParser]").Error(tail.Err())
	}
}
