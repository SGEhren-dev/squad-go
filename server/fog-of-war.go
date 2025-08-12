package server

import (
	"fmt"
	"time"
)

type FogOfWarPlugin struct {
	Plugin
	tickerHandle  *time.Ticker
	tickerChannel chan bool
}

func NewFogOfWarPlugin(server *SquadServer) *FogOfWarPlugin {
	return &FogOfWarPlugin{
		Plugin: Plugin{
			Enabled:     true,
			Name:        "FogOfWar",
			Description: "Automate setting the FogOfWar mode.",
			SquadServer: server,
		},
		tickerHandle:  nil,
		tickerChannel: nil,
	}
}

func (plugin *FogOfWarPlugin) Boot() {
	delay, ok := plugin.GetSettings()["delay"].(float64)

	if !ok {
		return
	}

	plugin.tickerHandle = time.NewTicker(time.Duration(delay) * time.Second)
	plugin.tickerChannel = make(chan bool)

	go func() {
		for {
			select {
			case <-plugin.tickerHandle.C:
				plugin.SquadServer.Rcon.Execute(fmt.Sprintf("AdminSetFogOfWar %v", plugin.GetSettings()["mode"].(float64)))
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
