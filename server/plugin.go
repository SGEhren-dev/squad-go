package server

import "github.com/charmbracelet/log"

type Plugin struct {
	Enabled     bool
	Name        string
	Description string
	SquadServer *SquadServer
}

type IPlugin interface {
	Boot()
	GetBase() *Plugin
	Shutdown()
}

func (plugin *Plugin) Boot() {}

func (plugin *Plugin) LogWithPrefix(message any) {
	log.WithPrefix("[" + plugin.Name + "]").Info(message)
}

func (plugin *Plugin) Shutdown() {}

// #region Getters and Setters

func (plugin *Plugin) GetBase() *Plugin {
	return plugin
}

func (plugin *Plugin) GetSettings() map[string]any {
	for _, confPlugin := range plugin.SquadServer.Config.Plugins {
		if plugin.Name == confPlugin.Name {
			return confPlugin.Settings
		}
	}

	return nil
}

func (plugin *Plugin) SetSquadServer(server *SquadServer) {
	if server == plugin.SquadServer {
		return
	}

	plugin.SquadServer = server
}

// #endregion
