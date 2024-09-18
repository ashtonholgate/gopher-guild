package gopher

import (
	"gopherService/customErrors"
	"gopherService/utilities"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Controller interface {
	CreateGopherEndpoint(r http.ResponseWriter, w *http.Request)
	ReadGopherEndpoint(r http.ResponseWriter, w *http.Request)
}

type CommandServiceContract interface {
	Create(gopher IncomingGopher) (OutgoingGopher, error)
	Read(id int) (OutgoingGopher, error)
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

func (cs *controllerImpl) ReadGopherEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		customErrors.Handle(w, &customErrors.URLParsingError{PathParam: "id"})
		return
	}

	intId, err := strconv.Atoi(id)
	if err != nil {
		customErrors.Handle(w, err)
		return
	}

	fetchedGopher, err := cs.CommandService.Read(intId)
	if err != nil {
		customErrors.Handle(w, err)
		return
	}

	if err := utilities.WriteJSONResponse(w, http.StatusCreated, fetchedGopher); err != nil {
		customErrors.Handle(w, err)
	}
}

func NewGopherController(commandService CommandService) Controller {
	return &controllerImpl{CommandService: commandService}
}
