package repository

import (
	"errors"
	"log"
	"time"

	"github.com/renatospaka/go-clean-architecture/entity/person"
	uuid "github.com/satori/go.uuid"
)

type personInMemory struct{}

var (
	personInMemoryDb []entity.Person
)

func NewPersonInMemoryRepository() PersonRepository {
	return &personInMemory{}
}

func (*personInMemory) FindById(id string) (*entity.Person, error) {
	if len(personInMemoryDb) == 0 {
		err := errors.New(entity.ERROR_PERSON_BASE_EMPTY)
		log.Printf("person.findbyid.error: %v", err)
		return &entity.Person{}, err
	}

	thisGuy := entity.Person{}
	for idx, value := range personInMemoryDb {
		if value.ID == id {
			thisGuy = personInMemoryDb[idx]
			break
		}
	} 
		
	if thisGuy.ID == "" {
		err := errors.New(entity.ERROR_PERSON_INVALID_ID)
		log.Printf("person.findbyid.error: %v", err)
		return &entity.Person{}, err
	}

	return &thisGuy, nil
}

func (*personInMemory) Add(newPerson *entity.Person) (*entity.Person, error) {
	err := newPerson.IsValid()
	if err != nil {
		log.Printf("person.add.error: %v", err)
		return &entity.Person{}, err
	}

	newPerson.ID = uuid.NewV4().String()
	newPerson.Responsible = entity.Someone
	newPerson.CreatedAt = time.Now()
	newPerson.UpdatedAt = time.Now()

	personInMemoryDb = append(personInMemoryDb, *newPerson)
	return newPerson, nil
}
