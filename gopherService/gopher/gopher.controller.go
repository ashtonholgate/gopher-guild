package gopher

import (
	"gopherService/customErrors"
	"gopherService/utilities"
	"net/http"
)

type Controller interface {
	CreateGopherEndpoint(r http.ResponseWriter, w *http.Request)
}

type CommandServiceContract interface {
	Create(gopher IncomingGopher) (OutgoingGopher, error)
}

type controllerImpl struct {
	CommandService CommandServiceContract
}

func (cs *controllerImpl) CreateGopherEndpoint(w http.ResponseWriter, r *http.Request) {
	var newGopher IncomingGopher
	if err := utilities.ParseBodyAndValidate(r, &newGopher); err != nil {
		customErrors.Handle(w, err)
		return
	}

	createdGopher, err := cs.CommandService.Create(newGopher)
	if err != nil {
		customErrors.Handle(w, err)
		return
	}
	if err := utilities.WriteJSONResponse(w, http.StatusCreated, createdGopher); err != nil {
		customErrors.Handle(w, err)
	}

}

func NewGopherController(commandService CommandService) Controller {
	return &controllerImpl{CommandService: commandService}
}
