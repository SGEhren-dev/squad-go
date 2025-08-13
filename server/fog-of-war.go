package server

import (
	"encoding/json"
	"fmt"
	"time"
)

type FogOfWarSettings struct {
	Delay float64 `json:"delay"`
	Mode  float64 `json:"mode"`
}

type FogOfWarPlugin struct {
	Plugin
	settings      FogOfWarSettings
	tickerHandle  *time.Ticker
	tickerChannel chan bool
}

func init() {
	var name string = "FogOfWar"

	RegisterPlugin(name, func(server *SquadServer, rawSettings json.RawMessage) IPlugin {
		var settings FogOfWarSettings

		if err := json.Unmarshal(rawSettings, &settings); err != nil {
			return nil
		}

		return &FogOfWarPlugin{
			Plugin: Plugin{
				Name:        name,
				Description: "Automatically sets the FogOfWar to the given mode after a delay.",
				SquadServer: server,
			},
			settings:      settings,
			tickerHandle:  nil,
			tickerChannel: make(chan bool),
		}
	})
}

func (plugin *FogOfWarPlugin) Boot() {
	plugin.tickerHandle = time.NewTicker(time.Duration(plugin.settings.Delay) * time.Second)
	plugin.tickerChannel = make(chan bool)

	go func() {
		for {
			select {
			case <-plugin.tickerHandle.C:
				if plugin.SquadServer.Rcon != nil {
					plugin.SquadServer.Rcon.Execute(fmt.Sprintf("AdminSetFogOfWar %v", plugin.settings.Mode))
				}
			case <-plugin.tickerChannel:
				plugin.tickerHandle.Stop()
				return
			}
		}
	}()
}

func (plugin *FogOfWarPlugin) Shutdown() {
	if plugin.tickerHandle != nil && plugin.tickerChannel != nil {
		plugin.tickerChannel <- true
	}
}
