package server

import (
	"fmt"
	"time"
)

type AutomatedBroadcastPlugin struct {
	Plugin
	tickerHandle  *time.Ticker
	tickerChannel chan bool
}

var lastDispatched int = 0

func NewAutomatedBroadcastPlugin(server *SquadServer) *AutomatedBroadcastPlugin {
	return &AutomatedBroadcastPlugin{
		Plugin: Plugin{
			Enabled:     true,
			Name:        "AutomatedBroadcast",
			Description: "Automates broadcast messages.",
			SquadServer: server,
		},
		tickerHandle:  nil,
		tickerChannel: nil,
	}
}

func (plugin *AutomatedBroadcastPlugin) Boot() {
	delay, ok := plugin.GetSettings()["delay"].(float64)

	if !ok {
		return
	}

	plugin.tickerHandle = time.NewTicker(time.Duration(delay) * time.Second)
	plugin.tickerChannel = make(chan bool)
	messages := plugin.GetSettings()["messages"].([]any)

	go func() {
		for {
			select {
			case <-plugin.tickerHandle.C:
				plugin.SquadServer.Rcon.Execute(fmt.Sprintf("AdminBroadcast %s", messages[lastDispatched]))

				if lastDispatched+1 >= len(messages) {
					lastDispatched = 0
				} else {
					lastDispatched++
				}
			case <-plugin.tickerChannel:
				plugin.tickerHandle.Stop()
				return
			}
		}
	}()
}
