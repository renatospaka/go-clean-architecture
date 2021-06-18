package main

import (
	"fmt"
	"log"
	"time"

	"github.com/renatospaka/clean-arch/application/repository"
	"github.com/renatospaka/clean-arch/entity"
)

var (
	personRepo repository.PersonRepository = repository.NewPersonInMemoryRepository()

)
func main() {
	log.Println("EMR Server up & running")

	person := entity.NewPerson()
	person.Name = "Renato"
	person.MiddleName = "Costa"
	person.LastName = "Spakauskas"
	person.DOB = time.Date(1970, 11, 14, 0, 0, 0, 0, time.UTC)
	person.Email = "renato@email.com"
	person.Responsible = entity.Self
	
	//add new person
	p, err := personRepo.Add(&person)
	if err != nil {
		log.Fatalf("Error on add a new person: %v", err)
	}
	fmt.Println("NEW => ", p.ID, ", ", p.Name)
	
	//consult this new person
	existingPerson, err := personRepo.FindById(p.ID)
	if err != nil {
		log.Fatalf("Error consulting an ID: %v", err)
	}
	fmt.Println("EXISTING => ", existingPerson.ID, ", ", existingPerson.Name, ", included on ", existingPerson.CreatedAt)
}
