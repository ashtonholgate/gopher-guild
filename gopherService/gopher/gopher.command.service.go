package gopher

import (
	"fmt"
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

func (cs *commandServiceImpl) Create(gopher IncomingGopher) (OutgoingGopher, error) {
	createdGopher, err := cs.repositoryService.Create(gopher)
	if err != nil {
		return OutgoingGopher{}, fmt.Errorf("gopherCommandService failed to create gopher: %w", err)
	}
	return createdGopher, nil
}

func NewGopherCommandService(repositoryService RepositoryServiceContract) CommandService {
	return &commandServiceImpl{repositoryService: repositoryService}
}
