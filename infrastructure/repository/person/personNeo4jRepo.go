package repository

import (
	"errors"
	"log"
	"time"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	uuid "github.com/satori/go.uuid"

	"github.com/renatospaka/go-clean-architecture/entity/person"
	"github.com/renatospaka/go-clean-architecture/infrastructure/db"
)

type personNeo4j struct{}

var (
	personNeo4jDbArray []entity.Person
	personNeo4jDb db.Neo4jSession //= db.NewNeo4jSessionWrite()
)

func NewPersonNeo4jRepository() PersonRepository {
	return &personNeo4j{}
}

func (*personNeo4j) FindById(id string) (*entity.Person, error) {
	personNeo4jDb = db.NewNeo4jSession()
	if err := personNeo4jDb.Connect(db.AccessRead); err != nil {
		log.Printf("person.FindById.Connect.error: %v", err)
		return &entity.Person{}, err
	}

	if err := personNeo4jDb.IsValid(); err != nil {
		log.Printf("person.FindById.IsValid.error: %v", err)
		return &entity.Person{}, err
	}
	defer personNeo4jDb.Session.Close()
	
	tx, err := personNeo4jDb.Session.BeginTransaction()
	if err != nil {
		log.Printf("person.FindById.BeginTransaction.error: %v", err)
		return &entity.Person{}, err
	}
	defer tx.Rollback()
	
	cypher := "MATCH (one:Person {uuid: $uuid}) " +
						"WITH one " +
						"MATCH (one)-[:BIRTH]->(dob:Day) " +
						"RETURN one.uuid As uuid, one.name AS name, one.middleName AS middleName, one.lastName AS lastName, one.gender AS gender, one.email AS email, dob.date AS dob "
	params := map[string]interface{}{
		"uuid": id,
	}
	result, err := tx.Run(cypher, params)
	if err != nil {
		log.Printf("person.FindById.Run.error: %v", err)
		return &entity.Person{}, err
	}
	
	found := false
	thisGuy := entity.Person{}
	for result.Next() {
		var dob2 dbtype.Date
		record := result.Record()
		uuid, ok := record.Get("uuid")
		name, ok := record.Get("name")
		middleName, ok := record.Get("middleName")
		lastName, ok := record.Get("lastName")
		email, ok := record.Get("email")
		if !ok || email == nil {
			email = ""
		}
		gender, ok := record.Get("gender")
		if !ok || gender == nil {
			gender = ""
		}
		//here is a little tricky
		dob, ok := record.Get("dob")
		if dob != nil && ok {
			dob2 = dob.(dbtype.Date)
		}
		if !ok || dob == nil { 
			dob2 = dbtype.Date{} 
		}

		found = true
		thisGuy = entity.Person{
			ID:         uuid.(string),
			Name:       name.(string),
			MiddleName: middleName.(string),
			LastName:   lastName.(string),
			Email:      email.(string),
			Gender:     gender.(string),
			DOB:        time.Time(dob2),	// that gave me 2 entire days of a huge headache
		}
	}

	if thisGuy.ID == "" || !found {
		err := errors.New(entity.ERROR_PERSON_INVALID_ID)
		log.Printf("person.FindById.error: %v", err)
		// err2 := tx.Rollback()
		// if err2 != nil {
		// 	log.Printf("person.FindById.result.Single.Rollback.error: %v", err2)
		// }
		return &entity.Person{}, err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("person.FindById.Commit.error: %v", err)
		return &entity.Person{}, err
	}

	err = tx.Close()
	if err != nil {
		log.Printf("person.FindById.Close.error: %v", err)
		return &entity.Person{}, err
	}

	personNeo4jDb.Session.Close()
	return &thisGuy, nil
}

//Add a new person and link him/her to the proper family
func (*personNeo4j) Add(person *entity.Person) (*entity.Person, error) {
	err := person.IsValid()
	if err != nil {
		log.Printf("person.add.error: %v", err)
		return &entity.Person{}, err
	}

	person.ID = uuid.NewV4().String()
	person.Responsible = entity.Someone
	person.CreatedAt = time.Now()
	person.UpdatedAt = time.Now()

	return person, nil
}
