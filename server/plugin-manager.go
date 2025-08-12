package server

import (
	"fmt"

	"github.com/charmbracelet/log"
)

type PluginManager struct {
	Plugins []IPlugin
	server  *SquadServer
}

func NewPluginManager(server *SquadServer) PluginManager {
	return PluginManager{
		Plugins: []IPlugin{
			NewAutomatedBroadcastPlugin(server),
			NewFogOfWarPlugin(server),
		},
		server: server,
	}
}

func (m *PluginManager) BootAll() {
	for _, plugin := range m.Plugins {
		base := plugin.GetBase()

		if !base.Enabled {
			log.Info(fmt.Sprintf("Plugin [%s] not enabled", base.Name))
			continue
		}

		plugin.Boot()

		log.Info(fmt.Sprintf("Plugin [%s] booted", base.Name))
	}
}

func (manager *PluginManager) Shutdown() {
	for _, plugin := range manager.Plugins {
		plugin.Shutdown()
	}
}
