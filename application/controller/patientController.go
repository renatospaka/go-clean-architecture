package controller

import (
	"encoding/json"
	"net/http"

	"github.com/renatospaka/go-clean-architecture/application/service"
	"github.com/renatospaka/go-clean-architecture/framework/utils"
	"github.com/renatospaka/go-clean-architecture/entity"
)

type PersonController interface{
	GetPerson(resp http.ResponseWriter, req *http.Request) 
	AddPerson(resp http.ResponseWriter, req *http.Request) 
}

type personId struct {
	ID string `json: "person_id"`
}

type controller struct {}


func NewPersonController() PersonController {
	return &controller{}
}

func (*controller) GetPerson(resp http.ResponseWriter, req *http.Request) {
	personService := service.NewPersonService()
	//thisGuy := entity.Person{}
	var thisId personId

	resp.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(thisId)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(utils.ServiceError{Message: "Error unmarshalling the person id"})
		return 
	}

	thisGuy, err := personService.GetPerson(thisId.ID)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(utils.ServiceError{Message: "Error getting the person"})
		return 
	}

	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(thisGuy)
}

func (*controller) AddPerson(resp http.ResponseWriter, req *http.Request) {
	return
}
