package repository

import "github.com/renatospaka/go-clean-architecture/entity"

type PersonRepository interface {
	FindById(id string) (*entity.Person, error)
	Add(person *entity.Person) (*entity.Person, error)
	//Update(person *entity.Person) (*entity.Person, error)
	//Remove(id string) (*entity.Person, error)
}
