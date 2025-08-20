package server

import (
	"encoding/json"
	"fmt"
	"squad-go/configuration"
	"squad-go/events"
	"squad-go/layers"
	"squad-go/parser"

	rcon "github.com/SquadGO/squad-rcon-go/v2"
	"github.com/SquadGO/squad-rcon-go/v2/rconEvents"
	"github.com/SquadGO/squad-rcon-go/v2/rconTypes"
	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/iamalone98/eventEmitter"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SquadServer struct {
	eventEmitter.EventEmitter
	Config   configuration.Config
	Database *gorm.DB
	Discord  *discordgo.Session
	Layers   *layers.Layers
	Parser   *parser.LogParser
	Players  rconTypes.Players
	Rcon     *rcon.Rcon
	manager  *PluginManager
}

func NewSquadServer() SquadServer {
	squadServer := SquadServer{
		Config:       configuration.Config{},
		EventEmitter: eventEmitter.NewEventEmitter(),
		Layers:       layers.New(),
	}

	err := squadServer.Config.LoadConfig()

	if err != nil {
		panic(err)
	}

	log.Info("Configuration loaded, starting Squad Server")

	return squadServer
}

func (server *SquadServer) Boot() {
	log.Info("Booting Squad Server...")

	server.Parser = parser.NewLogParser()
	server.manager = &PluginManager{}

	server.manager.Load(server)

	server.initializeConnectors()
	server.setupRcon()
	server.setupLogParser()
	server.manager.Boot()
}

func (server *SquadServer) Shutdown() {
	if server.Rcon != nil {
		server.Rcon.Close()

		log.Info("RCON connection closed.")
	}

	if server.Parser != nil && server.Parser.TailHandle != nil {
		server.Parser.TailHandle.Close()
		log.Info("Log file parser stopped.")
	}

	log.Info("Squad server shutdown successfully.")
}

func (server *SquadServer) initializeConnectors() {
	for connector, config := range server.Config.Connectors {
		switch connector {
		case "discord":
			{
				var token string

				if err := json.Unmarshal(config, &token); err != nil {
					continue
				}

				discord, err := discordgo.New("Bot " + token)

				if err != nil {
					log.Errorf("An error occurred when initializing the Discord connector: %s", err.Error())

					continue
				}

				server.Discord = discord

				err = server.Discord.Open()

				if err != nil {
					log.Errorf("An error occurred when initializing the Discord connector: %s", err.Error())

					continue
				}

				break
			}

		case "mysql":
			{
				var databaseConfig DatabaseConnector

				if err := json.Unmarshal(config, &databaseConfig); err != nil {
					continue
				}

				username := databaseConfig.Username
				password := databaseConfig.Password
				hostname := databaseConfig.Hostname
				database := databaseConfig.DatabaseName
				port := databaseConfig.Port
				dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, hostname, port, database)
				db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

				if err != nil {
					log.Errorf("Failed to initialize Database: %s", err.Error())

					return
				}

				server.Database = db

				break
			}

		default:
			log.Infof("Unknown connector %s", connector)
		}
	}
}

func (server *SquadServer) setupRcon() {
	log.WithPrefix("[RCON]").Info("Setting up RCON connection...")

	rconHandle, err := rcon.NewRcon(rcon.RconConfig{
		Host:               server.Config.Rcon.Host,
		Port:               server.Config.Rcon.Port,
		Password:           server.Config.Rcon.Password,
		AutoReconnect:      true,
		AutoReconnectDelay: 5,
	})

	if err != nil {
		log.WithPrefix("[RCON]").Error("Failed to setup RCON.")

		return
	}

	server.Rcon = rconHandle

	rconHandle.Emitter.On(rconEvents.CONNECTED, func(_ any) {
		log.WithPrefix("[RCON]").Infof(
			"Connected to RCON server at %s:%s",
			server.Config.Rcon.Host,
			server.Config.Rcon.Port,
		)
	})

	rconHandle.Emitter.On(rconEvents.RECONNECTING, func(_ any) {
		log.WithPrefix("[RCON]").Info("Attempting to reconnect to RCON.")
	})

	rconHandle.Emitter.On(rconEvents.CLOSE, func(_ any) {
		log.WithPrefix("[RCON]").Info("RCON connection closed.")
	})

	rconHandle.Emitter.On(rconEvents.ERROR, func(err any) {
		log.WithPrefix("[RCON]").Info("Error: " + err.(error).Error())
	})

	rconHandle.Emitter.On(rconEvents.LIST_PLAYERS, func(payload any) {
		if players, ok := payload.(rconTypes.Players); ok {
			server.Players = players
		}
	})

	log.WithPrefix("[RCON]").Infof(
		"Connection established to %s:%s",
		server.Config.Rcon.Host,
		server.Config.Rcon.Port,
	)
}

func (server *SquadServer) setupLogParser() {
	server.Parser.ParseLogFile(server.Config.LogFilePath)

	server.Parser.On(events.PLAYER_CONNECTED, func(payload any) {
		data := payload.(parser.PlayerConnected)

		log.Infof(
			"Player connected at - %s - Steam ID: %s - EOS ID: %s - IP: %s",
			data.Time,
			data.SteamID,
			data.EOSID,
			data.IP,
		)

		server.Emit(events.PLAYER_CONNECTED, data)
	})

	server.Parser.On(events.PLAYER_DIED, func(payload any) {
		if data, ok := payload.(parser.PlayerDied); ok {
			attacker := server.GetPlayerByEosId(data.Attacker)
			victim := server.GetPlayerByEosId(data.Victim)

			if victim.TeamID == attacker.TeamID {
				server.Emit(events.TEAMKILL, Teamkill{
					AttackerName: attacker.PlayerName,
					TeamID:       victim.TeamID,
					VictimName:   victim.PlayerName,
					Weapon:       data.Weapon,
				})
			}
		}
	})

	server.Parser.On(events.ADMIN_BROADCAST, func(payload any) {
		server.Emit(events.ADMIN_BROADCAST, payload)
	})

	server.Parser.On(events.PLAYER_DAMAGED, func(payload any) {
		if data, ok := payload.(parser.PlayerDamaged); ok {
			server.Emit(events.PLAYER_DAMAGED, PlayerDamaged{
				Attacker: server.GetPlayerByEosId(data.Victim),
				Victim:   server.GetPlayerByEosId(data.AttackerName),
			})
		}
	})
}

func (server *SquadServer) GetPlayerWithPredicate(predicate func(rconTypes.Player) bool) rconTypes.Players {
	result := make(rconTypes.Players, 0, len(server.Players))

	for _, player := range server.Players {
		if predicate(player) {
			result = append(result, player)
		}
	}

	return result
}

func (server *SquadServer) GetPlayerByEosId(identifier string) rconTypes.Player {
	foundPlayers := server.GetPlayerWithPredicate(func(player rconTypes.Player) bool {
		return player.EosID == identifier
	})

	if len(foundPlayers) == 0 {
		return rconTypes.Player{}
	}

	return foundPlayers[0]
}
