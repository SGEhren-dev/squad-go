package server

import (
	"github.com/bwmarrin/discordgo"
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

// #region Discord Plugin

type DiscordPlugin struct {
	Plugin
	discordClient *discordgo.Session
}

type IDiscordPlugin interface {
	IPlugin
	SendMessage(string)
}

func (plugin *DiscordPlugin) Boot() {}

func (plugin *DiscordPlugin) SetupDiscordClient() {
	discord, err := discordgo.New("Bot " + plugin.SquadServer.Config.Discord.Token)

	if err != nil {
		plugin.LogWithPrefix("Failed to initialize Discord client")
	}

	plugin.discordClient = discord

	err = discord.Open()

	if err != nil {
		plugin.LogWithPrefix("Failed to open Discord connection")
	}
}

func (plugin *DiscordPlugin) Shutdown() {
	plugin.discordClient.Close()
}

// #endregion
