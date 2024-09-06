package gopher

import (
	"gopherService/customErrors"
	"gopherService/utilities"
	"reflect"
	"testing"
)

type MockRepositoryService struct {
	CreateFunc func(gopher IncomingGopher) (OutgoingGopher, error)
}

func (m *MockRepositoryService) Create(gopher IncomingGopher) (OutgoingGopher, error) {
	return m.CreateFunc(gopher)
}

func TestValidateGopher(t *testing.T) {
	tests := []struct {
		name    string
		gopher  IncomingGopher
		wantErr bool
	}{
		{
			name: "Valid gopher",
			gopher: IncomingGopher{
				BaseGopher: BaseGopher{
					Name:  "Gophie",
					Age:   utilities.ToPointer(5),
					Color: "Blue",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid color (red)",
			gopher: IncomingGopher{
				BaseGopher: BaseGopher{
					Name:  "Redy",
					Age:   utilities.ToPointer(3),
					Color: "Red",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validateGopher(tt.gopher)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateGopher() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if _, ok := err.(*customErrors.GopherColorInvalidError); !ok {
					t.Errorf("validateGopher() error is not of type GopherColorInvalidError")
				}
			}
		})
	}
}

func TestCommandServiceImpl_Create(t *testing.T) {
	tests := []struct {
		name           string
		gopher         IncomingGopher
		mockCreateFunc func(gopher IncomingGopher) (OutgoingGopher, error)
		want           OutgoingGopher
		wantErr        bool
	}{
		{
			name: "Successful creation",
			gopher: IncomingGopher{
				BaseGopher: BaseGopher{
					Name:  "Gophie",
					Age:   utilities.ToPointer(5),
					Color: "Blue",
				},
			},
			mockCreateFunc: func(gopher IncomingGopher) (OutgoingGopher, error) {
				return OutgoingGopher{
					BaseGopher: gopher.BaseGopher,
					Id:         1,
				}, nil
			},
			want: OutgoingGopher{
				BaseGopher: BaseGopher{
					Name:  "Gophie",
					Age:   utilities.ToPointer(5),
					Color: "Blue",
				},
				Id: 1,
			},
			wantErr: false,
		},
		{
			name: "Invalid color",
			gopher: IncomingGopher{
				BaseGopher: BaseGopher{
					Name:  "Redy",
					Age:   utilities.ToPointer(3),
					Color: "Red",
				},
			},
			mockCreateFunc: func(gopher IncomingGopher) (OutgoingGopher, error) {
				return OutgoingGopher{}, nil
			},
			want:    OutgoingGopher{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepositoryService{
				CreateFunc: tt.mockCreateFunc,
			}
			cs := &commandServiceImpl{
				repositoryService: mockRepo,
			}
			got, err := cs.Create(tt.gopher)
			if (err != nil) != tt.wantErr {
				t.Errorf("commandServiceImpl.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("commandServiceImpl.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGopherCommandService(t *testing.T) {
	mockRepo := &MockRepositoryService{}
	cs := NewGopherCommandService(mockRepo)

	if cs == nil {
		t.Error("NewGopherCommandService() returned nil")
	}

	if _, ok := cs.(*commandServiceImpl); !ok {
		t.Error("NewGopherCommandService() did not return a *commandServiceImpl")
	}
}
