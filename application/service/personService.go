package service

import (
	//"gopkg.in/multierror.v1"

	"errors"

	"github.com/renatospaka/go-clean-architecture/entity"
	"github.com/renatospaka/go-clean-architecture/infrastructure/repository"
)

type PersonService interface {
	GetPerson(id string) (*entity.Person, error)
	AddPerson(person *entity.Person) (*entity.Person, error)
}

type person struct{}

var repo repository.PersonRepository = repository.NewPersonInMemoryRepository()

func NewPersonService() PersonService {
	return &person{}
}


//Get one person identified by his/her ID
func (*person) GetPerson(id string) (*entity.Person, error) {
	thisGuy, err := repo.FindById(id)
	if err != nil {
		desc := "personrepository.getperson.error: " + err.Error()
		err2 := errors.New(desc)
		return &entity.Person{}, err2
	}

	return thisGuy, nil
}


//add a new person who later can be a patient, a responsible for someone, a medic
func (*person) AddPerson(person *entity.Person) (*entity.Person, error) {
	thisGuy, err := repo.Add(person)
	if err != nil {
		desc := "personrepository.addperson.error: " + err.Error()
		err2 := errors.New(desc)
		return &entity.Person{}, err2
	}
	return thisGuy, nil
}
