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
	DiscordPlugin
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
			DiscordPlugin: DiscordPlugin{
				Plugin: Plugin{
					Name:        name,
					Description: "Send all admin broadcasts to Discord.",
					SquadServer: server,
				},
				discordClient: nil,
			},
			settings: settings,
		}
	})
}

func (plugin *DiscordAdminBroadcastPlugin) Boot() {
	plugin.SetupDiscordClient()

	plugin.SquadServer.Parser.Emitter.On(events.ADMIN_BROADCAST, func(payload any) {
		data := payload.(parser.AdminBroadcast)

		if plugin.discordClient == nil {
			return
		}

		plugin.discordClient.ChannelMessageSendEmbed(plugin.settings.Channel, &discordgo.MessageEmbed{
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
