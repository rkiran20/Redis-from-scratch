package main

import (
	"flag"
	"log"

	"github.com/dicedb/dice/config"
	"github.com/dicedb/dice/server"
)

func setupFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "Host for the Dice server")
	flag.IntVar(&config.Port, "port", 8123, "Port for the Dice server")
	flag.Parse()
}

func main() {
	setupFlags()
	log.Println("Starting Dice ")
	server.RunSyncTCPServer()
}
