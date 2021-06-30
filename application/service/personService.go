package service

import (
	"log"

	"github.com/renatospaka/go-clean-architecture/core/entity/person"
	"github.com/renatospaka/go-clean-architecture/core/repository/person"
)

type PersonService interface {
	GetPerson(id string) (*entity.Person, error)
	AddPerson(person *entity.Person) (*entity.Person, error)
}

type person struct{}

//var repo repository.PersonRepository = repository.NewPersonInMemoryRepository()
var repo repository.PersonRepository //= repository.NewPersonNeo4jRepository()

func NewPersonService() PersonService {
	repo = repository.NewPersonNeo4jRepository()
	return &person{}
}


//Get one person identified by his/her ID
func (*person) GetPerson(id string) (*entity.Person, error) {
	thisGuy, err := repo.FindById(id)
	if err != nil {
		log.Printf("personController.GetPerson.error: %v", err.Error())
		// desc := "personController.GetPerson.error: " + err.Error()
		// err2 := errors.New(desc)
		return &entity.Person{}, err
	}
	
	return thisGuy, nil
}


//add a new person who later can be a patient, a responsible for someone, a medic
func (*person) AddPerson(person *entity.Person) (*entity.Person, error) {
	thisGuy, err := repo.Add(person)
	if err != nil {
		// desc := "personController.AddPerson.error: " + err.Error()
		// // err2 := errors.New(desc)
		return &entity.Person{}, err
	}
	return thisGuy, nil
}
