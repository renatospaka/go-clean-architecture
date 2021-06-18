package repository

import (
	"errors"
	"log"
	"time"

	"github.com/renatospaka/clean-arch/entity"
	uuid "github.com/satori/go.uuid"
)

type person struct{}

var (
	personInMemory []entity.Person
)

func NewPersonInMemoryRepository() PersonRepository {
	return &person{}
}

func (*person) FindById(id string) (*entity.Person, error) {
	if len(personInMemory) == 0 {
		err := errors.New(entity.ERROR_PERSON_BASE_EMPTY)
		log.Printf("Error: %v", err)
		return &entity.Person{}, err
	}

	thisGuy := entity.Person{}
	for idx, value := range personInMemory {
		if value.ID == id {
			thisGuy = personInMemory[idx]
			break
		}
	} 

	if thisGuy.ID == "" {
		err := errors.New(entity.ERROR_PERSON_INVALID_ID)
		log.Printf("Error: %v", err)
		return &entity.Person{}, err
	}

	return &thisGuy, nil
}

func (*person) Add(person *entity.Person) (*entity.Person, error) {
	err := person.IsValid()
	if err != nil {
		log.Printf("Failed to add a new person: %v", err)
		return &entity.Person{}, err
	}

	person.ID = uuid.NewV4().String()
	person.Responsible = entity.Someone
	person.CreatedAt = time.Now()
	person.UpdatedAt = time.Now()

	personInMemory = append(personInMemory, *person)
	return person, nil
}
