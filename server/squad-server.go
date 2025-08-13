package server

import (
	"fmt"
	"squad-go/configuration"
	"squad-go/parser"

	rcon "github.com/SquadGO/squad-rcon-go/v2"
	"github.com/SquadGO/squad-rcon-go/v2/rconEvents"
	"github.com/charmbracelet/log"
	"github.com/iamalone98/eventEmitter"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SquadServer struct {
	Config   configuration.Config
	Database *gorm.DB
	Emitter  eventEmitter.EventEmitter
	Parser   *parser.LogParser
	Players  map[string][]any
	manager  *PluginManager
	Rcon     *rcon.Rcon
}

func NewSquadServer() *SquadServer {
	squadServer := SquadServer{
		Config:  configuration.Config{},
		Emitter: eventEmitter.NewEventEmitter(),
	}

	err := squadServer.Config.LoadConfig()

	if err != nil {
		log.Error("Failed to create SquadServer due to config failure.")

		return nil
	}

	log.Info("Configuration loaded, starting Squad Server")

	return &squadServer
}

func (server *SquadServer) Boot() {
	log.Info("Booting Squad Server...")

	server.Parser = parser.NewLogParser()
	server.manager = &PluginManager{}

	server.manager.Load(server)

	server.setupDatabase()
	server.setupRcon()
	server.manager.Boot()
	server.Parser.ParseLogFile(server.Config.LogFilePath)
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

func (server *SquadServer) setupDatabase() {
	log.Info("Setting up database connection")
	dialect := server.Config.Database.Dialect
	database := server.Config.Database.Name

	switch dialect {
	case "sqlite":
		{
			db, err := gorm.Open(sqlite.Open(database+".db"), &gorm.Config{})

			if err != nil {
				log.Error("Failed to initialize Database.")

				return
			}

			server.Database = db

			break
		}

	case "mysql":
		{
			username := server.Config.Database.Username
			password := server.Config.Database.Password
			hostname := server.Config.Database.Host
			port := server.Config.Database.Port
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, hostname, port, database)
			db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

			if err != nil {
				log.Error("Failed to initialize Database.")

				return
			}

			server.Database = db

			break
		}

	default:
		log.Warn("An unknown dialect was used: " + dialect)
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
		log.Error("Failed to setup RCON.")

		return
	}

	server.Rcon = rconHandle

	rconHandle.Emitter.On(rconEvents.CONNECTED, func(_ any) {
		log.WithPrefix("[RCON]").Info("Connected to RCON server at " + server.Config.Rcon.Host + ":" + server.Config.Rcon.Port)
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
}
