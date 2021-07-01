package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/renatospaka/go-clean-architecture/application/adapter/router"
	"github.com/renatospaka/go-clean-architecture/application/service"
	"github.com/renatospaka/go-clean-architecture/core/entity/person"
	"github.com/renatospaka/go-clean-architecture/infrastructure/utils"
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
	httpRouter    router.Router = router.NewMuxRouter()
	personService service.PersonService
)

func NewPersonController(service service.PersonService) PersonController {
	personService = service
	return &controller{}
}

//Get a person by ID. If ID is missing or not found, then an error is raised
func (*controller) GetPerson(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	id := httpRouter.GetParam(req, "id")
	if id == utils.ERROR_MISSING_PARAMETER {
		log.Println("personController.GetPerson: person id is missing")
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(utils.ServiceError{Message: "person id is missing"})
		return
	}
	
	person, err := personService.GetPerson(id)
	if err != nil {
		log.Printf("personController.GetPerson: %v", err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(utils.ServiceError{Message: err.Error()})
		return
	}
	
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(person)
}

//Add a person. For now, just expecting JSON parse partial object based on struct Person
func (*controller) AddPerson(resp http.ResponseWriter, req *http.Request) {
	var person entity.Person
	
	resp.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&person)
	if err != nil {
		log.Printf("personController.AddPerson: %v", err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(utils.ServiceError{Message: err.Error()})
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
