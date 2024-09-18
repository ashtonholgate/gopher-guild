package customErrors

import (
	"errors"
	"log"
	"net/http"
)

func Handle(w http.ResponseWriter, err error) {
	log.Printf("%+v", err)

	var JSONDecodingError *JSONDecodingError
	var validationError *ValidationError
	var noRowsError *NoRowsError
	var databaseErr *DatabaseError

	var statusCode int
	var message string

	switch {
	case errors.As(err, &JSONDecodingError):
		statusCode = http.StatusBadRequest
		message = JSONDecodingError.Error()
	case errors.As(err, &validationError):
		statusCode = http.StatusBadRequest
		message = validationError.Error()
	case errors.As(err, &noRowsError):
		statusCode = http.StatusNotFound
	case errors.As(err, &databaseErr):
		statusCode = http.StatusInternalServerError
		message = "Database operation failed."
	default:
		statusCode = http.StatusInternalServerError
		message = "An unexpected error occurred. Please try again later or contact support if the problem persists."
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}
