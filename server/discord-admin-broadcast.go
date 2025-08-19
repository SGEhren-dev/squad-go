package server

import (
	"encoding/json"
	"fmt"
	"squad-go/events"
	"squad-go/parser"

	"github.com/bwmarrin/discordgo"
)

type DiscordAdminBroadcastSettings struct {
	Channel string `json:"channel"`
}

type DiscordAdminBroadcastPlugin struct {
	Plugin
	settings DiscordAdminBroadcastSettings
}

func init() {
	var name string = "DiscordAdminBroadcast"

	RegisterPlugin(name, func(server *SquadServer, rawSettings json.RawMessage) IPlugin {
		var settings DiscordAdminBroadcastSettings

		if err := json.Unmarshal(rawSettings, &settings); err != nil {
			return nil
		}

		return &DiscordAdminBroadcastPlugin{
			Plugin: Plugin{
				Name:        name,
				Description: "Send all admin broadcasts to Discord.",
				SquadServer: server,
			},
			settings: settings,
		}
	})
}

func (plugin *DiscordAdminBroadcastPlugin) Boot() {
	plugin.SquadServer.Parser.On(events.ADMIN_BROADCAST, func(payload any) {
		data := payload.(parser.AdminBroadcast)

		if plugin.SquadServer.Discord == nil && len(plugin.settings.Channel) > 0 {
			return
		}

		plugin.SquadServer.Discord.ChannelMessageSendEmbed(plugin.settings.Channel, &discordgo.MessageEmbed{
			Title:       "Admin Broadcast Sent",
			Description: fmt.Sprintf("Admin '%s' sent a broadcast", data.From),
			Timestamp:   data.Time,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Message",
					Value: data.Message,
				},
			},
		})
	})
}
