package customErrors

import "fmt"

type EnvVarNotFoundError struct {
	EnvVar string
}

func (e *EnvVarNotFoundError) Error() string {
	return fmt.Sprintf("Environment variable %v not found", e.EnvVar)
}

type EnvVarWrongTypeError struct {
	EnvVar string
	Value  any
	Type   string
}

func (e *EnvVarWrongTypeError) Error() string {
	return fmt.Sprintf("Environment variable %v provided as %v, but should be %v", e.EnvVar, e.Value, e.Type)
}

type DatabaseError struct {
	Action      string
	ErrorString string
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("Error returned from database when %v: %v", e.Action, e.ErrorString)
}

type GopherColorInvalidError struct {
	Color string
}

func (e *GopherColorInvalidError) Error() string {
	return fmt.Sprintf("Gopher color is invalid. Gophers can not be %v", e.Color)
}
