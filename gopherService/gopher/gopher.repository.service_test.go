package gopher

import (
	"database/sql"
	"errors"
	"gopherService/customErrors"
	"gopherService/utilities"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRepositoryServiceImpl_Create(t *testing.T) {
	// Create a new SQL mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new repository service with the mocked db
	repoService := NewGopherRepositoryService(db)

	tests := []struct {
		name          string
		input         IncomingGopher
		mockBehavior  func(mock sqlmock.Sqlmock, input IncomingGopher)
		expected      OutgoingGopher
		expectedError error
	}{
		{
			name: "Successful Gopher Creation",
			input: IncomingGopher{
				BaseGopher: BaseGopher{
					Name:  "TestGopher",
					Age:   utilities.ToPointer(5),
					Color: "Blue",
				},
			},
			mockBehavior: func(mock sqlmock.Sqlmock, input IncomingGopher) {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO gophers").
					WithArgs(input.Name, input.Age, input.Color).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "age", "color"}).
						AddRow(1, input.Name, input.Age, input.Color))
				mock.ExpectCommit()
			},
			expected: OutgoingGopher{
				BaseGopher: BaseGopher{
					Name:  "TestGopher",
					Age:   utilities.ToPointer(5),
					Color: "Blue",
				},
				Id: 1,
			},
			expectedError: nil,
		},
		{
			name: "Database Error",
			input: IncomingGopher{
				BaseGopher: BaseGopher{
					Name:  "ErrorGopher",
					Age:   utilities.ToPointer(3),
					Color: "Red",
				},
			},
			mockBehavior: func(mock sqlmock.Sqlmock, input IncomingGopher) {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO gophers").
					WithArgs(input.Name, input.Age, input.Color).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expected:      OutgoingGopher{},
			expectedError: &customErrors.DatabaseError{Action: "inserting into gophers table", ErrorString: sql.ErrConnDone.Error()},
		},
		{
			name: "Transaction Begin Error",
			input: IncomingGopher{
				BaseGopher: BaseGopher{
					Name:  "TransactionErrorGopher",
					Age:   utilities.ToPointer(4),
					Color: "Green",
				},
			},
			mockBehavior: func(mock sqlmock.Sqlmock, input IncomingGopher) {
				mock.ExpectBegin().WillReturnError(sql.ErrConnDone)
			},
			expected:      OutgoingGopher{},
			expectedError: sql.ErrConnDone,
		},
		{
			name: "Transaction Commit Error",
			input: IncomingGopher{
				BaseGopher: BaseGopher{
					Name:  "CommitErrorGopher",
					Age:   utilities.ToPointer(6),
					Color: "Yellow",
				},
			},
			mockBehavior: func(mock sqlmock.Sqlmock, input IncomingGopher) {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO gophers").
					WithArgs(input.Name, input.Age, input.Color).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "age", "color"}).
						AddRow(1, input.Name, input.Age, input.Color))
				mock.ExpectCommit().WillReturnError(sql.ErrTxDone)
			},
			expected:      OutgoingGopher{},
			expectedError: sql.ErrTxDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the mock behavior for this test case
			tt.mockBehavior(mock, tt.input)

			// Call the method being tested
			result, err := repoService.Create(tt.input)

			// Check the results
			if tt.expectedError != nil {
				assert.Error(t, err)
				var dbErr *customErrors.DatabaseError
				if errors.As(err, &dbErr) {
					assert.Equal(t, tt.expectedError.(*customErrors.DatabaseError).Action, dbErr.Action)
					assert.Contains(t, dbErr.ErrorString, tt.expectedError.(*customErrors.DatabaseError).ErrorString)
				} else {
					assert.ErrorIs(t, err, tt.expectedError)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			// Ensure all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
