package gopher

import (
	"gopherService/customErrors"
	"strings"
)

type Address struct {
	Street   string
	Postcode string
}

type BaseGopher struct {
	Name    string   `json:"name"`
	Age     *int     `json:"age"`
	Color   string   `json:"color"`
	Address *Address `json:"address"`
}

type IncomingGopher struct {
	BaseGopher
}

type OutgoingGopher struct {
	BaseGopher
	Id int `json:"id"`
}

func (g *IncomingGopher) Validate() error {
	var issues []string

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

	if g.Address == nil {
		issues = append(issues, "address is required")
	} else {
		if g.Address.Street == "" {
			issues = append(issues, "address.street is required")
		}
		if g.Address.Postcode == "" {
			issues = append(issues, "address.postcode is required")
		}
	}

	if len(issues) > 0 {
		return &customErrors.ValidationError{Issues: strings.Join(issues, ", ")}
	}
	return nil
}
