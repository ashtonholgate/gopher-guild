package app

import (
	"gopherService/config"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

type App struct {
	Router       *gin.Engine
	Dependencies AppDependencies
	Config       config.Config
}

func New(config config.Config) (*App, error) {
	dependencies, err := initialiseDependencies(config)
	if err != nil {
		return nil, err
	}

	router := gin.Default()

	router.Use(errorHandler())

	app := &App{
		Router:       router,
		Dependencies: dependencies,
		Config:       config,
	}

	app.setupRoutes(dependencies)

	app.Router.Use(errorHandler())
	return app, nil
}

func (a *App) setupRoutes(dependencies AppDependencies) {
	a.Router.GET("/health", healthHandler)

	a.Router.POST("/gophers", dependencies.GopherRouterService.CreateGopherEndpoint())
}

func (a *App) Run(addr string) error {
	return a.Router.Run(addr)
}

func healthHandler(c *gin.Context) {
	c.String(200, "healthy")
}
