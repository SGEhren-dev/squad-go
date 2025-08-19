package server

import (
	"encoding/json"
	"fmt"
	"time"
)

type AutomatedBroadcastSettings struct {
	Delay    float64  `json:"delay"`
	Messages []string `json:"messages"`
}

type AutomatedBroadcastPlugin struct {
	Plugin
	settings      AutomatedBroadcastSettings
	tickerHandle  *time.Ticker
	tickerChannel chan bool
}

var lastDispatched int = 0

func init() {
	var name string = "AutomatedBroadcast"

	RegisterPlugin(name, func(server *SquadServer, rawSettings json.RawMessage) IPlugin {
		var settings AutomatedBroadcastSettings
		if err := json.Unmarshal(rawSettings, &settings); err != nil {
			return nil
		}

		return &AutomatedBroadcastPlugin{
			Plugin: Plugin{
				Name:        name,
				Description: "Automatically sends a broadcast at a set interval.",
				SquadServer: server,
			},
			settings:      settings,
			tickerHandle:  nil,
			tickerChannel: make(chan bool),
		}
	})
}

func (plugin *AutomatedBroadcastPlugin) Boot() {
	plugin.tickerHandle = time.NewTicker(time.Duration(plugin.settings.Delay) * time.Second)
	messages := plugin.settings.Messages

	go func() {
		for {
			select {
			case <-plugin.tickerHandle.C:
				if plugin.SquadServer.Rcon != nil {
					plugin.SquadServer.Rcon.Execute(fmt.Sprintf("AdminBroadcast %s", messages[lastDispatched]))
				}

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
