package main

import (
	"squad-go/server"

	"github.com/charmbracelet/log"
)

func main() {
	log.SetPrefix("[SquadGO]")
	log.Info("Starting SquadGO!")

	server := server.NewSquadServer()

	server.Boot()

	select {}
}
