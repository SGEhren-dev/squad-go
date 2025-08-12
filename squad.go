package main

import (
	"squad-go/configuration"
	"squad-go/server"

	"github.com/charmbracelet/log"
)

func main() {
	log.SetPrefix("[SquadGO]")
	log.Info("Starting SquadGO!")
	log.Info("Loading configuration file")

	config := &configuration.Config{}

	err := config.LoadConfig()

	if err != nil {
		return
	}

	log.Info("Configuration loaded successfully.")

	server := &server.SquadServer{
		Config: config,
	}

	server.Boot()

	select {}
}
