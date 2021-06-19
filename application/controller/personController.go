package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/renatospaka/go-clean-architecture/application/service"
	"github.com/renatospaka/go-clean-architecture/entity"
	"github.com/renatospaka/go-clean-architecture/framework/utils"
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
	var thisId personId
	//fmt.Println("GetPerson Body:", req.Body)

	resp.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&thisId)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(utils.ServiceError{Message: "Error unmarshalling the person id"})
		return 
	}
	
	fmt.Println("GetPerson thisId.ID:", thisId.ID)
	person, err := personService.GetPerson(string(thisId.ID))
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
