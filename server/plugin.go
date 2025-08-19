package server

import (
	"github.com/charmbracelet/log"
)

// #region Plugin

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

func (plugin *Plugin) GetBase() *Plugin {
	return plugin
}

func (plugin *Plugin) LogWithPrefix(message any) {
	log.WithPrefix("[" + plugin.Name + "]").Info(message)
}

func (plugin *Plugin) Shutdown() {}

// #endregion
