package main

import (
	"ecommerce-app/config"
	"ecommerce-app/internal/api"
	"log"
)

func main() {

	cfg, err := config.SetupEnvironment()

	if err != nil {
		log.Fatalf("Error setting up environment: %v", err)
		return
	}
	api.StartServer(cfg)
}
