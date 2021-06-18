package entity_test

import (
	"testing"
	"time"

	"github.com/renatospaka/go-clean-architecture/entity"
	"github.com/stretchr/testify/require"
)

func TestPerson_IsValid(t *testing.T) {
	person := entity.NewPerson()
	person.Name = "Renato"
	person.MiddleName = "Costa"
	person.LastName = "Spakauskas"
	person.DOB = time.Date(1970, 11, 14, 0, 0, 0, 0, time.UTC)
	person.Email = "renato@email.com"
	person.Responsible = entity.Self

	//all working
	err := person.IsValid()
	require.Nil(t, err)
}

func TestPerson_NameInvalid(t *testing.T) {
	person := entity.NewPerson()
	person.Name = ""
	person.MiddleName = "Costa"
	person.LastName = "Spakauskas"
	person.DOB = time.Date(1970, 11, 14, 0, 0, 0, 0, time.UTC)
	person.Email = "renato@email.com"
	person.Responsible = entity.Self

	err := person.IsValid()
	require.EqualError(t, err, entity.ERROR_NAME_MISSING)

	person.Name = "Re"
	err = person.IsValid()
	require.EqualError(t, err, entity.ERROR_NAME_TOO_SHORT)

	person.Name = "RenatoRenatoRenatoRenatoRenatoRenatoRenatoRenatoRenato"
	err = person.IsValid()
	require.EqualError(t, err, entity.ERROR_NAME_TOO_LONG)

	person.Name = "Renato"
	person.LastName = ""
	err = person.IsValid()
	require.EqualError(t, err, entity.ERROR_LAST_NAME_MISSING)

	person.LastName = "Sp"
	err = person.IsValid()
	require.EqualError(t, err, entity.ERROR_LAST_NAME_TOO_SHORT)

	person.LastName = "SpakauskasSpakauskasSpakauskasSpakauskasSpakauskasSpakauskas"
	err = person.IsValid()
	require.EqualError(t, err, entity.ERROR_LAST_NAME_TOO_LONG)
}

func TestPerson_EmailInvalid(t *testing.T) {
	person := entity.NewPerson()
	person.Name = "Renato"
	person.MiddleName = "Costa"
	person.LastName = "Spakauskas"
	person.DOB = time.Date(1970, 11, 14, 0, 0, 0, 0, time.UTC)
	person.Email = ""
	person.Responsible = entity.Self

	err := person.IsValid()
	require.EqualError(t, err, entity.ERROR_EMAIL_MISSING)

	person.Email = "email"
	err = person.IsValid()
	require.EqualError(t, err, entity.ERROR_EMAIL_INVALID)

	person.Email = "email@"
	err = person.IsValid()
	require.EqualError(t, err, entity.ERROR_EMAIL_INVALID)

	person.Email = "em@a"
	err = person.IsValid()
	require.EqualError(t, err, entity.ERROR_EMAIL_INVALID)

	person.Responsible = ""
	person.Email = ""
	err = person.IsValid()
	require.Nil(t, err)

	person.Email = "email"
	err = person.IsValid()
	require.Nil(t, err)

	person.Email = "email@"
	err = person.IsValid()
	require.Nil(t, err)

	person.Email = "em@a"
	err = person.IsValid()
	require.Nil(t, err)
}

func TestPerson_DOBInvalid(t *testing.T) {
	person := entity.NewPerson()
	person.Name = "Renato"
	person.MiddleName = "Costa"
	person.LastName = "Spakauskas"
	person.DOB = time.Time{}
	person.Email = "renato@email.com"
	person.Responsible = entity.Self

	err := person.IsValid()
	require.EqualError(t, err, entity.ERROR_DOB_MISSING)

	person.DOB = time.Now()
	err = person.IsValid()
	require.EqualError(t, err, entity.ERROR_DOB_INVALID)

	person.DOB = time.Now().AddDate(0, 0, 1) //amanhã
	err = person.IsValid()
	require.EqualError(t, err, entity.ERROR_DOB_INVALID)

	person.DOB = time.Now().AddDate(-5, 3, -12) //5 anos antes
	err = person.IsValid()
	require.Nil(t, err)
}