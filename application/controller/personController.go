package controller

import (
	"encoding/json"
	_ "fmt"
	"net/http"

	router "github.com/renatospaka/go-clean-architecture/application/http"
	"github.com/renatospaka/go-clean-architecture/application/service"
	"github.com/renatospaka/go-clean-architecture/entity"
	"github.com/renatospaka/go-clean-architecture/framework/utils"
)

type PersonController interface {
	GetPerson(resp http.ResponseWriter, req *http.Request)
	AddPerson(resp http.ResponseWriter, req *http.Request)
}

type personId struct {
	ID string `json: "person_id"`
}

type controller struct{}

var (
	httpRouter router.Router = router.NewMuxRouter()
)

func NewPersonController() PersonController {
	return &controller{}
}

func (*controller) GetPerson(resp http.ResponseWriter, req *http.Request) {
	personService := service.NewPersonService()

	resp.Header().Set("Content-Type", "application/json")

	id := httpRouter.GetParam(req, "id")
	if id == router.ERROR_MISSING_OR_NOT_FOUND_PARAMETER {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(utils.ServiceError{Message: "Person id is missing"})
		return
	}

	person, err := personService.GetPerson(id)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(utils.ServiceError{Message: "Error getting the person"})
		return
	}

	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(person)
}

func (*controller) AddPerson(resp http.ResponseWriter, req *http.Request) {
	personService := service.NewPersonService()
	var person entity.Person

	resp.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&person)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(utils.ServiceError{Message: "Error unmarshalling the person"})
		return
	}

	newPerson, err := personService.AddPerson(&person)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(utils.ServiceError{Message: err.Error()})
		return
	}

	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(newPerson)
}
