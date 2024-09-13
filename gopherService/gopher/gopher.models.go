package gopher

import (
	"gopherService/customErrors"
	"strings"
)

type BaseGopher struct {
	Name  string `json:"name"`
	Age   *int   `json:"age"`
	Color string `json:"color"`
}

type IncomingGopher struct {
	BaseGopher
}

type OutgoingGopher struct {
	BaseGopher
	Id int `json:"id"`
}

func (g *IncomingGopher) Validate() error {
	issues := []string{}

	if g.Name == "" {
		issues = append(issues, "name is required")
	}
	if g.Age == nil {
		issues = append(issues, "age is required")
	} else if g.Age != nil {
		if *g.Age < 0 {
			issues = append(issues, "age must be non-negative")
		}
	}
	if g.Color == "" {
		issues = append(issues, "color is required")
	}
	if strings.ToLower(g.Color) == "red" {
		issues = append(issues, "color can not be red")
	}
	if len(issues) > 0 {
		return &customErrors.ValidationError{Issues: strings.Join(issues, ", ")}
	}
	return nil
}
