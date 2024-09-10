package app

import (
	"errors"
	"gopherService/customErrors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			log.Printf("error: %+v", err)

			var gopherColorInvalidErr *customErrors.GopherColorInvalidError
			var databaseErr *customErrors.DatabaseError

			switch {
			case errors.As(err, &gopherColorInvalidErr):
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid gopher color",
					"details": gopherColorInvalidErr.Error(),
				})
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
		}
	}
}
