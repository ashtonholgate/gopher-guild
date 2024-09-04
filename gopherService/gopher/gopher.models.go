package gopher

import (
	"gopherService/config"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type BaseGopher struct {
	Name  string `json:"name" binding:"required,validateNameLength"`
	Age   *int   `json:"age" binding:"required,min=0"`
	Color string `json:"color" binding:"required,validateColorLength"`
}

type IncomingGopher struct {
	BaseGopher
}

type OutgoingGopher struct {
	BaseGopher
	Id int `json:"id"`
}

func validateNameLength(fl validator.FieldLevel, max int) bool {
	name := fl.Field().String()
	return len(name) <= max
}

func validateColorLength(fl validator.FieldLevel, max int) bool {
	color := fl.Field().String()
	return len(color) <= max
}

// InitialiseValidator initialises the custom validators necessary for parsing gophers
func InitialiseValidator(configuration config.Config) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validateNameLength", func(fl validator.FieldLevel) bool {
			return validateNameLength(fl, configuration.MAXIMUM_GOPHER_NAME_LENGTH)
		})
		v.RegisterValidation("validateColorLength", func(fl validator.FieldLevel) bool {
			return validateColorLength(fl, configuration.MAXIMUM_GOPHER_COLOR_LENGTH)
		})
	}
}
