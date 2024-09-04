package gopher

import (
	"fmt"
	"gopherService/customErrors"
	"strings"
)

type CommandService interface {
	Create(gopher IncomingGopher) (OutgoingGopher, error)
}

type RepositoryServiceContract interface {
	Create(IncomingGopher) (OutgoingGopher, error)
}

type commandServiceImpl struct {
	repositoryService RepositoryServiceContract
}

func validateGopher(gopher IncomingGopher) (IncomingGopher, error) {
	if strings.ToLower(gopher.Color) == "red" {
		return IncomingGopher{}, &customErrors.GopherColorInvalidError{Color: gopher.Color}
	}
	return gopher, nil
}

func (cs *commandServiceImpl) Create(gopher IncomingGopher) (OutgoingGopher, error) {
	if _, err := validateGopher(gopher); err != nil {
		return OutgoingGopher{}, err
	}
	createdGopher, err := cs.repositoryService.Create(gopher)
	if err != nil {
		return OutgoingGopher{}, fmt.Errorf("gopherCommandService failed to create gopher: %w", err)
	}
	return createdGopher, nil
}

func NewGopherCommandService(repositoryService RepositoryServiceContract) CommandService {
	return &commandServiceImpl{repositoryService: repositoryService}
}
