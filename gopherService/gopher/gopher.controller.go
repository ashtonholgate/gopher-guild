package gopher

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	CreateGopherEndpoint() gin.HandlerFunc
}

type CommandServiceContract interface {
	Create(gopher IncomingGopher) (OutgoingGopher, error)
}

type controllerImpl struct {
	CommandService CommandServiceContract
}

func (cs *controllerImpl) CreateGopherEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newGopher IncomingGopher

		if err := c.ShouldBindJSON(&newGopher); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		createdGopher, err := cs.CommandService.Create(newGopher)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, createdGopher)
	}
}

func NewGopherController(commandService CommandService) Controller {
	return &controllerImpl{CommandService: commandService}
}
