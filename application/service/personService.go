package service

import (
	//"gopkg.in/multierror.v1"

	"errors"

	"github.com/renatospaka/go-clean-architecture/entity"
	"github.com/renatospaka/go-clean-architecture/infrastrucure/repository"
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

func (*person) GetPerson(id string) (*entity.Person, error) {
	thisGuy, err := repo.FindById(id)
	if err != nil {
		desc := "personrepository.getperson.error: failed to retrieve person\n" + err.Error()
		err2 := errors.New(desc)
		return &entity.Person{}, err2
	}

	return thisGuy, nil
}

func (*person) AddPerson(person *entity.Person) (*entity.Person, error) {
	thisGuy, err := repo.Add(person)
	if err != nil {
		desc := "personrepository.add.error: failed to add this person\n" + err.Error()
		err2 := errors.New(desc)
		return &entity.Person{}, err2
	}
	return thisGuy, nil
}
