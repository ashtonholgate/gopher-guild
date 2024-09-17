package utilities

import (
	"encoding/json"
	"gopherService/customErrors"
	"net/http"
)

type validateableStruct interface {
	Validate() error
}

// ParseBodyAndValidate extracts the body of the HTTP request, unmarshals it, and runs the validation function that it expects to find on the struct that it is given.
func ParseBodyAndValidate(r *http.Request, newStruct validateableStruct) error {
	if err := json.NewDecoder(r.Body).Decode(&newStruct); err != nil {
		return &customErrors.JSONDecodingError{Err: err}
	}
	if err := newStruct.Validate(); err != nil {
		return err
	}
	return nil
}

// WriteJSONResponse will marshal a struct that it is given and send it in the body of a HTTP request with the status code provided.
func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}
