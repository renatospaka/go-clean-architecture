package repository

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/renatospaka/go-clean-architecture/entity"
	uuid "github.com/satori/go.uuid"
)

type person struct{}

var (
	personInMemory []entity.Person
)

func init() {
	thisPerson := entity.Person{
		ID: "1",
		Name: "Renato",
		MiddleName: "Costa",
		LastName: "Spakauskas",
		DOB: time.Date(1970, 11, 14, 0, 0, 0, 0, time.UTC),
		Gender: entity.GenderMale,
		Email: "renato@email.com",
	}
	personInMemory = append(personInMemory, thisPerson)

	thisPerson = entity.Person{
		ID: "2",
		Name: "Camila",
		MiddleName: "Pinho",
		LastName: "Spakauskas",
		DOB: time.Date(1995, 2, 6, 0, 0, 0, 0, time.UTC),
		Gender: entity.GenderFemale,
		Email: "camila@email.com",
	}
	personInMemory = append(personInMemory, thisPerson)
	
	fmt.Println("personInMemory: ", len(personInMemory))
}

func NewPersonInMemoryRepository() PersonRepository {
	return &person{}
}

func (*person) FindById(id string) (*entity.Person, error) {
	if len(personInMemory) == 0 {
		err := errors.New(entity.ERROR_PERSON_BASE_EMPTY)
		log.Printf("person.findbyid.error: %v", err)
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
		log.Printf("person.findbyid.error: %v", err)
		return &entity.Person{}, err
	}

	return &thisGuy, nil
}

func (*person) Add(person *entity.Person) (*entity.Person, error) {
	err := person.IsValid()
	if err != nil {
		log.Printf("person.add.error: %v", err)
		return &entity.Person{}, err
	}

	person.ID = uuid.NewV4().String()
	person.Responsible = entity.Someone
	person.CreatedAt = time.Now()
	person.UpdatedAt = time.Now()

	personInMemory = append(personInMemory, *person)
	return person, nil
}
