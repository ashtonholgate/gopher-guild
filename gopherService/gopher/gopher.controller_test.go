package gopher

import (
	"bytes"
	"encoding/json"
	"gopherService/customErrors"
	"gopherService/utilities"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCommandService is a mock of CommandService interface
type MockCommandService struct {
	mock.Mock
}

func (m *MockCommandService) Create(gopher IncomingGopher) (OutgoingGopher, error) {
	args := m.Called(gopher)
	return args.Get(0).(OutgoingGopher), args.Error(1)
}

func TestCreateGopherEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		inputGopher    IncomingGopher
		mockReturn     OutgoingGopher
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "Successful Gopher Creation",
			inputGopher: IncomingGopher{BaseGopher: BaseGopher{
				Name:  "TestGopher",
				Age:   utilities.ToPointer(5),
				Color: "Blue",
			}},
			mockReturn:     OutgoingGopher{BaseGopher: BaseGopher{Name: "TestGopher", Age: utilities.ToPointer(5), Color: "Blue"}, Id: 1},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedBody:   OutgoingGopher{BaseGopher: BaseGopher{Name: "TestGopher", Age: utilities.ToPointer(5), Color: "Blue"}, Id: 1},
		},
		{
			name: "Database Error",
			inputGopher: IncomingGopher{BaseGopher: BaseGopher{
				Name:  "TestGopher",
				Age:   utilities.ToPointer(5),
				Color: "Blue",
			}},
			mockReturn: OutgoingGopher{},
			mockError: &customErrors.DatabaseError{
				Action:      "creating gopher",
				ErrorString: "connection timeout",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: gin.H{
				"error":   "Database operation failed",
				"details": "An error occurred while processing your request. Please try again later.",
			},
		},
		{
			name: "Unexpected Error",
			inputGopher: IncomingGopher{BaseGopher: BaseGopher{
				Name:  "TestGopher",
				Age:   utilities.ToPointer(5),
				Color: "Blue",
			}},
			mockReturn:     OutgoingGopher{},
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   gin.H{"error": "An unexpected error occurred", "details": "Please try again later or contact support if the problem persists."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCommandService)
			mockService.On("Create", tt.inputGopher).Return(tt.mockReturn, tt.mockError)

			controller := NewGopherController(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.inputGopher)
			c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			controller.CreateGopherEndpoint()(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedStatus == http.StatusCreated {
				assert.Equal(t, tt.expectedBody.(OutgoingGopher).Id, int(response["id"].(float64)))
				assert.Equal(t, tt.expectedBody.(OutgoingGopher).Name, response["name"])
				assert.Equal(t, tt.expectedBody.(OutgoingGopher).Color, response["color"])
				assert.Equal(t, float64(*tt.expectedBody.(OutgoingGopher).Age), response["age"])
			} else {
				assert.Equal(t, tt.expectedBody.(gin.H)["error"], response["error"])
				assert.Equal(t, tt.expectedBody.(gin.H)["details"], response["details"])
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestInvalidJSONInput(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockCommandService)
	controller := NewGopherController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	invalidJSON := []byte(`{"name": "TestGopher", "age": "invalid", "color": "Blue"}`)
	c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.CreateGopherEndpoint()(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "json: cannot unmarshal")
}
