package app

import (
	"gopherService/config"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

type App struct {
	Router       *mux.Router
	Dependencies AppDependencies
	Config       config.Config
}

func New(config config.Config) (*App, error) {
	dependencies, err := initialiseDependencies(config)
	if err != nil {
		return nil, err
	}

	app := App{
		Router:       mux.NewRouter(),
		Dependencies: dependencies,
		Config:       config,
	}

	setupRoutes(app.Router, dependencies)

	return &app, nil
}

func setupRoutes(r *mux.Router, dependencies AppDependencies) {
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/gophers", dependencies.GopherController.CreateGopherEndpoint).Methods("POST")
	r.HandleFunc("/gophers/{id}", dependencies.GopherController.ReadGopherEndpoint).Methods("GET")
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Healthy"}`))
}
