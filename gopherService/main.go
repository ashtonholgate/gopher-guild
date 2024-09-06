package main

import (
	"gopherService/app"
	"gopherService/config"
	"log"
	"os"
)

func main() {
	configuration, err := config.New()
	if err != nil {
		log.Printf("Failed to load configuration: %v", err)
		os.Exit(1)
	}

	app, err := app.New(configuration)
	if err != nil {
		log.Printf("Failed to initialize application: %v", err)
		os.Exit(1)
	}
	defer app.Dependencies.DB.Close()

	if err := app.Run(":8080"); err != nil {
		log.Printf("Failed to run server: %v", err)
		os.Exit(1)
	}
}
