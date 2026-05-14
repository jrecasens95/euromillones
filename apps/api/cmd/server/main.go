package main

import (
	"log"

	"euromillones/internal/platform/config"
	"euromillones/internal/server"
)

func main() {
	config.Load()

	app, err := server.New()
	if err != nil {
		log.Fatalf("failed to bootstrap server: %v", err)
	}

	log.Fatal(app.Listen(":" + config.Current.Port))
}
