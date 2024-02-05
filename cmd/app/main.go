package main

import (
	"inditilla/config"
	"inditilla/internal/app"
	"log"
)

// Get config and run application with that config
func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
