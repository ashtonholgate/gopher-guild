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

type NoRowsError struct{}

func (e *NoRowsError) Error() string {
	return "Not found"
}

type DatabaseError struct {
	Action      string
	ErrorString string
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("Error returned from database when %v: %v", e.Action, e.ErrorString)
}

type URLParsingError struct {
	PathParam string
}

func (e *URLParsingError) Error() string {
	return fmt.Sprintf("Could not find %v in path", e.PathParam)
}

type JSONDecodingError struct {
	Err error
}

func (e *JSONDecodingError) Error() string {
	return e.Err.Error()
}

type ValidationError struct {
	Issues string
}

func (e *ValidationError) Error() string {
	return e.Issues
}
