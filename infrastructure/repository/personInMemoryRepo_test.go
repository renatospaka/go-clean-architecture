package repository_test

import (
	"testing"
	"time"

	"github.com/renatospaka/go-clean-architecture/entity"
	"github.com/renatospaka/go-clean-architecture/infrastructure/repository"
	"github.com/stretchr/testify/assert"
)

func TestAdd_IsValid(t *testing.T) {
	personRepo := repository.NewPersonInMemoryRepository()
	person := entity.NewPerson()
	person.Name = "Renato"
	person.MiddleName = "Costa"
	person.LastName = "Spakauskas"
	person.DOB = time.Date(1970, 11, 14, 0, 0, 0, 0, time.UTC)
	person.Email = "renato@email.com"
	person.Responsible = entity.Self

	p, err := personRepo.Add(&person)
	assert.Nil(t, err)
	assert.NotNil(t, p.ID)
	assert.Equal(t, person.ID, p.ID)
}

func TestAdd_IsInvalid(t *testing.T) {
	personRepo := repository.NewPersonInMemoryRepository()
	person := entity.NewPerson()
	person.Name = ""
	person.MiddleName = "Costa"
	person.LastName = "Spakauskas"
	person.DOB = time.Date(1970, 11, 14, 0, 0, 0, 0, time.UTC)
	person.Email = "renato@email.com"
	person.Responsible = entity.Self

	p, err := personRepo.Add(&person)
	assert.NotNil(t, err)
	assert.Equal(t, "name is missing", err.Error())
	assert.Empty(t, p.ID)
}
