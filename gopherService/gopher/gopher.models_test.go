package gopher

import (
	"encoding/json"
	"gopherService/customErrors"
	"gopherService/utilities"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateGopherModels(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		output          IncomingGopher
		unmarshalError  error
		validationError error
	}{
		{
			name:  "Correct Input",
			input: `{"name": "Gophie", "age": 1, "color": "blue"}`,
			output: IncomingGopher{
				BaseGopher: BaseGopher{
					Name:  "Gophie",
					Age:   utilities.ToPointer(1),
					Color: "blue",
				},
			},
		},
		{
			name:  "Missing Name",
			input: `{"age": 1, "color": "blue"}`,
			output: IncomingGopher{
				BaseGopher: BaseGopher{
					Age:   utilities.ToPointer(1),
					Color: "blue",
				},
			},
			validationError: &customErrors.ValidationError{Issues: "name is required"},
		},
		{
			name:  "Missing Age",
			input: `{"name": "Gophie", "color": "blue"}`,
			output: IncomingGopher{
				BaseGopher: BaseGopher{
					Name:  "Gophie",
					Color: "blue",
				},
			},
			validationError: &customErrors.ValidationError{Issues: "age is required"},
		},
		{
			name:  "Missing Color",
			input: `{"name": "Gophie", "age": 1}`,
			output: IncomingGopher{
				BaseGopher: BaseGopher{
					Name: "Gophie",
					Age:  utilities.ToPointer(1),
				},
			},
			validationError: &customErrors.ValidationError{Issues: "color is required"},
		},
		{
			name:            "Missing All",
			input:           `{}`,
			output:          IncomingGopher{},
			validationError: &customErrors.ValidationError{Issues: "name is required, age is required, color is required"},
		},
		{
			name:  "Age can't be below zero",
			input: `{"name": "Gophie", "age": -1, "color": "blue"}`,
			output: IncomingGopher{
				BaseGopher: BaseGopher{
					Name:  "Gophie",
					Age:   utilities.ToPointer(-1),
					Color: "blue",
				},
			},
			validationError: &customErrors.ValidationError{Issues: "age must be non-negative"},
		},
		{
			name:  "Color can not be red",
			input: `{"name": "Gophie", "age": 1, "color": "red"}`,
			output: IncomingGopher{
				BaseGopher: BaseGopher{
					Name:  "Gophie",
					Age:   utilities.ToPointer(1),
					Color: "red",
				},
			},
			validationError: &customErrors.ValidationError{Issues: "color can not be red"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := IncomingGopher{}
			unmarshalError := json.Unmarshal([]byte(test.input), &got)
			validationError := got.Validate()
			assert.Equal(t, got, test.output)
			assert.Equal(t, unmarshalError, test.unmarshalError)
			assert.Equal(t, validationError, test.validationError)
		})
	}
}
