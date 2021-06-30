package entity

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"

	"github.com/renatospaka/go-clean-architecture/infrastructure/utils"
)

type Person struct {
	ID          string    `json: "person_id" valid:"uuid,optional"`
	Name        string    `json: "name" valid:"alpha,required,stringlength(3|30)"`
	MiddleName  string    `json: "middle_name" valid:"alpha,optional"`
	LastName    string    `json: "last_name" valid:"alpha,required,stringlength(3|30)"`
	Email       string    `json: "email" valid:"email,required"`
	Gender      string    `json: "gender" valid:"alpha,optional"`
	DOB         time.Time `json: "day_of_birth" valid:"-"`
	AgeInMonths int       `json: "age_in_months" valid:"-"`
	Age         int       `json: "age" valid:"-"`
	Responsible string    `json: "responsible_id" valid:"-"`
	CreatedAt   time.Time `json: "created_at" valid:"-"`
	UpdatedAt   time.Time `json: "updated_at" valid:"-"`
}

type person struct{}

func NewPerson() Person {
	p := Person{}
	return p
}

func init() {
	govalidator.SetFieldsRequiredByDefault(false)
}

//check wheather person is filled with required information - name, lastname, email and day of birth
func (p *Person) IsValid() error {
	if govalidator.IsNull(p.Name) {
	//if strings.TrimSpace(p.Name) == "" {
		return errors.New(ERROR_NAME_MISSING)
	}

	if !govalidator.MinStringLength(p.Name, "3") {
	//if len(strings.TrimSpace(p.Name)) < 3 {
		return errors.New(ERROR_NAME_TOO_SHORT)
	}

	if !govalidator.MaxStringLength(p.Name, "20") {
	//if len(strings.TrimSpace(p.Name)) > 20 {
		return errors.New(ERROR_NAME_TOO_LONG)
	}

	if govalidator.IsNull(p.LastName) {
	//if strings.TrimSpace(p.LastName) == "" {
		return errors.New(ERROR_LAST_NAME_MISSING)
	}

	if !govalidator.MinStringLength(p.LastName, "3") {
	//if len(strings.TrimSpace(p.LastName)) < 3 {
		return errors.New(ERROR_LAST_NAME_TOO_SHORT)
	}

	if !govalidator.MaxStringLength(p.LastName, "20") {
	//if len(strings.TrimSpace(p.LastName)) > 20 {
		return errors.New(ERROR_LAST_NAME_TOO_LONG)
	}

	if !govalidator.StringLength(p.MiddleName, "0", "30") {
		return errors.New(ERROR_MIDDLE_NAME_TOO_LONG)
	}

	if p.Responsible == Self {
		if govalidator.IsNull(p.Email) {
		//if strings.TrimSpace(p.Email) == "" {
			return errors.New(ERROR_EMAIL_MISSING)
		}

		if !govalidator.IsEmail(p.Email) {
		//if !utils.IsEmailValid(p.Email) {
			return errors.New(ERROR_EMAIL_INVALID)
		}
	}

	if p.DOB.IsZero() {
		return errors.New(ERROR_DOB_MISSING)
	}

	if utils.IsDateEqualToday(p.DOB) {
		return errors.New(ERROR_DOB_INVALID)
	}

	if utils.IsDateGreaterToday(p.DOB) {
		return errors.New(ERROR_DOB_INVALID)
	}

	return nil
}

func (p *Person) isValid() error {
	_, err := govalidator.ValidateStruct(p)
	if err != nil {
		return err
	}
	return nil
}