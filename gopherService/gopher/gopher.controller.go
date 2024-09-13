package gopher

import (
	"errors"
	"fmt"
	"gopherService/customErrors"
	"log"
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

		fmt.Println(newGopher.Validate())

		if err := newGopher.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		createdGopher, err := cs.CommandService.Create(newGopher)
		if err != nil {
			log.Printf("error: %+v", err)

			var databaseErr *customErrors.DatabaseError

			switch {
			case errors.As(err, &databaseErr):
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Database operation failed",
					"details": "An error occurred while processing your request. Please try again later.",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "An unexpected error occurred",
					"details": "Please try again later or contact support if the problem persists.",
				})
			}
			return
		}

		c.JSON(http.StatusCreated, createdGopher)
	}
}

func NewGopherController(commandService CommandService) Controller {
	return &controllerImpl{CommandService: commandService}
}
