package main

import (
	"gopherService/app"
	"gopherService/config"
	"log"
	"net/http"
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

	http.Handle("/", app.Router)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
