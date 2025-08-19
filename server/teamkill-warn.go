package server

import (
	"squad-go/events"
	"squad-go/parser"
)

type TeamKillWarnPlugin struct {
	Plugin
}

func NewTeamKillWarnPlugin(server *SquadServer) *TeamKillWarnPlugin {
	return &TeamKillWarnPlugin{
		Plugin{
			Enabled:     false,
			Name:        "TeamKillWarn",
			Description: "Warns players when they team kill.",
			SquadServer: server,
		},
	}
}

func (plugin *TeamKillWarnPlugin) Boot() {
	plugin.SquadServer.Parser.On(events.PLAYER_CONNECTED, plugin.HandlePlayerConnected)
}

func (plugin *TeamKillWarnPlugin) HandlePlayerConnected(payload any) {
	data := payload.(*parser.PlayerConnected)

	plugin.LogWithPrefix(data.Time + " " + data.IP + " " + data.PlayerController)
}

func (plugin *TeamKillWarnPlugin) Shutdown() {
	plugin.SquadServer.Parser.RemoveListener(events.PLAYER_CONNECTED, plugin.HandlePlayerConnected)
}
