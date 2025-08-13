package server

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/log"
)

type PluginFactory func(server *SquadServer, rawSettings json.RawMessage) IPlugin

var pluginRegistry = map[string]PluginFactory{}

func RegisterPlugin(name string, factory PluginFactory) {
	pluginRegistry[name] = factory
}

type PluginManager struct {
	Plugins []IPlugin
}

func (manager *PluginManager) Load(server *SquadServer) {
	for name, config := range server.Config.Plugins {
		factory, exists := pluginRegistry[name]

		if !exists {
			log.Warn(fmt.Sprintf("No factory registered for plugin [%s]", name))
			continue
		}

		plugin := factory(server, config.Settings)
		plugin.GetBase().Enabled = config.Enabled
		manager.Plugins = append(manager.Plugins, plugin)

		log.Info(fmt.Sprintf("Factory built for plugin [%s]", name))
	}
}

func (manager *PluginManager) Boot() {
	for _, plugin := range manager.Plugins {
		if plugin.GetBase().Enabled {
			plugin.Boot()

			log.Info(fmt.Sprintf("Plugin [%s] was booted", plugin.GetBase().Name))
		}
	}
}
