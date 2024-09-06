package app

import (
	"database/sql"
	"fmt"
	"gopherService/config"
	"gopherService/gopher"
)

type AppDependencies struct {
	DB                      *sql.DB
	GopherRepositoryService gopher.RepositoryService
	GopherCommandService    gopher.CommandService
	GopherRouterService     gopher.Controller
}

type DatabaseInitializationError struct{}

func (e *DatabaseInitializationError) Error() string {
	return "Failed to connect to the database"
}

func initialiseDB(configuration config.Config) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		configuration.DB_HOST, configuration.DB_PORT, configuration.DB_USER, configuration.DB_PASSWORD, configuration.DB_NAME)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db, db.Ping()
}

func initialiseDependencies(configuration config.Config) (AppDependencies, error) {
	db, err := initialiseDB(configuration)
	gopher.InitialiseValidator(configuration)
	if err != nil {
		return AppDependencies{}, &DatabaseInitializationError{}
	}
	gopherRepositoryService := gopher.NewGopherRepositoryService(db)
	gopherCommandService := gopher.NewGopherCommandService(gopherRepositoryService)
	gopherRouterService := gopher.NewGopherController(gopherCommandService)
	return AppDependencies{
		DB:                      db,
		GopherRepositoryService: gopherRepositoryService,
		GopherCommandService:    gopherCommandService,
		GopherRouterService:     gopherRouterService,
	}, nil
}
