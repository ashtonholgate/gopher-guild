package gopher

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCommandService is a mock implementation of CommandServiceContract
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
		expectedBody   OutgoingGopher
	}{
		{
			name: "Successful creation",
			inputGopher: IncomingGopher{
				BaseGopher: BaseGopher{
					Name:  "TestGopher",
					Age:   new(int),
					Color: "Brown",
				},
			},
			mockReturn: OutgoingGopher{
				BaseGopher: BaseGopher{
					Name:  "TestGopher",
					Age:   new(int),
					Color: "Brown",
				},
				Id: 1,
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedBody: OutgoingGopher{
				BaseGopher: BaseGopher{
					Name:  "TestGopher",
					Age:   new(int),
					Color: "Brown",
				},
				Id: 1,
			},
		},
		{
			name: "Creation error",
			inputGopher: IncomingGopher{
				BaseGopher: BaseGopher{
					Name:  "TestGopher",
					Age:   new(int),
					Color: "Brown",
				},
			},
			mockReturn:     OutgoingGopher{},
			mockError:      errors.New("creation error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   OutgoingGopher{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCommandService)
			mockService.On("Create", tt.inputGopher).Return(tt.mockReturn, tt.mockError)

			controller := NewGopherController(mockService)

			router := gin.New()
			router.POST("/gophers", controller.CreateGopherEndpoint())
			router.Use(ErrorMiddleware())

			w := httptest.NewRecorder()
			jsonData, _ := json.Marshal(tt.inputGopher)
			req, _ := http.NewRequest("POST", "/gophers", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.mockError == nil {
				var response OutgoingGopher
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			} else {
				assert.Contains(t, w.Body.String(), tt.mockError.Error())
			}

			mockService.AssertExpectations(t)
		})
	}
}

// ErrorMiddleware handles errors and returns appropriate responses
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}
}
