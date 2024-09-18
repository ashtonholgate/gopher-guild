package gopher

import (
	"fmt"
)

type CommandService interface {
	Create(gopher IncomingGopher) (OutgoingGopher, error)
	Read(id int) (OutgoingGopher, error)
}

type RepositoryServiceContract interface {
	Create(IncomingGopher) (OutgoingGopher, error)
	Read(int) (OutgoingGopher, error)
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

func (cs *commandServiceImpl) Read(id int) (OutgoingGopher, error) {
	fetchedGopher, err := cs.repositoryService.Read(id)
	if err != nil {
		return OutgoingGopher{}, fmt.Errorf("gopherCommandService failed to read gopher with id %v: %w", id, err)
	}
	return fetchedGopher, nil
}

func NewGopherCommandService(repositoryService RepositoryServiceContract) CommandService {
	return &commandServiceImpl{repositoryService: repositoryService}
}
